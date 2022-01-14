package types

import (
	"github.com/gotk3/gotk3/gtk"
	"main/util"
)

// InfoBar contains the GTK components for the GTKInfoBar that is displayed at
// the top of the application window.
type InfoBar struct {
	InfoBar *gtk.InfoBar
	Label   *gtk.Label
	Button  *gtk.Button
}

// BuildInfoBar constructs the GTKInfoBar for the application.
func BuildInfoBar(builder *gtk.Builder) *InfoBar {
	return &InfoBar{
		InfoBar: util.BuilderGetInfoBar(builder, "info_bar"),
		Label:   util.BuilderGetLabel(builder, "info_bar_label"),
		Button:  util.BuilderGetButton(builder, "info_bar_button"),
	}
}

// DisplayMessage displays the info bar with the specified text and type.
func (infoBar InfoBar) DisplayMessage(text string, messageType gtk.MessageType) {
	infoBar.HideMessage()
	infoBar.Label.SetText(text)
	infoBar.InfoBar.SetMessageType(messageType)
	infoBar.InfoBar.Show()
}

// HideMessage hides the info bar.
func (infoBar InfoBar) HideMessage() {
	infoBar.InfoBar.Hide()
}
