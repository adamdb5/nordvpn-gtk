package types

import (
	"github.com/adamdb5/opennord"
	"github.com/gotk3/gotk3/gtk"
)

type Application struct {
	Client *opennord.Client
	Ui     *UI
}

type UI struct {
	Window *gtk.Window

	InfoBar       *gtk.InfoBar
	InfoBarLabel  *gtk.Label
	InfoBarButton *gtk.Button

	CountriesComboBoxText *gtk.ComboBoxText
	CitiesComboBoxText    *gtk.ComboBoxText

	StatusLabel *gtk.Label
}
