package types

import (
	"errors"
	"github.com/gotk3/gotk3/gtk"
	"main/util"
)

// ConnectTab contains the GTK components for the 'Connect' GTKNotebook page.
type ConnectTab struct {
	StatusLabel           *gtk.Label
	CountriesComboBoxText *gtk.ComboBoxText
	CitiesComboBoxText    *gtk.ComboBoxText
	GroupsComboBoxText    *gtk.ComboBoxText
	ServerEntry           *gtk.Entry
	DisconnectButton      *gtk.Button
	CountryConnectButton  *gtk.Button
	CityConnectButton     *gtk.Button
	GroupConnectButton    *gtk.Button
	ServerConnectButton   *gtk.Button
	BestConnectButton     *gtk.Button
	SaveButton            *gtk.Button
}

// BuildConnectTab constructs the GTKNotebook page for the 'Connect' tab from
// the provided builder.
func BuildConnectTab(builder *gtk.Builder) *ConnectTab {
	return &ConnectTab{
		StatusLabel: util.BuilderGetLabel(builder,
			"connect_status_label"),
		CountriesComboBoxText: util.BuilderGetComboBoxText(builder,
			"connect_country_combo_text"),
		CitiesComboBoxText: util.BuilderGetComboBoxText(builder,
			"connect_city_combo_text"),
		GroupsComboBoxText: util.BuilderGetComboBoxText(builder,
			"connect_group_combo_text"),
		ServerEntry: util.BuilderGetEntry(builder,
			"connect_server_entry"),
		DisconnectButton: util.BuilderGetButton(builder,
			"connect_disconnect_button"),
		CountryConnectButton: util.BuilderGetButton(builder,
			"connect_country_connect_button"),
		CityConnectButton: util.BuilderGetButton(builder,
			"connect_city_connect_button"),
		GroupConnectButton: util.BuilderGetButton(builder,
			"connect_group_connect_button"),
		ServerConnectButton: util.BuilderGetButton(builder,
			"connect_server_connect_button"),
		BestConnectButton: util.BuilderGetButton(builder,
			"connect_best_connect_button"),
		SaveButton: util.BuilderGetButton(builder,
			"connect_save_button"),
	}
}

func ConnectSaveClicked(app *Application) error {
	serverText, _ := app.Window.ConnectTab.ServerEntry.GetText()
	app.Config.Connect = &Connect{
		Country: app.Window.ConnectTab.CountriesComboBoxText.GetActiveText(),
		City:    app.Window.ConnectTab.CitiesComboBoxText.GetActiveText(),
		Group:   app.Window.ConnectTab.GroupsComboBoxText.GetActiveText(),
		Server:  serverText,
	}
	return SaveConfig(app)
}

// DisconnectClicked is invoked whenever the 'Disconnect' button on the
// 'Connect' tab is clicked. This function will disconnect the user from their
// current VPN session.
func DisconnectClicked(app *Application) error {
	status, err := app.Client.Status()
	infoBar := app.Window.InfoBar

	if err != nil {
		util.LogError("Lost connection to NordVPN daemon", err)
		infoBar.Button.SetLabel("Reconnect")
		infoBar.Button.Connect("clicked", app.ConnectToDaemon)
		infoBar.DisplayMessage("Lost connection to NordVPN daemon",
			gtk.MESSAGE_ERROR)
		return err
	}

	// The user probably used some other tool to disconnect
	if status.GetState() == "Disconnected" {
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("You are not connected to a VPN",
			gtk.MESSAGE_ERROR)
		return nil
	}

	err = app.Client.Disconnect()
	if err != nil {
		util.LogError("Could not disconnect from VPN", err)
		infoBar.Button.SetLabel("Reconnect")
		infoBar.Button.Connect("clicked", app.ConnectToDaemon)
		infoBar.DisplayMessage("Could not disconnect from VPN",
			gtk.MESSAGE_ERROR)
		return err
	}

	infoBar.Button.SetLabel("Dismiss")
	infoBar.Button.Connect("clicked", infoBar.HideMessage)
	infoBar.DisplayMessage("Successfully disconnected from "+status.
		GetHostname(), gtk.MESSAGE_INFO)
	_ = app.UpdateConnectionStatus()
	return nil
}

// CountrySelected is invoked whenever a country is selected from the country
// combo box on the 'Connect' tab. This function updates the cities combo box
// with the relevant cities.
func CountrySelected(app *Application) {
	_ = app.PopulateCities()
}

// ConnectToCountry is invoked whenever the 'Connect to Country' button on the
// 'Connect' tab is clicked. This function will attempt to connect the user to
// the chosen country.
func ConnectToCountry(app *Application) error {
	return app.Connect(
		app.Window.ConnectTab.CountriesComboBoxText.GetActiveText())
}

// ConnectToCity is invoked whenever the 'Connect to City' button on the
// 'Connect' tab is clicked. This function will attempt to connect the user to
// the chosen city.
func ConnectToCity(app *Application) error {
	return app.Connect(
		app.Window.ConnectTab.CitiesComboBoxText.GetActiveText())
}

// ConnectToGroup is invoked whenever the 'Connect to Group' button on the
// 'Connect' tab is clicked. This function will attempt to connect the user to
// the chosen group.
func ConnectToGroup(app *Application) error {
	return app.Connect(
		app.Window.ConnectTab.GroupsComboBoxText.GetActiveText())
}

// ConnectToServer is invoked whenever the 'Connect to Server' button on the
// 'Connect' tab is clicked. This function will attempt to connect the user to
// the specified server.
func ConnectToServer(app *Application) error {
	text, _ := app.Window.ConnectTab.ServerEntry.GetText()
	if len(text) == 0 {
		util.LogError("No server specified", nil)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("No server specified", gtk.MESSAGE_ERROR)
		return errors.New("no server specified")
	}

	return app.Connect(text)
}
