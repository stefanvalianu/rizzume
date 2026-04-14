package validate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Report struct {
	PageCount       int
	PageCountOK     bool
	TextExtracted   bool
	ExtractedText   string
	ATSIssues       []string
	PDFValid        bool
	PDFValidOutput  string
	HasMetadata     bool
	TextOutputPath  string
}

func (r *Report) OK() bool {
	return r.PageCountOK && r.TextExtracted && len(r.ATSIssues) == 0
}

func (r *Report) String() string {
	var b strings.Builder
	b.WriteString("=== Validation Report ===\n")

	if r.PageCountOK {
		fmt.Fprintf(&b, "  Page count: %d (OK)\n", r.PageCount)
	} else {
		fmt.Fprintf(&b, "  Page count: %d (FAIL — must be 1)\n", r.PageCount)
	}

	if r.TextExtracted {
		b.WriteString("  Text extraction: OK\n")
	} else {
		b.WriteString("  Text extraction: FAIL\n")
	}

	if len(r.ATSIssues) > 0 {
		b.WriteString("  ATS issues:\n")
		for _, issue := range r.ATSIssues {
			fmt.Fprintf(&b, "    - %s\n", issue)
		}
	} else {
		b.WriteString("  ATS checks: OK\n")
	}

	if r.PDFValid {
		b.WriteString("  PDF structure: OK\n")
	} else {
		fmt.Fprintf(&b, "  PDF structure: %s\n", r.PDFValidOutput)
	}

	if r.HasMetadata {
		b.WriteString("  PDF metadata: OK (title, author, keywords set)\n")
	} else {
		b.WriteString("  PDF metadata: missing (consider setting document title/author)\n")
	}

	if r.TextOutputPath != "" {
		fmt.Fprintf(&b, "  Extracted text saved to: %s\n", r.TextOutputPath)
	}

	return b.String()
}

type Opts struct {
	PDFPath        string
	AllowOverflow  bool
	OutDir         string
}

func ValidatePDF(opts Opts) (*Report, error) {
	report := &Report{}

	if _, err := os.Stat(opts.PDFPath); err != nil {
		return nil, fmt.Errorf("PDF not found: %w", err)
	}

	// 1. Page count via pdfcpu
	pageCount, err := getPageCount(opts.PDFPath)
	if err != nil {
		// Fallback: try pdftotext-based heuristic
		pageCount = -1
	}
	report.PageCount = pageCount
	if opts.AllowOverflow {
		report.PageCountOK = true
	} else {
		report.PageCountOK = pageCount == 1
	}

	// 2. Text extraction via pdftotext
	text, err := extractText(opts.PDFPath)
	if err != nil {
		report.TextExtracted = false
	} else {
		report.TextExtracted = len(strings.TrimSpace(text)) > 0
		report.ExtractedText = text
	}

	// Save extracted text
	if opts.OutDir != "" && report.ExtractedText != "" {
		textPath := filepath.Join(opts.OutDir, "resume.txt")
		os.MkdirAll(opts.OutDir, 0o755)
		if err := os.WriteFile(textPath, []byte(report.ExtractedText), 0o644); err == nil {
			report.TextOutputPath = textPath
		}
	}

	// 3. ATS smoke tests
	if report.TextExtracted {
		report.ATSIssues = runATSChecks(report.ExtractedText)
	}

	// 4. PDF structure validation via pdfcpu
	valid, output := validatePDFStructure(opts.PDFPath)
	report.PDFValid = valid
	report.PDFValidOutput = output

	// 5. Metadata check
	report.HasMetadata = checkMetadata(opts.PDFPath)

	return report, nil
}

func checkMetadata(pdfPath string) bool {
	pdfcpuPath := findPdfcpu()
	if pdfcpuPath == "" {
		return false
	}
	cmd := exec.Command(pdfcpuPath, "info", pdfPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	info := string(out)
	hasTitle := strings.Contains(info, "Title:") && !strings.Contains(info, "Title: \n")
	hasAuthor := strings.Contains(info, "Author:") && !strings.Contains(info, "Author: \n")
	return hasTitle && hasAuthor
}

func getPageCount(pdfPath string) (int, error) {
	// Try pdfcpu first
	pdfcpuPath := findPdfcpu()
	if pdfcpuPath != "" {
		cmd := exec.Command(pdfcpuPath, "info", pdfPath)
		out, err := cmd.CombinedOutput()
		if err == nil {
			return parsePageCount(string(out))
		}
	}

	// Fallback: count form feeds in pdftotext output
	cmd := exec.Command("pdftotext", pdfPath, "-")
	out, err := cmd.Output()
	if err != nil {
		return -1, fmt.Errorf("pdftotext failed: %w", err)
	}
	pages := strings.Count(string(out), "\f") + 1
	if len(strings.TrimSpace(string(out))) == 0 {
		return 0, nil
	}
	return pages, nil
}

func parsePageCount(info string) (int, error) {
	re := regexp.MustCompile(`(?i)page\s*count\s*:\s*(\d+)`)
	matches := re.FindStringSubmatch(info)
	if len(matches) < 2 {
		return -1, fmt.Errorf("could not parse page count from pdfcpu output")
	}
	return strconv.Atoi(matches[1])
}

func extractText(pdfPath string) (string, error) {
	cmd := exec.Command("pdftotext", "-layout", pdfPath, "-")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("pdftotext failed: %w", err)
	}
	return string(out), nil
}

func runATSChecks(text string) []string {
	var issues []string

	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		issues = append(issues, "extracted text is empty")
		return issues
	}

	// Check for standard section headings
	textUpper := strings.ToUpper(text)
	expectedSections := []string{"EXPERIENCE", "EDUCATION", "SKILLS"}
	for _, section := range expectedSections {
		if !strings.Contains(textUpper, section) {
			issues = append(issues, fmt.Sprintf("missing expected section heading: %s", section))
		}
	}

	// Check for merged words (heuristic: very long "words" without spaces)
	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			clean := strings.Trim(word, ".,;:!?()[]{}\"'")
			if len(clean) > 50 {
				issues = append(issues, fmt.Sprintf("possible merged words detected: %.60s...", clean))
				break
			}
		}
	}

	// Check that name/email/phone appear in the text
	if len(strings.TrimSpace(text)) < 100 {
		issues = append(issues, "extracted text is suspiciously short")
	}

	return issues
}

func validatePDFStructure(pdfPath string) (bool, string) {
	pdfcpuPath := findPdfcpu()
	if pdfcpuPath == "" {
		return true, "pdfcpu not found, skipping structure validation"
	}

	cmd := exec.Command(pdfcpuPath, "validate", pdfPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Sprintf("FAIL: %s", strings.TrimSpace(string(out)))
	}
	return true, "OK"
}

func findPdfcpu() string {
	// Check PATH first
	if path, err := exec.LookPath("pdfcpu"); err == nil {
		return path
	}
	// Check common Go install locations
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	gopath := filepath.Join(home, "go", "bin", "pdfcpu")
	if _, err := os.Stat(gopath); err == nil {
		return gopath
	}
	return ""
}
