package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/stefanvalianu/rizzume/internal/validate"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a generated resume PDF",
	RunE:  runValidate,
}

var (
	validatePDF           string
	validateAllowOverflow bool
)

func init() {
	validateCmd.Flags().StringVar(&validatePDF, "pdf", "out/resume.pdf", "path to PDF to validate")
	validateCmd.Flags().BoolVar(&validateAllowOverflow, "allow-overflow", false, "skip the one-page check")
}

func runValidate(cmd *cobra.Command, args []string) error {
	rootDir, err := findRootDir()
	if err != nil {
		return err
	}

	pdfPath := absPath(rootDir, validatePDF)
	outDir := filepath.Dir(pdfPath)

	report, err := validate.ValidatePDF(validate.Opts{
		PDFPath:       pdfPath,
		AllowOverflow: validateAllowOverflow,
		OutDir:        outDir,
	})
	if err != nil {
		return err
	}

	fmt.Print(report.String())

	if !report.OK() {
		return fmt.Errorf("validation failed")
	}

	return nil
}
