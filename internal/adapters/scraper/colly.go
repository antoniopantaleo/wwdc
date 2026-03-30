package scraper

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/antoniopantaleo/wwdc/internal/domain"
	"github.com/gocolly/colly/v2"
)

type CollyScraper struct {
	baseURL string
}

func NewCollyScraper(baseURL string) *CollyScraper {
	return &CollyScraper{baseURL: baseURL}
}

func (s *CollyScraper) Scrape() ([]domain.WWDCEvent, error) {
	mu := sync.Mutex{}
	eventsMap := make(map[string]*domain.WWDCEvent)
	eventsScraper := colly.NewCollector()

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 50,
		MaxConnsPerHost:     50,
	}

	singleVideoScraper := colly.NewCollector(
		colly.Async(true),
		colly.AllowURLRevisit(),
	)
	singleVideoScraper.WithTransport(transport)
	singleVideoScraper.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 20})
	videosScraper := colly.NewCollector(
		colly.Async(true),
		colly.AllowURLRevisit(),
	)
	log.Println("Visiting", s.baseURL)
	videosScraper.WithTransport(transport)
	videosScraper.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 20})

	re := regexp.MustCompile(`wwdc(\d{4})`)
	eventsScraper.OnHTML("a[href^=\"/videos/wwdc\"].vc-card", func(h *colly.HTMLElement) {
		eventURL := h.Request.AbsoluteURL(h.Attr("href"))
		coverURL := h.Request.AbsoluteURL(h.ChildAttr("img.vc-card__image", "src"))
		title := h.ChildText("span.vc-card__tag--event")
		matches := re.FindStringSubmatch(eventURL)
		year := "unknown"
		if len(matches) > 1 {
			year = matches[1]
		}
		intYear, _ := strconv.Atoi(year)
		log.Printf("Found event: %s %s (%s) - %s\n", year, title, coverURL, eventURL)
		event := &domain.WWDCEvent{
			Title:    title,
			Year:     intYear,
			CoverURL: coverURL,
			Videos:   []domain.WWDCVideo{},
		}
		mu.Lock()
		if _, exists := eventsMap[year]; !exists {
			eventsMap[year] = event
		}
		mu.Unlock()
		ctx := colly.NewContext()
		ctx.Put("eventYear", year)
		videosScraper.Request("GET", eventURL, nil, ctx, nil)
	})

	videosScraper.OnHTML("a[href^=\"/videos/play/wwdc\"].vc-card", func(h *colly.HTMLElement) {
		videoURL := h.Request.AbsoluteURL(h.Attr("href"))
		ctx := h.Request.Ctx
		log.Printf("Found video page: %s\n", videoURL)
		singleVideoScraper.Request("GET", videoURL, nil, ctx, nil)
	})

	singleVideoScraper.OnHTML("li.download li:first-child", func(h *colly.HTMLElement) {
		videoURL := h.ChildAttr("li a", "href")
		h.Request.Ctx.Put("videoURL", videoURL)
		log.Printf("Found video URL: %s\n", videoURL)
	})

	singleVideoScraper.OnHTML("li.supplement.details", func(h *colly.HTMLElement) {
		title := h.ChildText("h1")
		content := h.ChildText("p")
		ctx := h.Request.Ctx
		ctx.Put("videoTitle", title)
		ctx.Put("videoContent", content)
		log.Printf("Found video details: %s - %s\n", title, content)
	})

	singleVideoScraper.OnScraped(func(r *colly.Response) {
		year := r.Ctx.Get("eventYear")
		videoURL := r.Ctx.Get("videoURL")
		title := r.Ctx.Get("videoTitle")
		content := r.Ctx.Get("videoContent")
		video := &domain.WWDCVideo{
			Title:    title,
			VideoURL: videoURL,
			Content:  content,
		}
		mu.Lock()
		if event, exists := eventsMap[year]; exists {
			event.Videos = append(event.Videos, *video)
		}
		mu.Unlock()
		log.Printf("%s, Finished scraping video page: %s\n", year, r.Request.URL)
	})

	eventsScraper.Visit(s.baseURL + "/videos")
	videosScraper.Wait()
	singleVideoScraper.Wait()
	var events []domain.WWDCEvent
	for _, event := range eventsMap {
		events = append(events, *event)
	}
	return events, nil
}
