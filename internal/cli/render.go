package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/stefva/rizzume/internal/render"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render a resume PDF from content and theme files",
	RunE:  runRender,
}

var (
	renderContent  string
	renderTheme    string
	renderTemplate string
	renderOut      string
	renderColor    string
)

func init() {
	renderCmd.Flags().StringVar(&renderContent, "content", "", "path to content YAML file (required)")
	renderCmd.Flags().StringVar(&renderTheme, "theme", "themes/default.yaml", "path to theme YAML file")
	renderCmd.Flags().StringVar(&renderTemplate, "template", "templates/resume.typ", "path to Typst template")
	renderCmd.Flags().StringVar(&renderOut, "out", "out/resume.pdf", "output PDF path")
	renderCmd.Flags().StringVar(&renderColor, "color", "", "override accent color (hex, e.g. #E74C3C)")
	renderCmd.MarkFlagRequired("content")
}

func runRender(cmd *cobra.Command, args []string) error {
	rootDir, err := findRootDir()
	if err != nil {
		return err
	}

	opts := render.Opts{
		ContentPath:   absPath(rootDir, renderContent),
		ThemePath:     absPath(rootDir, renderTheme),
		TemplatePath:  absPath(rootDir, renderTemplate),
		OutPath:       absPath(rootDir, renderOut),
		RootDir:       rootDir,
		ColorOverride: renderColor,
	}

	if err := render.RenderPDF(opts); err != nil {
		return err
	}

	fmt.Printf("Theme: %s\n", renderTheme)
	if renderColor != "" {
		fmt.Printf("Color: %s (override)\n", renderColor)
	}
	fmt.Printf("PDF:   %s\n", opts.OutPath)
	return nil
}

func findRootDir() (string, error) {
	// Use current directory, walking up to find go.mod as a heuristic
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Fall back to cwd
	return os.Getwd()
}

func absPath(rootDir, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(rootDir, path)
}
