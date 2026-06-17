package main

import (
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed icon.ico
var trayIcon []byte

func (a *App) runSystray() {
	systray.Run(func() {
		systray.SetIcon(trayIcon)
		systray.SetTooltip("Monitor Screen")

		mShow := systray.AddMenuItem("Show", "Show window")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Exit application")

		go func() {
			for {
				select {
				case <-mShow.ClickedCh:
					a.ShowWindow()
				case <-mQuit.ClickedCh:
					systray.Quit()
					a.QuitApp()
				}
			}
		}()
	}, func() {
		// on exit
	})
}
