package i18n

import (
	"os"
	"testing"
)

func TestDetectLang(t *testing.T) {
	// Save all locale env vars
	savedLANG := os.Getenv("LANG")
	savedLCALL := os.Getenv("LC_ALL")
	savedLCMESSAGES := os.Getenv("LC_MESSAGES")
	defer func() {
		os.Setenv("LANG", savedLANG)
		os.Setenv("LC_ALL", savedLCALL)
		os.Setenv("LC_MESSAGES", savedLCMESSAGES)
	}()

	// Clear all locale env vars first
	os.Unsetenv("LC_ALL")
	os.Unsetenv("LC_MESSAGES")

	os.Setenv("LANG", "zh_CN.UTF-8")
	if got := DetectLang(); got != LangZH {
		t.Errorf("expected zh, got %s", got)
	}

	os.Setenv("LANG", "en_US.UTF-8")
	if got := DetectLang(); got != LangEN {
		t.Errorf("expected en, got %s", got)
	}

	os.Unsetenv("LANG")
	if got := DetectLang(); got != LangEN {
		t.Errorf("expected en fallback, got %s", got)
	}
}

func TestNew(t *testing.T) {
	b := New(LangEN)
	if b.Lang() != LangEN {
		t.Errorf("expected en, got %s", b.Lang())
	}

	b2 := New(LangZH)
	if b2.Lang() != LangZH {
		t.Errorf("expected zh, got %s", b2.Lang())
	}
}

func TestLookup(t *testing.T) {
	b := New(LangEN)

	if got := b.T("app.name"); got != "wutils" {
		t.Errorf("expected wutils, got %s", got)
	}
	if got := b.T("app.description"); got != "Weidows Utilities" {
		t.Errorf("expected 'Weidows Utilities', got %s", got)
	}
	if got := b.T("dsg.description"); got != "Prevent hard drives from sleeping" {
		t.Errorf("expected 'Prevent hard drives from sleeping', got %s", got)
	}
	if got := b.T("nonexistent.key"); got != "nonexistent.key" {
		t.Errorf("expected key itself as fallback, got %s", got)
	}
}

func TestLookupZH(t *testing.T) {
	b := New(LangZH)

	if got := b.T("app.description"); got != "Weidows 工具箱" {
		t.Errorf("expected 'Weidows 工具箱', got %s", got)
	}
	if got := b.T("gui.start_all"); got != "全部启动" {
		t.Errorf("expected '全部启动', got %s", got)
	}
}

func TestTF(t *testing.T) {
	b := New(LangEN)

	if got := b.TF("gui.running", 3, 8); got != "3/8 running" {
		t.Errorf("expected '3/8 running', got %s", got)
	}
	if got := b.TF("nonexistent.key", "arg"); got != "nonexistent.key" {
		t.Errorf("expected key itself, got %s", got)
	}
}

func TestGlobal(t *testing.T) {
	InitGlobal(LangEN)
	if got := G("app.name"); got != "wutils" {
		t.Errorf("expected wutils, got %s", got)
	}

	InitGlobal(LangZH)
	if got := G("gui.dashboard"); got != "仪表盘" {
		t.Errorf("expected '仪表盘', got %s", got)
	}
}

func TestResolveLang(t *testing.T) {
	if got := ResolveLang("en"); got != LangEN {
		t.Errorf("expected en, got %s", got)
	}
	if got := ResolveLang("zh"); got != LangZH {
		t.Errorf("expected zh, got %s", got)
	}
	// auto depends on system, just verify it returns a valid Lang
	lang := ResolveLang("auto")
	if lang != LangEN && lang != LangZH {
		t.Errorf("expected en or zh, got %s", lang)
	}
}
