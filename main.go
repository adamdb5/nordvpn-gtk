package main

import (
	"errors"
	"fmt"
	"github.com/adamdb5/opennord"
	"github.com/adamdb5/opennord/pb"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"io"
	"log"
	"main/types"
	"math"
	"os"
	"strings"
	"time"
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

func populateGroups(app *types.Application) error {
	groups, err := app.Client.Groups(&pb.GroupsRequest{
		Protocol:  pb.ProtocolEnum_UDP,
		Obfuscate: false,
	})
	if err != nil {
		log.Printf("Error: Lost connection to NordVPN daemon: %s", err)
		app.Ui.InfoBarButton.SetLabel("Reconnect")
		app.Ui.InfoBarButton.Connect("clicked", func() { connectToDaemon(app) })
		showInfoBar("Lost connection to NordVPN daemon", gtk.MESSAGE_ERROR, app.Ui)
		return err
	}

	list, err := gtk.ListStoreNew(glib.TYPE_STRING)
	for _, group := range groups.GetGroups() {
		appendIter := list.Append()
		list.SetValue(appendIter, 0, group)
	}

	renderer, _ := gtk.CellRendererTextNew()

	app.Ui.GroupsComboBoxText.SetModel(list)
	app.Ui.GroupsComboBoxText.PackStart(renderer, true)
	app.Ui.GroupsComboBoxText.SetActive(0)

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

	connectedText := status.GetState()
	if connectedText == "Connected" {
		app.Ui.StatusLabel.SetMarkup("<span foreground=\"green\">" + connectedText + "</span>")
		app.Ui.DisconnectButton.SetSensitive(true)
	} else {
		app.Ui.StatusLabel.SetMarkup("<span foreground=\"red\">" + connectedText + "</span>")
		app.Ui.DisconnectButton.SetSensitive(false)
	}

	// Update session tab
	//onSessionRefresh(app)

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
	populateCities(app)
	populateGroups(app)

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
				GroupsComboBoxText:    getComboBoxText(builder, "group_combo"),
				ServerEntry:           getEntry(builder, "server_entry"),
				DisconnectButton:      getButton(builder, "disconnect_button"),
				CountryConnectButton:  getButton(builder, "country_connect_button"),
				CityConnectButton:     getButton(builder, "city_connect_button"),
				GroupConnectButton:    getButton(builder, "group_connect_button"),
				ServerConnectButton:   getButton(builder, "server_connect_button"),

				// Session
				SessionStatusLabel:        getLabel(builder, "session_status_label"),
				SessionServerLabel:        getLabel(builder, "session_server_label"),
				SessionCountryLabel:       getLabel(builder, "session_country_label"),
				SessionCityLabel:          getLabel(builder, "session_city_label"),
				SessionServerIPLabel:      getLabel(builder, "session_server_ip_label"),
				SessionTechnologyLabel:    getLabel(builder, "session_technology_label"),
				SessionProtocolLabel:      getLabel(builder, "session_protocol_label"),
				SessionBytesReceivedLabel: getLabel(builder, "session_bytes_received_label"),
				SessionBytesSentLabel:     getLabel(builder, "session_bytes_sent_label"),
				SessionUptimeLabel:        getLabel(builder, "session_uptime_label"),
			},
		}

		if connectToDaemon(&app) != nil {
			app.Ui.InfoBarButton.Connect("clicked", func() { connectToDaemon(&app) })
		}

		app.Ui.CountriesComboBoxText.Connect("changed", func() { onCountrySelected(&app) })
		app.Ui.DisconnectButton.Connect("clicked", func() { onDisconnectClicked(&app) })
		app.Ui.CountryConnectButton.Connect("clicked", func() { onCountryConnect(&app) })
		app.Ui.CityConnectButton.Connect("clicked", func() { onCityConnect(&app) })
		app.Ui.GroupConnectButton.Connect("clicked", func() { onGroupConnect(&app) })
		app.Ui.ServerConnectButton.Connect("clicked", func() { onServerConnect(&app) })

		go checkSession(&app)

		app.Ui.Window.Show()
		application.AddWindow(app.Ui.Window)
	})

	os.Exit(application.Run(os.Args))
}

