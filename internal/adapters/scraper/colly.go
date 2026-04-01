package scraper

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/antoniopantaleo/wwdc/internal/domain"
	"github.com/gocolly/colly/v2"
)

type CollyScraper struct {
	baseURL  string
	reporter domain.ProgressReporter
}

func NewCollyScraper(baseURL string, reporter domain.ProgressReporter) *CollyScraper {
	return &CollyScraper{baseURL: baseURL, reporter: reporter}
}

func (s *CollyScraper) Scrape() ([]domain.WWDCEvent, error) {
	mu := sync.Mutex{}
	eventsMap := make(map[string]*domain.WWDCEvent)
	var scrapeErr error
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
	videosScraper.WithTransport(transport)
	videosScraper.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 20})

	re := regexp.MustCompile(`wwdc(\d{4})`)

	eventsScraper.OnError(func(r *colly.Response, err error) {
		scrapeErr = fmt.Errorf("failed to fetch events list: %w", err)
	})

	videosScraper.OnError(func(r *colly.Response, err error) {
		year := r.Ctx.Get("eventYear")
		s.reporter.Warning(fmt.Sprintf("Failed to fetch event page for %s: %s", year, err))
	})

	singleVideoScraper.OnError(func(r *colly.Response, err error) {
		s.reporter.Warning(fmt.Sprintf("Failed to fetch video page %s: %s", r.Request.URL, err))
	})

	s.reporter.Info("Scraping events from Apple Developer website...")

	eventsScraper.OnHTML("a[href^=\"/videos/wwdc\"].vc-card", func(h *colly.HTMLElement) {
		eventURL := h.Request.AbsoluteURL(h.Attr("href"))
		coverURL := h.Request.AbsoluteURL(h.ChildAttr("img.vc-card__image", "src"))
		title := h.ChildText("span.vc-card__tag--event")
		matches := re.FindStringSubmatch(eventURL)
		if len(matches) < 2 {
			s.reporter.Warning(fmt.Sprintf("Could not extract year from URL: %s", eventURL))
			return
		}
		year := matches[1]
		intYear, err := strconv.Atoi(year)
		if err != nil {
			s.reporter.Warning(fmt.Sprintf("Invalid year %s in URL: %s", year, eventURL))
			return
		}
		event := &domain.WWDCEvent{
			Title:    title,
			Year:     intYear,
			CoverURL: coverURL,
			Videos:   []domain.WWDCVideo{},
		}
		mu.Lock()
		if _, exists := eventsMap[year]; !exists {
			eventsMap[year] = event
			s.reporter.Info(fmt.Sprintf("Found event: %s (%s)", title, year))
		}
		mu.Unlock()
		ctx := colly.NewContext()
		ctx.Put("eventYear", year)
		videosScraper.Request("GET", eventURL, nil, ctx, nil)
	})

	videosScraper.OnHTML("a[href^=\"/videos/play/wwdc\"].vc-card", func(h *colly.HTMLElement) {
		videoURL := h.Request.AbsoluteURL(h.Attr("href"))
		ctx := h.Request.Ctx
		singleVideoScraper.Request("GET", videoURL, nil, ctx, nil)
	})

	singleVideoScraper.OnHTML("li.download li:first-child", func(h *colly.HTMLElement) {
		videoURL := h.ChildAttr("li a", "href")
		h.Request.Ctx.Put("videoURL", videoURL)
	})

	singleVideoScraper.OnHTML("li.supplement.details", func(h *colly.HTMLElement) {
		title := h.ChildText("h1")
		content := h.ChildText("p")
		ctx := h.Request.Ctx
		ctx.Put("videoTitle", title)
		ctx.Put("videoContent", content)
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
	})

	eventsScraper.Visit(s.baseURL + "/videos")
	if scrapeErr != nil {
		return nil, scrapeErr
	}
	videosScraper.Wait()
	singleVideoScraper.Wait()
	var events []domain.WWDCEvent
	for _, event := range eventsMap {
		events = append(events, *event)
	}
	return events, nil
}
