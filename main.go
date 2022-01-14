package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"main/types"
	"main/util"
	"os"
)

func main() {
	application, _ := gtk.ApplicationNew(types.AppId,
		glib.APPLICATION_FLAGS_NONE)

	application.Connect("activate", func() {
		builder, err := gtk.BuilderNewFromFile("ui/window.glade")
		if err != nil {
			util.LogFatal("Could not build interface", err)
		}

		app := types.BuildApplication(builder)
		err = app.ConnectToDaemon()
		if err != nil {
			infoBar := app.Window.InfoBar
			infoBar.Button.SetLabel("Retry")
			infoBar.Button.Connect("clicked", app.ConnectToDaemon)
			infoBar.DisplayMessage("Could not connect to NordVPN daemon",
				gtk.MESSAGE_ERROR)
		}

		go app.UpdateSessionStatus()

		gtkWindow := app.Window.Window
		gtkWindow.Show()
		application.AddWindow(gtkWindow)
	})

	os.Exit(application.Run(os.Args))
}
