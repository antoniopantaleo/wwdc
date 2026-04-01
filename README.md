<h1 align="center">wwdc</h1>

<p align="center">
  <img src="assets/logo.png" alt="wwdc logo" width="200">
</p>

<p align="center">
  A CLI tool to scrape WWDC session videos from the Apple Developer website and export them in multiple formats.
</p>

---

## What it does

`wwdc` scrapes all WWDC events and their session videos from [developer.apple.com/videos](https://developer.apple.com/videos/) and exports them in your format of choice:

- **JSON** — structured data output to stdout or saved to a file
- **Markdown** — organized folders (one per event) containing `.md` files (one per video)

## Installation

### Building from source

Requires Go 1.25.6 or later.

```bash
git clone https://github.com/antoniopantaleo/wwdc.git
cd wwdc
go build -ldflags "-s -w -X github.com/antoniopantaleo/wwdc/cmd.version=$(cat VERSION)" -o wwdc .
```

### Using mise

If you use [mise](https://mise.jdx.dev/), all tool versions are pinned in `mise.toml`:

```bash
git clone https://github.com/antoniopantaleo/wwdc.git
cd wwdc
mise trust && install
mise run build:prod
```

### Docker

```bash
docker build -t wwdc -f Dockerfile .
```

## Usage

### Export as JSON

Print to stdout:

```bash
wwdc export json
```

Save to a file:

```bash
wwdc export json --output wwdc.json
```

### Export as Markdown

Creates a `WWDC/` directory with subfolders for each event:

```bash
wwdc export markdown
```

If your markdown viewer uses the filename as a title, you can omit the title heading from the generated files:

```bash
wwdc export markdown --omit-title
```

## Performance

`wwdc` uses concurrent HTTP requests to scrape hundreds of video pages in parallel. A full export of all WWDC events typically completes in under 20 seconds:

```
$ time wwdc export json > /dev/null

real    0m19.298s
user    0m6.00s
sys     0m1.14s
```

## License

See [LICENSE](LICENSE) for details.
