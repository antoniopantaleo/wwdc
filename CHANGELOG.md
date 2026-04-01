## [0.1.4] - 2026-04-01

### ⚙️ Miscellaneous Tasks

- Generate full changelog
## [0.1.3] - 2026-04-01

### 📚 Documentation

- Update changelog for v0.1.3

### ⚙️ Miscellaneous Tasks

- Add ci git user info
## [0.1.2] - 2026-04-01

### ⚙️ Miscellaneous Tasks

- Commit CHANGELOG.md to repo
## [0.1.1] - 2026-04-01

### ⚙️ Miscellaneous Tasks

- Use git-cliff in its own job to generate changelog
## [0.1.0] - 2026-04-01

### 🚀 Features

- Add domain models
- Add ports
- Add ScrapeAndExportUseCase
- *(json)* Add exporter
- *(cobra)* Add cobra root command using stub scraper
- *(cobra)* Check for --format flag value
- *(json)* Change indentation from tab to 2 spaces
- *(md)* Add markdown exporter
- Create separate `export` command
- Add filesystem os
- *(md)* Embed video instead of just displaying its URL
- *(colly)* Add colly scraper
- *(colly)* Use video download href instead of streaming url for better compatibility
- *(md)* Use <video> tag instead of markdown embedding for better compatibility
- *(md)* Add option to omit title during export
- *(json)* Extend option to omit title to json export format too
- Split format options in subcommands, add output flag to JSON and remove omitTitle from it
- *(md)* Sanitize path before writing to file system
- *(colly)* Add progress reporter

### 💼 Other

- Add mise
- Use git-cliff to bump version
- Separate build task in two different tasks, one for prod and one for dev
- Reduce binary size stripping debug info
- Add Dockerfile

### 📚 Documentation

- Add info to commands
- Improve long command description
- Add README

### ⚡ Performance

- *(cobra)* Use os.Stdout instead of bytes.Buffer
- *(colly)* Faster scraping using custom HTTP transport

### 🧪 Testing

- Check if returned error is the expected one
- *(json)* Add testcase for empty events
- *(md)* Add one event multiple videos test
- *(md)* Add multiple events multiple videos test
- *(md)* Add one event no videos test
- *(md)* Add no events test
- *(md)* Add wrong year test
- *(md)* Add makeDir error test
- *(md)* Add writeFile error test
- *(colly)* Improve test obtaining exact instance of event and video

### ⚙️ Miscellaneous Tasks

- Initial commit
- Create project structure
- Add cobra dependency
- *(cobra)* Manually check if --format flag is set to avoid conflicts with other commands
- Replace stub with real filesystem
- Moved composition to root command
- Use colly scraper in prod
- *(colly)* Inject logger
- Remove cliff.toml
- Add logo asset
- Add LICENSE
- Add gitignore
- Setup workflows
