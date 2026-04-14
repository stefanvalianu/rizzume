package render

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/stefva/rizzume/internal/config"
	"github.com/stefva/rizzume/internal/content"
)

type Opts struct {
	ContentPath  string
	ThemePath    string
	TemplatePath string
	OutPath      string
	RootDir      string
	ColorOverride string
}

func RenderPDF(opts Opts) error {
	c, err := content.Load(opts.ContentPath)
	if err != nil {
		return err
	}

	t, err := config.LoadTheme(opts.ThemePath)
	if err != nil {
		return err
	}

	// Apply color override if specified
	if opts.ColorOverride != "" {
		t.Colors.Accent = opts.ColorOverride
		t.Colors.AccentBg = ""
		t.Colors.AccentFg = ""
	}

	// Derive accent_dark and accent_light from accent
	config.DeriveColors(t)

	// Write intermediate JSON files
	tmpDir := filepath.Join(opts.RootDir, ".rizzume-tmp")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	contentJSON, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling content JSON: %w", err)
	}
	contentJSONPath := filepath.Join(tmpDir, "content.json")
	if err := os.WriteFile(contentJSONPath, contentJSON, 0o644); err != nil {
		return fmt.Errorf("writing content JSON: %w", err)
	}

	themeJSON, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling theme JSON: %w", err)
	}
	themeJSONPath := filepath.Join(tmpDir, "theme.json")
	if err := os.WriteFile(themeJSONPath, themeJSON, 0o644); err != nil {
		return fmt.Errorf("writing theme JSON: %w", err)
	}

	// Make output directory
	if err := os.MkdirAll(filepath.Dir(opts.OutPath), 0o755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	// Build Typst paths relative to root (prefixed with / for Typst root-relative resolution)
	contentRelPath, err := filepath.Rel(opts.RootDir, contentJSONPath)
	if err != nil {
		return fmt.Errorf("computing relative path: %w", err)
	}
	themeRelPath, err := filepath.Rel(opts.RootDir, themeJSONPath)
	if err != nil {
		return fmt.Errorf("computing relative path: %w", err)
	}

	// Invoke Typst
	fontsDir := filepath.Join(opts.RootDir, "fonts")
	args := []string{
		"compile",
		"--root", opts.RootDir,
		"--font-path", fontsDir,
		"--input", fmt.Sprintf("content-path=/%s", contentRelPath),
		"--input", fmt.Sprintf("theme-path=/%s", themeRelPath),
		opts.TemplatePath,
		opts.OutPath,
	}

	cmd := exec.Command("typst", args...)
	cmd.Dir = opts.RootDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("typst compile failed: %w\n%s", err, string(out))
	}

	// Print any warnings
	if len(out) > 0 {
		fmt.Fprintf(os.Stderr, "%s", string(out))
	}

	return nil
}
