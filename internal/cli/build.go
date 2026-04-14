package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stefva/rizzume/internal/render"
	"github.com/stefva/rizzume/internal/validate"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Render, validate, and check a resume PDF",
	RunE:  runBuild,
}

var (
	buildContent       string
	buildTheme         string
	buildTemplate      string
	buildOut           string
	buildAllowOverflow bool
	buildColor         string
)

func init() {
	buildCmd.Flags().StringVar(&buildContent, "content", "", "path to content YAML file (required)")
	buildCmd.Flags().StringVar(&buildTheme, "theme", "themes/default.yaml", "path to theme YAML file")
	buildCmd.Flags().StringVar(&buildTemplate, "template", "templates/resume.typ", "path to Typst template")
	buildCmd.Flags().StringVar(&buildOut, "out", "out/resume.pdf", "output PDF path")
	buildCmd.Flags().BoolVar(&buildAllowOverflow, "allow-overflow", false, "skip the one-page check")
	buildCmd.Flags().StringVar(&buildColor, "color", "", "override accent color (hex, e.g. #E74C3C)")
	buildCmd.MarkFlagRequired("content")
}

func runBuild(cmd *cobra.Command, args []string) error {
	rootDir, err := findRootDir()
	if err != nil {
		return err
	}

	outPath := absPath(rootDir, buildOut)

	err = render.RenderPDF(render.Opts{
		ContentPath:   absPath(rootDir, buildContent),
		ThemePath:     absPath(rootDir, buildTheme),
		TemplatePath:  absPath(rootDir, buildTemplate),
		OutPath:       outPath,
		RootDir:       rootDir,
		ColorOverride: buildColor,
	})
	if err != nil {
		return fmt.Errorf("render failed: %w", err)
	}

	report, err := validate.ValidatePDF(validate.Opts{
		PDFPath:       outPath,
		AllowOverflow: buildAllowOverflow,
		OutDir:        absPath(rootDir, "out"),
	})
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if !report.PageCountOK {
		fmt.Print(report.String())
		return fmt.Errorf("content does not fit on one page (use --allow-overflow to bypass)")
	}

	fmt.Print(report.String())
	if buildColor != "" {
		fmt.Printf("Color: %s (override)\n", buildColor)
	}
	fmt.Printf("PDF:   %s\n", outPath)
	return nil
}
