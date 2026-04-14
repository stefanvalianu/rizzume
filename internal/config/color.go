package config

import (
	"fmt"
	"math"
	"strings"
)

// DeriveColors computes accent_bg and accent_fg from the primary accent color
// if they are not already set in the theme.
//
// accent    — the bright primary color (last name, company names)
// accent_bg — light pastel tint for section heading backgrounds
// accent_fg — dark shade for text on section heading backgrounds
func DeriveColors(t *Theme) {
	if t.Colors.AccentBg == "" {
		t.Colors.AccentBg = tint(t.Colors.Accent, 0.88, 0.35)
	}
	if t.Colors.AccentFg == "" {
		t.Colors.AccentFg = shade(t.Colors.Accent, 0.28)
	}
}

// tint produces a light pastel: push lightness up, pull saturation down.
func tint(hex string, targetL, satScale float64) string {
	r, g, b, err := parseHex(hex)
	if err != nil {
		return hex
	}
	h, s, _ := rgbToHSL(r, g, b)
	nr, ng, nb := hslToRGB(h, s*satScale, targetL)
	return fmt.Sprintf("#%02X%02X%02X", nr, ng, nb)
}

// shade produces a dark version: push lightness down, keep saturation.
func shade(hex string, targetL float64) string {
	r, g, b, err := parseHex(hex)
	if err != nil {
		return hex
	}
	h, s, _ := rgbToHSL(r, g, b)
	nr, ng, nb := hslToRGB(h, math.Min(1, s*1.1), targetL)
	return fmt.Sprintf("#%02X%02X%02X", nr, ng, nb)
}

func parseHex(hex string) (uint8, uint8, uint8, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0, fmt.Errorf("invalid hex color: %s", hex)
	}
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b, err
}

func rgbToHSL(r, g, b uint8) (float64, float64, float64) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	l := (max + min) / 2.0

	if max == min {
		return 0, 0, l
	}

	d := max - min
	s := d / (1 - math.Abs(2*l-1))

	var h float64
	switch max {
	case rf:
		h = (gf - bf) / d
		if gf < bf {
			h += 6
		}
	case gf:
		h = (bf-rf)/d + 2
	case bf:
		h = (rf-gf)/d + 4
	}
	h /= 6.0

	return h, s, l
}

func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	if s == 0 {
		v := uint8(math.Round(l * 255))
		return v, v, v
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	r := hueToRGB(p, q, h+1.0/3.0)
	g := hueToRGB(p, q, h)
	b := hueToRGB(p, q, h-1.0/3.0)

	return uint8(math.Round(r * 255)), uint8(math.Round(g * 255)), uint8(math.Round(b * 255))
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}