func getEntry(builder *gtk.Builder, name string) *gtk.Entry {
	obj, _ := builder.GetObject(name)
	return obj.(*gtk.Entry)
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

func onDisconnectClicked(app *types.Application) error {
	status, err := app.Client.Status()
	if err != nil {
		log.Printf("Error: Lost connection to NordVPN daemon: %s", err)
		app.Ui.InfoBarButton.SetLabel("Reconnect")
		app.Ui.InfoBarButton.Connect("clicked", func() { connectToDaemon(app) })
		showInfoBar("Lost connection to NordVPN daemon", gtk.MESSAGE_ERROR, app.Ui)
		return err
	}

	// The user probably used some other tool to disconnect
	if status.GetState() == "Disconnected" {
		showInfoBar("You are not connected to a VPN", gtk.MESSAGE_WARNING, app.Ui)
		return nil
	}

	err = app.Client.Disconnect()
	if err != nil {
		log.Printf("Error: Could not disconnect from VPN: %s", err)
		app.Ui.InfoBarButton.SetLabel("Reconnect")
		app.Ui.InfoBarButton.Connect("clicked", func() { connectToDaemon(app) })
		showInfoBar("Could not disconnect from VPN", gtk.MESSAGE_ERROR, app.Ui)
		return err
	}

	showInfoBar("Successfully disconnected from "+status.GetHostname(), gtk.MESSAGE_INFO, app.Ui)
	app.Ui.InfoBarButton.SetLabel("Dismiss")
	app.Ui.InfoBarButton.Connect("clicked", func() { app.Ui.InfoBar.Hide() })
	updateConnectionStatus(app)
	return nil
}

func genericConnect(app *types.Application, tag string) error {
	client, err := app.Client.Connect(&pb.ConnectRequest{
		ServerTag: tag,
		Protocol:  pb.ProtocolEnum_UDP,
		Obfuscate: false,
		CyberSec:  false,
		Dns:       nil,
		WhiteList: nil,
	})

	app.Ui.InfoBarButton.Connect("clicked", func() { app.Ui.InfoBar.Hide() })
	app.Ui.InfoBarButton.SetLabel("Dismiss")
	if err != nil {
		log.Printf("Error: Could not connect to VPN: %s", err)
		showInfoBar("Could not connect to VPN: "+err.Error(), gtk.MESSAGE_ERROR, app.Ui)
		return err
	}

	count := 0
	for {
		msg, err := client.Recv()
		count++

		if err == io.EOF {
			break
		}

		if err != nil {
			// Format common error
			if strings.Contains(err.Error(), "You are not logged in") {
				err = errors.New("You are not logged in")
			}

			log.Printf("Error: Could not connect to VPN: %s", err)
			showInfoBar("Could not connect to VPN: "+err.Error(), gtk.MESSAGE_ERROR, app.Ui)
			return err
		}

		if count == 2 {
			showInfoBar("Connected to "+msg.GetMessages()[1], gtk.MESSAGE_INFO, app.Ui)
			log.Printf("Info: Connected to %s", msg.GetMessages()[1])
		}
	}

	updateConnectionStatus(app)

	return nil
}

func onCountryConnect(app *types.Application) error {
	return genericConnect(app, app.Ui.CountriesComboBoxText.GetActiveText())
}

func onCityConnect(app *types.Application) error {
	return genericConnect(app, app.Ui.CitiesComboBoxText.GetActiveText())
}

func onServerConnect(app *types.Application) error {
	text, _ := app.Ui.ServerEntry.GetText()
	if len(text) == 0 {
		log.Println("Error: No server specified")
		app.Ui.InfoBarButton.SetLabel("Dismiss")
		app.Ui.InfoBarButton.Connect("clicked", func() { app.Ui.InfoBar.Hide() })
		showInfoBar("No server specified", gtk.MESSAGE_ERROR, app.Ui)
		return errors.New("no server specified")
	}

	return genericConnect(app, text)
}

func onGroupConnect(app *types.Application) error {
	return genericConnect(app, app.Ui.GroupsComboBoxText.GetActiveText())
}

func checkSession(app *types.Application) {
	for {
		status, err := app.Client.Status()
		if err == nil && status.GetState() == "Connected" {
			app.Ui.SessionStatusLabel.SetMarkup("<span foreground=\"green\">" + status.GetState() + "</span>")
			app.Ui.SessionServerLabel.SetText(status.GetHostname())
			app.Ui.SessionCountryLabel.SetText(status.GetCountry())
			app.Ui.SessionCityLabel.SetText(status.GetCity())
			app.Ui.SessionServerIPLabel.SetText(status.GetIp())
			app.Ui.SessionTechnologyLabel.SetText(status.GetTechnology().String())
			app.Ui.SessionProtocolLabel.SetText(status.GetProtocol().String())
			app.Ui.SessionBytesReceivedLabel.SetText(byteCountIEC(status.GetDownload()))
			app.Ui.SessionBytesSentLabel.SetText(byteCountIEC(status.GetUpload()))
			app.Ui.SessionUptimeLabel.SetText(formatDuration(status.GetUptime()))
		} else {
			app.Ui.SessionStatusLabel.SetMarkup("<span foreground=\"red\">Disconnected</span>")
		}
		time.Sleep(1 * time.Second)
	}
}

func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func formatDuration(nanoseconds int64) string {
	duration := time.Duration(nanoseconds)
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	var parts []string

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}
