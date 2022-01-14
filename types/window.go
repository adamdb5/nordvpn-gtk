package types

import (
	"github.com/gotk3/gotk3/gtk"
	"main/util"
)

// Window contains the GTK components for the root GTKWindow.
type Window struct {
	Window     *gtk.Window
	InfoBar    *InfoBar
	ConnectTab *ConnectTab
	SessionTab *SessionTab
	AccountTab *AccountTab
	AboutTab   *AboutTab
}

// BuildWindow constructs the root GTKWindow for the application.
func BuildWindow(builder *gtk.Builder) *Window {
	window := util.BuilderGetWindow(builder, "main_window")
	window.SetTitle(AppName)

	return &Window{
		Window:     window,
		InfoBar:    BuildInfoBar(builder),
		ConnectTab: BuildConnectTab(builder),
		SessionTab: BuildSessionTab(builder),
		AccountTab: BuildAccountTab(builder),
		AboutTab:   BuildAboutTab(builder),
	}
}
