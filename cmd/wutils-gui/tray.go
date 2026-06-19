package main

import (
	"log"

	"github.com/Weidows/wutils/internal/i18n"
	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// trayReady is called by systray when the tray menu is ready.
func (a *App) trayReady() {
	systray.SetIcon(trayIconData)
	systray.SetTitle("wutils")
	systray.SetTooltip("wutils — " + i18n.G("app.description"))

	// --- Menu items ---
	mShow := systray.AddMenuItem(i18n.G("gui.tray_show"), "Show wutils window")
	mHide := systray.AddMenuItem(i18n.G("gui.tray_hide"), "Hide wutils window to tray")
	systray.AddSeparator()
	mStartAll := systray.AddMenuItem(i18n.G("gui.start_all"), "Start all background services")
	mStopAll := systray.AddMenuItem(i18n.G("gui.stop_all"), "Stop all running services")
	systray.AddSeparator()
	mDashboard := systray.AddMenuItem(i18n.G("gui.tray_dashboard"), "Open the main dashboard")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem(i18n.G("gui.tray_quit"), "Quit wutils")

	// --- Event handlers (goroutine per channel to avoid blocking) ---

	go func() {
		for range mShow.ClickedCh {
			if a.ctx != nil {
				runtime.WindowShow(a.ctx)
				runtime.WindowCenter(a.ctx)
			}
		}
	}()

	go func() {
		for range mHide.ClickedCh {
			if a.ctx != nil {
				runtime.WindowHide(a.ctx)
			}
		}
	}()

	go func() {
		for range mStartAll.ClickedCh {
			a.registry.StartAll()
		}
	}()

	go func() {
		for range mStopAll.ClickedCh {
			a.registry.StopAll()
		}
	}()

	go func() {
		for range mDashboard.ClickedCh {
			if a.ctx != nil {
				runtime.WindowShow(a.ctx)
				runtime.WindowCenter(a.ctx)
				runtime.EventsEmit(a.ctx, "navigate", "dashboard")
			}
		}
	}()

	go func() {
		<-mQuit.ClickedCh
		if a.ctx != nil {
			runtime.Quit(a.ctx)
		}
		systray.Quit()
	}()
}

// trayExit is called when the tray is exiting.
func (a *App) trayExit() {
	// Cleanup — services already stopped via shutdown
}

// startTray launches the system tray in a background goroutine.
// Must be called before wails.Run() on Windows.
func (a *App) startTray() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("systray: recovered from panic: %v", r)
			}
		}()
		systray.Run(a.trayReady, a.trayExit)
	}()
}
