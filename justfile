content := "examples/content.sample.yaml"
out     := "out/resume.pdf"

# Build the CLI binary
build:
    go build -o resume-gen ./cmd

# Render a resume PDF (accepts --theme, --color, --out)
render *args: build
    ./resume-gen render --content {{content}} {{args}}

# Full build: render + validate + auto-fit across profiles (accepts --color)
resume *args: build
    ./resume-gen build --content {{content}} {{args}}

# Validate an existing PDF
validate pdf=out: build
    ./resume-gen validate --pdf {{pdf}}

# Render then open the PDF (accepts --theme, --color)
preview *args: build
    ./resume-gen render --content {{content}} {{args}}
    xdg-open {{out}} 2>/dev/null || open {{out}} 2>/dev/null || echo "open {{out}} manually"

# Remove build artifacts
clean:
    rm -rf resume-gen out/ .rizzume-tmp/

# Check aspell availability for spellcheck
[private]
check-aspell:
    @which aspell > /dev/null 2>&1 || (echo "aspell not found — install with: sudo dnf install aspell aspell-en" && exit 1)

# Spellcheck the extracted resume text
spellcheck: build check-aspell
    @test -f out/resume.txt || (echo "Run 'just resume' first to generate out/resume.txt" && exit 1)
    @echo "=== Spellcheck ==="
    @cat out/resume.txt | aspell list --lang=en --personal=/dev/null 2>/dev/null | sort -u | while read word; do echo "  misspelled: $word"; done || echo "  No issues found"
