package main

import (
	"github.com/adamdb5/opennord"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"main/types"
	"os"
)

const (
	appId = "net.adambruce.nordvpn-gtk"
)

func showInfoBar(text string, messageType gtk.MessageType, ui *types.UI) {
	ui.InfoBar.SetMessageType(messageType)
	ui.InfoBarLabel.SetText(text)
	ui.InfoBar.ShowAll()
}

func hideInfoBar(ui *types.UI) {
	ui.InfoBar.Hide()
}

func populateCountries(app *types.Application) error {
	countries, err := app.Client.Countries()
	if err != nil {
		showInfoBar(err.Error(), gtk.MESSAGE_ERROR, app.Ui)
		return err
	}

	list, err := gtk.ListStoreNew(glib.TYPE_STRING)
	for _, country := range countries.GetCountries() {
		appendIter := list.Append()
		list.SetValue(appendIter, 0, country)
	}

	renderer, _ := gtk.CellRendererTextNew()

	app.Ui.CountriesComboBoxText.SetModel(list)
	app.Ui.CountriesComboBoxText.PackStart(renderer, true)
	app.Ui.CountriesComboBoxText.SetActive(0)

	return nil
}

func populateCities(app *types.Application) error {
	cities, err := app.Client.Cities(app.Ui.CountriesComboBoxText.GetActiveText())
	if err != nil {
		log.Printf("Error: Lost connection to NordVPN daemon: %s", err)
		app.Ui.InfoBarButton.SetLabel("Reconnect")
		app.Ui.InfoBarButton.Connect("clicked", func() { connectToDaemon(app) })
		showInfoBar("Lost connection to NordVPN daemon", gtk.MESSAGE_ERROR, app.Ui)
		return err
	}

	list, err := gtk.ListStoreNew(glib.TYPE_STRING)
	for _, city := range cities.GetCities() {
		appendIter := list.Append()
		list.SetValue(appendIter, 0, city)
	}

	renderer, _ := gtk.CellRendererTextNew()

	app.Ui.CitiesComboBoxText.SetModel(list)
	app.Ui.CitiesComboBoxText.PackStart(renderer, true)
	app.Ui.CitiesComboBoxText.SetActive(0)

	return nil
}

func updateConnectionStatus(app *types.Application) error {
	status, err := app.Client.Status()
	if err != nil {
		log.Printf("Error: Lost connection to NordVPN daemon: %s", err)
		app.Ui.InfoBarButton.SetLabel("Reconnect")
		app.Ui.InfoBarButton.Connect("clicked", func() { connectToDaemon(app) })
		showInfoBar("Lost connection to NordVPN daemon", gtk.MESSAGE_ERROR, app.Ui)
		return err
	}

	app.Ui.StatusLabel.SetText(status.GetState())

	return nil
}

func connectToDaemon(app *types.Application) error {
	hideInfoBar(app.Ui)
	client, err := opennord.NewOpenNordClient()
	if err != nil {
		log.Printf("Error: Could not connect to NordVPN daemon: %s", err)
		app.Ui.InfoBarButton.SetLabel("Retry")
		showInfoBar("Could not connect to NordVPN daemon", gtk.MESSAGE_ERROR, app.Ui)
		return err
	}

	app.Client = &client

	// Check if we're already connected to the vpn
	updateConnectionStatus(app)

	// If successful, populate the countries
	populateCountries(app)

	return nil
}

func main() {
	application, _ := gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)

	application.Connect("startup", func() {
		log.Println("nordvpn-gtk startup")
	})

	application.Connect("activate", func() {
		builder, err := gtk.BuilderNewFromFile("ui/window.glade")
		if err != nil {
			log.Fatalf("Could not build interface: %s", err)
		}

		app := types.Application{
			Ui: &types.UI{
				// Global
				Window:        getWindow(builder, "main_window"),
				InfoBar:       getInfoBar(builder, "info_bar"),
				InfoBarLabel:  getLabel(builder, "info_bar_label"),
				InfoBarButton: getButton(builder, "info_bar_button"),

				// Connect
				StatusLabel:           getLabel(builder, "status_label"),
				CountriesComboBoxText: getComboBoxText(builder, "country_combo"),
				CitiesComboBoxText:    getComboBoxText(builder, "city_combo"),
			},
		}

		if connectToDaemon(&app) != nil {
			app.Ui.InfoBarButton.Connect("clicked", func() { connectToDaemon(&app) })
		}

		app.Ui.CountriesComboBoxText.Connect("changed", func() { onCountrySelected(&app) })

		app.Ui.Window.Show()
		application.AddWindow(app.Ui.Window)
	})

	os.Exit(application.Run(os.Args))
}

func getWindow(builder *gtk.Builder, name string) *gtk.Window {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.Window)
}

func getInfoBar(builder *gtk.Builder, name string) *gtk.InfoBar {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.InfoBar)
}

func getLabel(builder *gtk.Builder, name string) *gtk.Label {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.Label)
}

func getComboBoxText(builder *gtk.Builder, name string) *gtk.ComboBoxText {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.ComboBoxText)
}

func getButton(builder *gtk.Builder, name string) *gtk.Button {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.Button)
}

func onCountrySelected(app *types.Application) {
	populateCities(app)
}
