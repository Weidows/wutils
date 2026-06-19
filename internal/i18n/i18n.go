// Package i18n provides internationalization for wutils.
//
// Usage:
//
//	bundle := i18n.New(i18n.DetectLang())
//	fmt.Println(bundle.T("app.name"))
//	fmt.Println(bundle.TF("dsg.writing_to", "E:"))
//
// Translation files are embedded in the binary at compile time.
package i18n

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Lang represents a language code.
type Lang string

const (
	LangEN  Lang = "en"
	LangZH Lang = "zh"
)

//go:embed translations/*.yaml
var translationFS embed.FS

// Bundle holds translations for a single language.
type Bundle struct {
	lang Lang
	data map[string]interface{}
}

// globalBundle is the application-wide i18n bundle.
// Set during initialization via InitGlobal.
var globalBundle = New(LangEN)

// InitGlobal sets the global bundle. Called once at startup.
func InitGlobal(lang Lang) {
	globalBundle = New(lang)
}

// Global returns the application-wide translation bundle.
func Global() *Bundle {
	return globalBundle
}

// G is a short alias for Global().T.
func G(key string) string {
	return globalBundle.T(key)
}

// GF is a short alias for Global().TF.
func GF(key string, args ...interface{}) string {
	return globalBundle.TF(key, args...)
}

// ResolveLang picks the effective language from config and system.
// If cfgLocale is "auto" (or empty), uses DetectLang().
// Otherwise uses the configured locale.
func ResolveLang(cfgLocale string) Lang {
	switch cfgLocale {
	case "en":
		return LangEN
	case "zh":
		return LangZH
	default: // "auto" or empty
		return DetectLang()
	}
}

// DetectLang detects the system language from environment variables.
// On Windows it checks LANG env var; on POSIX it checks LC_ALL/LC_MESSAGES/LANG.
// Returns LangEN as fallback.
func DetectLang() Lang {
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		val := os.Getenv(env)
		if val == "" {
			continue
		}
		val = strings.ToLower(val)
		switch {
		case strings.HasPrefix(val, "zh"):
			return LangZH
		case strings.HasPrefix(val, "en"):
			return LangEN
		}
	}
	// Windows fallback: check LANG (often set by Git Bash / MSYS2)
	if v := os.Getenv("LANG"); v != "" {
		if strings.HasPrefix(strings.ToLower(v), "zh") {
			return LangZH
		}
	}
	return LangEN
}

// New creates a translation bundle for the given language.
// Falls back to English if the language is not supported.
func New(lang Lang) *Bundle {
	b := &Bundle{lang: lang}

	// Try loading the requested language
	if lang != LangEN {
		if data := loadLang(lang); data != nil {
			b.data = data
			return b
		}
	}

	// Fallback to English
	data := loadLang(LangEN)
	if data == nil {
		// Last resort: empty map so lookups don't panic
		b.data = make(map[string]interface{})
	} else {
		b.data = data
	}
	return b
}

func loadLang(lang Lang) map[string]interface{} {
	path := fmt.Sprintf("translations/%s.yaml", lang)
	data, err := translationFS.ReadFile(path)
	if err != nil {
		return nil
	}
	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}

// Lang returns the bundle's language code.
func (b *Bundle) Lang() Lang { return b.lang }

// T looks up a dotted key and returns the string value.
// Returns the key itself if not found.
//
// Examples:
//
//	b.T("app.name")              // "wutils"
//	b.T("dsg.description")       // "Disk Sleep Guard"
func (b *Bundle) T(key string) string {
	val := lookup(b.data, key)
	if val == "" {
		return key
	}
	return val
}

// TF looks up a key and formats it with fmt.Sprintf.
func (b *Bundle) TF(key string, args ...interface{}) string {
	tpl := b.T(key)
	if tpl == key {
		return key
	}
	return fmt.Sprintf(tpl, args...)
}

// lookup traverses the nested map for a dotted key.
func lookup(data map[string]interface{}, key string) string {
	parts := strings.Split(key, ".")
	current := data
	for i, part := range parts {
		if i == len(parts)-1 {
			if v, ok := current[part]; ok {
				if s, ok := v.(string); ok {
					return s
				}
				// Try formatting non-string values
				return fmt.Sprintf("%v", v)
			}
			return ""
		}
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return ""
		}
	}
	return ""
}
