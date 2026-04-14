package content

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Content struct {
	Basics     Basics      `yaml:"basics" json:"basics"`
	Summary    *Summary    `yaml:"summary,omitempty" json:"summary,omitempty"`
	Skills     *Skills     `yaml:"skills,omitempty" json:"skills,omitempty"`
	Experience []Company   `yaml:"experience,omitempty" json:"experience,omitempty"`
	Education  []Education `yaml:"education,omitempty" json:"education,omitempty"`
}

type Basics struct {
	Name     string `yaml:"name" json:"name"`
	Email    string `yaml:"email,omitempty" json:"email,omitempty"`
	Phone    string `yaml:"phone,omitempty" json:"phone,omitempty"`
	Location string `yaml:"location,omitempty" json:"location,omitempty"`
	LinkedIn string `yaml:"linkedin,omitempty" json:"linkedin,omitempty"`
	Website  string `yaml:"website,omitempty" json:"website,omitempty"`
	GitHub   string `yaml:"github,omitempty" json:"github,omitempty"`
}

type Summary struct {
	Text string `yaml:"text" json:"text"`
}

type Skills struct {
	Groups []SkillGroup `yaml:"groups" json:"groups"`
}

type SkillGroup struct {
	Name  string `yaml:"name" json:"name"`
	Items string `yaml:"items" json:"items"`
}

type Company struct {
	Company  string `yaml:"company" json:"company"`
	URL      string `yaml:"url,omitempty" json:"url,omitempty"`
	Location string `yaml:"location" json:"location"`
	Start    string `yaml:"start" json:"start"`
	End      string `yaml:"end" json:"end"`
	Roles    []Role `yaml:"roles" json:"roles"`
}

type Role struct {
	Title   string   `yaml:"title" json:"title"`
	Team    string   `yaml:"team,omitempty" json:"team,omitempty"`
	Start   string   `yaml:"start,omitempty" json:"start,omitempty"`
	End     string   `yaml:"end,omitempty" json:"end,omitempty"`
	Bullets []string `yaml:"bullets,omitempty" json:"bullets,omitempty"`
}

type Education struct {
	Degree   string `yaml:"degree" json:"degree"`
	School   string `yaml:"school" json:"school"`
	Location string `yaml:"location,omitempty" json:"location,omitempty"`
	Year     string `yaml:"year" json:"year"`
}

func Load(path string) (*Content, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading content file: %w", err)
	}

	var c Content
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parsing content YAML: %w", err)
	}

	if c.Basics.Name == "" {
		return nil, fmt.Errorf("content validation: basics.name is required")
	}

	return &c, nil
}
