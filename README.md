# rizzume

ATS-friendly single-page resume generator. Structured YAML content + theme tokens in, deterministic PDF out.

Built with **Typst** (rendering) and **Go** (orchestration/validation).

## Quick Start

```bash
just render            # render with sample content + default theme
just resume            # render + validate + auto-fit
just validate          # validate out/resume.pdf
just preview           # render then open the PDF
just clean             # remove build artifacts
```

Output goes to `out/resume.pdf` by default. All `just` recipes accept extra flags, e.g. `just render --theme themes/tight.yaml`.

## Prerequisites

- **Go** 1.21+
- **just** (`cargo install just` or package manager)
- **Typst** 0.11+ (`cargo install typst-cli`)
- **pdftotext** (from poppler-utils) — for text extraction validation
- **pdfcpu** (optional) — `go install github.com/pdfcpu/pdfcpu/cmd/pdfcpu@latest`

## How It Works

1. You write your resume content in a YAML file (`examples/content.sample.yaml`)
2. You pick or customize a theme (`themes/default.yaml`)
3. The CLI converts YAML → JSON, invokes Typst with the template, and validates the output
4. The `build` command will fail if the content overflows a single page with the specified theme (can be bypassed)

## Project Structure

```
cmd/resume-gen/       Go CLI entry point
internal/
  cli/                CLI commands (render, validate, build)
  content/            Content schema and YAML loading
  config/             Theme schema and loading
  render/             Typst invocation
  validate/           PDF validation (page count, text extraction, ATS checks)
templates/
  resume.typ          Main Typst template
  partials/           Template components (header, section, skills, experience, education)
themes/
  default.yaml        Normal spacing/sizing
examples/
  content.sample.yaml Sample resume content
```

## CLI Commands

### `render`

Renders a single PDF from content + theme.

```bash
./resume-gen render \
  --content examples/content.sample.yaml \
  --theme themes/default.yaml \
  --template templates/resume.typ \
  --out out/resume.pdf
```

### `validate`

Validates a generated PDF: page count, text extraction, ATS checks, PDF structure.

```bash
./resume-gen validate --pdf out/resume.pdf
./resume-gen validate --pdf out/resume.pdf --allow-overflow
```

### `build`

Render + validate. Fails if the resume overflows one page (unless `--allow-overflow`).

```bash
./resume-gen build \
  --content examples/content.sample.yaml \
  --theme themes/default.yaml \
  --out out/resume.pdf
```

## Content Schema

See `examples/content.sample.yaml` for the full schema. Key sections:

- `basics` — name, email, phone, location, links
- `summary` — short professional summary
- `skills` — grouped skill categories
- `experience` — companies with multiple roles, each with optional bullets
- `education` — degrees and schools

## Theme Tokens

Themes control all visual properties: fonts, sizes, colors, spacing, divider styles.

The built-in default theme is the one we intend to use for resume generation, but more themes can be created.

## ATS Design

The template follows ATS-safe practices:
- Single-column layout with linear reading order
- No tables, text boxes, or floating frames for content
- Standard section headings (Summary, Skills, Experience, Education)
- Real text for all content (no images/icons for meaning)
- Clean text extraction verified by `pdftotext`
