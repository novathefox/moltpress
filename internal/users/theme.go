package users

import (
	"errors"
	"regexp"
	"strings"
)

// Theme errors
var (
	ErrCSSBlocked        = errors.New("CSS contains blocked patterns")
	ErrCSSTooLarge       = errors.New("CSS exceeds maximum size (10KB)")
	ErrInvalidFontPreset = errors.New("invalid font preset")
	ErrInvalidHexColor   = errors.New("invalid hex color format")
)

// ValidFontPresets contains the allowed font preset names
var ValidFontPresets = []string{
	"inter",
	"georgia",
	"playfair",
	"roboto",
	"lora",
	"montserrat",
	"merriweather",
	"source-code-pro",
	"oswald",
	"raleway",
}

// AllowedCSSProperties contains the whitelist of allowed CSS properties
var AllowedCSSProperties = map[string]bool{
	"background-color": true,
	"color":            true,
	"font-family":      true,
	"font-size":        true,
	"font-weight":      true,
	"text-align":       true,
	"text-decoration":  true,
	"line-height":      true,
	"letter-spacing":   true,
	"border-color":     true,
	"border-radius":    true,
	"padding":          true,
	"padding-top":      true,
	"padding-bottom":   true,
	"padding-left":     true,
	"padding-right":    true,
	"margin":           true,
	"margin-top":       true,
	"margin-bottom":    true,
	"margin-left":      true,
	"margin-right":     true,
	"opacity":          true,
	"box-shadow":       true,
}

// BlockedCSSPatterns contains regex patterns for dangerous CSS
var BlockedCSSPatterns = []string{
	`url\s*\(`,
	`@import`,
	`expression\s*\(`,
	`javascript:`,
	`-moz-binding`,
	`behavior\s*:`,
	`position\s*:\s*(fixed|absolute)`,
}

// ThemeSettings represents user-customizable theme configuration
type ThemeSettings struct {
	Colors    *ThemeColors  `json:"colors,omitempty"`
	Fonts     *ThemeFonts   `json:"fonts,omitempty"`
	Toggles   *ThemeToggles `json:"toggles,omitempty"`
	CustomCSS *string       `json:"custom_css,omitempty"`
}

// ThemeColors contains color customization options
type ThemeColors struct {
	Background *string `json:"background,omitempty"`
	Text       *string `json:"text,omitempty"`
	Accent     *string `json:"accent,omitempty"`
	Link       *string `json:"link,omitempty"`
	Title      *string `json:"title,omitempty"`
}

// ThemeFonts contains font customization options
type ThemeFonts struct {
	Title *string `json:"title,omitempty"` // FontPreset
	Body  *string `json:"body,omitempty"`  // FontPreset
}

// ThemeToggles contains visibility toggles for profile elements
type ThemeToggles struct {
	ShowAvatar        *bool `json:"show_avatar,omitempty"`
	ShowStats         *bool `json:"show_stats,omitempty"`
	ShowFollowerCount *bool `json:"show_follower_count,omitempty"`
	ShowBio           *bool `json:"show_bio,omitempty"`
}

// IsValidFontPreset checks if a font name is in the allowed presets list
func IsValidFontPreset(font string) bool {
	for _, preset := range ValidFontPresets {
		if strings.EqualFold(font, preset) {
			return true
		}
	}
	return false
}

// IsValidHexColor validates hex color format (#RGB or #RRGGBB)
func IsValidHexColor(color string) bool {
	if color == "" {
		return false
	}
	color = strings.TrimPrefix(color, "#")
	if len(color) != 3 && len(color) != 6 {
		return false
	}
	for _, c := range color {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// Validate checks all theme settings for validity
func (ts *ThemeSettings) Validate() error {
	if ts == nil {
		return nil
	}

	// Validate fonts
	if ts.Fonts != nil {
		if ts.Fonts.Title != nil && !IsValidFontPreset(*ts.Fonts.Title) {
			return ErrInvalidFontPreset
		}
		if ts.Fonts.Body != nil && !IsValidFontPreset(*ts.Fonts.Body) {
			return ErrInvalidFontPreset
		}
	}

	// Validate colors
	if ts.Colors != nil {
		if ts.Colors.Background != nil && !IsValidHexColor(*ts.Colors.Background) {
			return ErrInvalidHexColor
		}
		if ts.Colors.Text != nil && !IsValidHexColor(*ts.Colors.Text) {
			return ErrInvalidHexColor
		}
		if ts.Colors.Accent != nil && !IsValidHexColor(*ts.Colors.Accent) {
			return ErrInvalidHexColor
		}
		if ts.Colors.Link != nil && !IsValidHexColor(*ts.Colors.Link) {
			return ErrInvalidHexColor
		}
		if ts.Colors.Title != nil && !IsValidHexColor(*ts.Colors.Title) {
			return ErrInvalidHexColor
		}
	}

	// Validate custom CSS
	if ts.CustomCSS != nil {
		_, err := SanitizeCSS(*ts.CustomCSS)
		if err != nil {
			return err
		}
	}

	return nil
}

// SanitizeCSS validates and sanitizes custom CSS input.
// It checks input size (max 10KB), blocks dangerous patterns,
// and filters to only allow whitelisted CSS properties.
func SanitizeCSS(input string) (string, error) {
	const maxSize = 10 * 1024

	if len(input) > maxSize {
		return "", ErrCSSTooLarge
	}

	for _, pattern := range BlockedCSSPatterns {
		matched, err := regexp.MatchString(pattern, input)
		if err != nil {
			return "", err
		}
		if matched {
			return "", ErrCSSBlocked
		}
	}

	var sanitizedRules []string
	rules := strings.Split(input, ";")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		colonIdx := strings.Index(rule, ":")
		if colonIdx == -1 {
			continue
		}

		property := strings.TrimSpace(strings.ToLower(rule[:colonIdx]))

		if AllowedCSSProperties[property] {
			sanitizedRules = append(sanitizedRules, rule)
		}
	}

	if len(sanitizedRules) == 0 {
		return "", nil
	}

	return strings.Join(sanitizedRules, "; ") + ";", nil
}
