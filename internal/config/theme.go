package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Theme struct {
	Page       PageSettings `yaml:"page" json:"page"`
	Typography Typography   `yaml:"typography" json:"typography"`
	Colors     Colors       `yaml:"colors" json:"colors"`
	Spacing    Spacing      `yaml:"spacing" json:"spacing"`
}

type PageSettings struct {
	Width   string `yaml:"width" json:"width"`
	Height  string `yaml:"height" json:"height"`
	MarginX string `yaml:"margin_x" json:"margin_x"`
	MarginY string `yaml:"margin_y" json:"margin_y"`
}

type Typography struct {
	FontBody        string  `yaml:"font_body" json:"font_body"`
	FontHeading     string  `yaml:"font_heading" json:"font_heading"`
	FontSizeBase    string  `yaml:"font_size_base" json:"font_size_base"`
	FontSizeName    string  `yaml:"font_size_name" json:"font_size_name"`
	FontSizeSection string  `yaml:"font_size_section" json:"font_size_section"`
	FontSizeSmall   string  `yaml:"font_size_small" json:"font_size_small"`
	LineHeight      float64 `yaml:"line_height" json:"line_height"`
}

type Colors struct {
	Text     string `yaml:"text" json:"text"`
	Muted    string `yaml:"muted" json:"muted"`
	Accent   string `yaml:"accent" json:"accent"`
	AccentBg string `yaml:"accent_bg,omitempty" json:"accent_bg"`
	AccentFg string `yaml:"accent_fg,omitempty" json:"accent_fg"`
}

type Spacing struct {
	SectionGap     string `yaml:"section_gap" json:"section_gap"`
	ItemGap        string `yaml:"item_gap" json:"item_gap"`
	BulletGap      string `yaml:"bullet_gap" json:"bullet_gap"`
	HeaderGap      string `yaml:"header_gap" json:"header_gap"`
	SkillsLabelW   string `yaml:"skills_label_width" json:"skills_label_width"`
}

func LoadTheme(path string) (*Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading theme file: %w", err)
	}

	var t Theme
	if err := yaml.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parsing theme YAML: %w", err)
	}

	if t.Page.Width == "" || t.Page.Height == "" {
		return nil, fmt.Errorf("theme validation: page.width and page.height are required")
	}

	return &t, nil
}
