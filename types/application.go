package types

import (
	"errors"
	"github.com/adamdb5/opennord"
	"github.com/adamdb5/opennord/pb"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"io"
	"main/util"
	"strconv"
	"strings"
	"time"
)

// Application contains references to the OpenNord client and all functional
// GTK Controls.
type Application struct {
	Client *opennord.Client
	Window *Window
	Config *Config
}

// BuildApplication instantiates the Application and registers the GTK
// components.
func BuildApplication(builder *gtk.Builder) *Application {
	window := BuildWindow(builder)
	app := &Application{
		Client: nil,
		Window: window,
		Config: LoadConfig(),
	}

	return app
}

// RegisterCallbacks connects the GUI controls to the corresponding callbacks.
func (app Application) RegisterCallbacks() {
	// Connect
	app.Window.ConnectTab.DisconnectButton.Connect("clicked",
		func() { _ = DisconnectClicked(&app) })
	app.Window.ConnectTab.CountriesComboBoxText.Connect("changed",
		func() { CountrySelected(&app) })
	app.Window.ConnectTab.CountryConnectButton.Connect("clicked",
		func() { _ = ConnectToCountry(&app) })
	app.Window.ConnectTab.CityConnectButton.Connect("clicked",
		func() { _ = ConnectToCity(&app) })
	app.Window.ConnectTab.GroupConnectButton.Connect("clicked",
		func() { _ = ConnectToGroup(&app) })
	app.Window.ConnectTab.ServerConnectButton.Connect("clicked",
		func() { _ = ConnectToServer(&app) })
	app.Window.ConnectTab.BestConnectButton.Connect("clicked",
		func() { _ = app.Connect("") })
	app.Window.ConnectTab.SaveButton.Connect("clicked",
		func() { _ = ConnectSaveClicked(&app) })

	// Configure
	app.Window.ConfigureTab.AutoConnectButton.Connect("clicked",
		func() { _ = AutoConnectClicked(&app) })
	app.Window.ConfigureTab.CyberSecSwitch.Connect("state-set",
		func() { _ = CyberSecSwitchToggled(&app) })
	app.Window.ConfigureTab.DnsButton.Connect("clicked",
		func() { _ = DNSButtonClicked(&app) })
	app.Window.ConfigureTab.FirewallSwitch.Connect("state-set",
		func() { _ = FirewallSwitchToggled(&app) })
	app.Window.ConfigureTab.IPv6Switch.Connect("state-set",
		func() { _ = IPv6SwitchToggled(&app) })
	app.Window.ConfigureTab.KillSwitchSwitch.Connect("state-set",
		func() { _ = KillSwitchSwitchToggled(&app) })
	app.Window.ConfigureTab.NotifySwitch.Connect("state-set",
		func() { _ = NotificationsSwitchToggled(&app) })
	app.Window.ConfigureTab.ObfuscationSwitch.Connect("state-set",
		func() { _ = ObfuscationSwitchToggled(&app) })
	app.Window.ConfigureTab.ProtocolComboText.Connect("changed",
		func() { _ = ProtocolComboTextChanged(&app) })
	app.Window.ConfigureTab.TechnologyComboText.Connect("changed",
		func() { _ = TechnologyComboTextChanged(&app) })

	// Whitelist
	app.Window.WhiteListTab.SubnetAddButton.Connect("clicked",
		func() { _ = SubnetAddButtonClicked(&app) })
	app.Window.WhiteListTab.SubnetRemoveButton.Connect("clicked",
		func() { _ = SubnetRemoveButtonClicked(&app) })
	app.Window.WhiteListTab.UDPAddButton.Connect("clicked",
		func() { _ = UDPAddButtonClicked(&app) })
	app.Window.WhiteListTab.UDPRemoveButton.Connect("clicked",
		func() { _ = UDPRemoveButtonClicked(&app) })
	app.Window.WhiteListTab.TCPAddButton.Connect("clicked",
		func() { _ = TCPAddButtonClicked(&app) })
	app.Window.WhiteListTab.TCPRemoveButton.Connect("clicked",
		func() { _ = TCPRemoveButtonClicked(&app) })

	// Account
	app.Window.AccountTab.RefreshButton.Connect("clicked",
		func() { _ = AccountRefreshClicked(&app) })
	app.Window.AccountTab.OAuthButton.Connect("clicked",
		func() { _ = GenerateOAuthClicked(&app) })
	app.Window.AccountTab.LoginButton.Connect("clicked",
		func() { _ = LoginClicked(&app) })
	app.Window.AccountTab.LogoutButton.Connect("clicked",
		func() { _ = LogoutClicked(&app) })
}

// PopulateCountries makes a request via the client to retrieve all countries
// supported by the daemon.
func (app Application) PopulateCountries() error {
	countries, err := app.Client.Countries()

	if err != nil {
		util.LogError("Unable to retrieve countries", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to retrieve countries",
			gtk.MESSAGE_ERROR)
		return err
	}

	list, _ := gtk.ListStoreNew(glib.TYPE_STRING)
	for _, country := range countries.GetCountries() {
		appendIter := list.Append()
		_ = list.SetValue(appendIter, 0, country)
	}

	renderer, _ := gtk.CellRendererTextNew()
	connectTab := app.Window.ConnectTab
	connectTab.CountriesComboBoxText.SetModel(list)
	connectTab.CountriesComboBoxText.PackStart(renderer, true)
	connectTab.CountriesComboBoxText.SetActive(0)

	return nil
}

// PopulateCities makes a request via the client to retrieve all cities
// supported by the daemon.
func (app Application) PopulateCities() error {
	connectTab := app.Window.ConnectTab

	cities, err := app.Client.Cities(
		connectTab.CountriesComboBoxText.GetActiveText())

	if err != nil {
		util.LogError("Unable to retrieve cities", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to retrieve cities",
			gtk.MESSAGE_ERROR)
		return err
	}

	list, _ := gtk.ListStoreNew(glib.TYPE_STRING)
	for _, city := range cities.GetCities() {
		appendIter := list.Append()
		_ = list.SetValue(appendIter, 0, city)
	}

	renderer, _ := gtk.CellRendererTextNew()
	connectTab.CitiesComboBoxText.SetModel(list)
	connectTab.CitiesComboBoxText.PackStart(renderer, true)
	connectTab.CitiesComboBoxText.SetActive(0)

	return nil
}

// PopulateGroups makes a request via the client to retrieve all groups
// supported by the daemon.
func (app Application) PopulateGroups() error {
	groups, err := app.Client.Groups(&pb.GroupsRequest{
		Protocol:  pb.ProtocolEnum_UDP,
		Obfuscate: false,
	})

	if err != nil {
		util.LogError("Unable to retrieve groups", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to retrieve groups",
			gtk.MESSAGE_ERROR)
		return err
	}

	list, _ := gtk.ListStoreNew(glib.TYPE_STRING)
	for _, group := range groups.GetGroups() {
		appendIter := list.Append()
		_ = list.SetValue(appendIter, 0, group)
	}

	renderer, _ := gtk.CellRendererTextNew()
	connectTab := app.Window.ConnectTab
	connectTab.GroupsComboBoxText.SetModel(list)
	connectTab.GroupsComboBoxText.PackStart(renderer, true)
	connectTab.GroupsComboBoxText.SetActive(0)

	return nil
}

// PopulateProtocols makes a request via the client to retrieve all protocols
// supported by the daemon.
func (app Application) PopulateProtocols() error {
	protocols, err := app.Client.SettingsProtocols()

	if err != nil {
		util.LogError("Unable to retrieve protocols", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to retrieve protocols",
			gtk.MESSAGE_ERROR)
		return err
	}

	list, _ := gtk.ListStoreNew(glib.TYPE_STRING)
	for _, group := range protocols.GetProtocols() {
		appendIter := list.Append()
		_ = list.SetValue(appendIter, 0, group)
	}

	renderer, _ := gtk.CellRendererTextNew()
	configureTab := app.Window.ConfigureTab
	configureTab.ProtocolComboText.SetModel(list)
	configureTab.ProtocolComboText.PackStart(renderer, true)
	configureTab.ProtocolComboText.SetActive(0)

	return nil
}

// PopulateTechnologies makes a request via the client to retrieve all
// technologies supported by the daemon.
func (app Application) PopulateTechnologies() error {
	technologies, err := app.Client.SettingsTechnologies()

	if err != nil {
		util.LogError("Unable to retrieve technologies", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to retrieve technologies",
			gtk.MESSAGE_ERROR)
		return err
	}

	list, _ := gtk.ListStoreNew(glib.TYPE_STRING)
	for _, group := range technologies.GetTechnologies() {
		appendIter := list.Append()
		_ = list.SetValue(appendIter, 0, group)
	}

	renderer, _ := gtk.CellRendererTextNew()
	configureTab := app.Window.ConfigureTab
	configureTab.TechnologyComboText.SetModel(list)
	configureTab.TechnologyComboText.PackStart(renderer, true)
	configureTab.TechnologyComboText.SetActive(0)

	return nil
}

// ConnectToDaemon attempts to connect to the NordVPN daemon. If the connection
// is successful, the connection status and account information will be updated.
// Additionally, the countries, cities and groups on the 'Connect' tab will be
// populated.
func (app *Application) ConnectToDaemon() error {
	infoBar := app.Window.InfoBar
	infoBar.HideMessage()
	client, err := opennord.NewOpenNordClient()
	if err != nil {
		util.LogError("Could not connect to NordVPN daemon", err)
		infoBar.DisplayMessage(
			"Could not connect to NordVPN daemon", gtk.MESSAGE_ERROR)
		return err
	}

	app.Client = &client

	// If we have a good client, we can register our callbacks
	app.RegisterCallbacks()

	// And update the GUI
	_ = app.UpdateConnectionStatus()
	_ = app.UpdateAccountInformation()
	_ = app.PopulateCountries()
	_ = app.PopulateCities()
	_ = app.PopulateGroups()
	_ = app.PopulateProtocols()
	_ = app.PopulateTechnologies()
	app.PopulateFromConfig()

	return nil
}

func (app Application) PopulateFromConfig() {
	connectTab := app.Window.ConnectTab

	// Populate countries
	model, _ := connectTab.CountriesComboBoxText.GetModel()
	treeModel := model.ToTreeModel()
	cnt := 0
	treeModel.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath,
		iter *gtk.TreeIter) bool {
		val, _ := treeModel.GetValue(iter, 0)
		strVal, _ := val.GetString()
		if app.Config.Connect.Country == strVal {
			connectTab.CountriesComboBoxText.SetActive(cnt)
		}
		cnt++
		return false
	})

	// Populate city
	model, _ = connectTab.CitiesComboBoxText.GetModel()
	treeModel = model.ToTreeModel()
	cnt = 0
	treeModel.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath,
		iter *gtk.TreeIter) bool {
		val, _ := treeModel.GetValue(iter, 0)
		strVal, _ := val.GetString()
		if app.Config.Connect.City == strVal {
			connectTab.CitiesComboBoxText.SetActive(cnt)
		}
		cnt++
		return false
	})

	// Populate group
	model, _ = connectTab.GroupsComboBoxText.GetModel()
	treeModel = model.ToTreeModel()
	cnt = 0
	treeModel.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath,
		iter *gtk.TreeIter) bool {
		val, _ := treeModel.GetValue(iter, 0)
		strVal, _ := val.GetString()
		if app.Config.Connect.Group == strVal {
			connectTab.GroupsComboBoxText.SetActive(cnt)
		}
		cnt++
		return false
	})

	// Populate server
	connectTab.ServerEntry.SetText(app.Config.Connect.Server)

	configureTab := app.Window.ConfigureTab
	// Populate Auto-connect
	configureTab.AutoConnectSwitch.SetActive(app.Config.AutoConnectEnabled)
	configureTab.AutoConnectServerEntry.SetText(app.Config.AutoConnectServerTag)
	configureTab.CyberSecSwitch.SetActive(app.Config.CyberSecEnabled)
	configureTab.DNSEntry.SetText(strings.Join(app.Config.DNSServers, ","))
	configureTab.FirewallSwitch.SetActive(app.Config.FirewallEnabled)
	configureTab.IPv6Switch.SetActive(app.Config.IPv6Enabled)
	configureTab.KillSwitchSwitch.SetActive(app.Config.KillSwitchEnabled)
	configureTab.NotifySwitch.SetActive(app.Config.NotificationsEnabled)
	configureTab.ObfuscationSwitch.SetActive(app.Config.ObfuscationEnabled)

	model, _ = configureTab.ProtocolComboText.GetModel()
	treeModel = model.ToTreeModel()
	cnt = 0
	treeModel.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath,
		iter *gtk.TreeIter) bool {
		val, _ := treeModel.GetValue(iter, 0)
		strVal, _ := val.GetString()
		if app.Config.Protocol == strVal {
			configureTab.ProtocolComboText.SetActive(cnt)
		}
		cnt++
		return false
	})

	model, _ = configureTab.TechnologyComboText.GetModel()
	treeModel = model.ToTreeModel()
	cnt = 0
	treeModel.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath,
		iter *gtk.TreeIter) bool {
		val, _ := treeModel.GetValue(iter, 0)
		strVal, _ := val.GetString()
		if app.Config.Technology == strVal {
			configureTab.TechnologyComboText.SetActive(cnt)
		}
		cnt++
		return false
	})

	// Populate whitelist subnets
	for _, subnet := range app.Config.WhiteList.Subnets {
		label, _ := gtk.LabelNew(subnet)
		row, _ := gtk.ListBoxRowNew()
		row.SetHAlign(gtk.ALIGN_START)
		row.Add(label)
		app.Window.WhiteListTab.SubnetListBox.Add(row)
		row.ShowAll()
	}

	// Populate whitelist UDP ports
	for _, port := range app.Config.WhiteList.UDPPorts {
		label, _ := gtk.LabelNew(strconv.FormatInt(int64(port), 10))
		row, _ := gtk.ListBoxRowNew()
		row.SetHAlign(gtk.ALIGN_START)
		row.Add(label)
		app.Window.WhiteListTab.UDPListBox.Add(row)
		row.ShowAll()
	}

	// Populate whitelist TCP ports
	for _, port := range app.Config.WhiteList.TCPPorts {
		label, _ := gtk.LabelNew(strconv.FormatInt(int64(port), 10))
		row, _ := gtk.ListBoxRowNew()
		row.SetHAlign(gtk.ALIGN_START)
		row.Add(label)
		app.Window.WhiteListTab.TCPListBox.Add(row)
		row.ShowAll()
	}
}

// UpdateConnectionStatus updates the 'Connect' tab status and enables /
// disables the 'Disconnect' button depending on the connection status.
func (app Application) UpdateConnectionStatus() error {
	status, err := app.Client.Status()
	if err != nil {
		util.LogError("Lost connection to NordVPN daemon", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Reconnect")
		infoBar.Button.Connect("clicked", app.ConnectToDaemon)
		infoBar.DisplayMessage(
			"Lost connection to NordVPN daemon", gtk.MESSAGE_ERROR)
		return err
	}

	connectedText := status.GetState()
	canDisconnect := connectedText == "Connected"
	connectTab := app.Window.ConnectTab
	connectTab.StatusLabel.SetText(connectedText)
	connectTab.DisconnectButton.SetSensitive(canDisconnect)

	return nil
}

// UpdateAccountInformation updates the 'Account' tab with the user's account
// information. If the user is logged in, the login options will be disabled. If
// the user is logged out, the login options will be enabled.
func (app Application) UpdateAccountInformation() error {
	isLoggedIn, err := app.Client.IsLoggedIn()
	infoBar := app.Window.InfoBar

	if err != nil {
		util.LogError("Lost connection to NordVPN daemon", err)
		infoBar.Button.SetLabel("Reconnect")
		infoBar.Button.Connect("clicked", app.ConnectToDaemon)
		infoBar.DisplayMessage(
			"Lost connection to NordVPN daemon", gtk.MESSAGE_ERROR)
		return err
	}

	accountTab := app.Window.AccountTab
	if !isLoggedIn.GetIsLoggedIn() {
		accountTab.EmailLabel.SetText("N/A")
		accountTab.ExpiresLabel.SetText("N/A")
		accountTab.LoginFrame.SetSensitive(true)
		accountTab.OAuthFrame.SetSensitive(true)
		accountTab.LogoutButton.SetSensitive(false)
		accountTab.StatusLabel.SetText("Not Logged In")

		util.LogError("You are not logged in to NordVPN", nil)
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage(
			"You are not logged in to NordVPN", gtk.MESSAGE_ERROR)
		return err
	}

	account, _ := app.Client.AccountInfo()
	accountTab.EmailLabel.SetText(account.GetEmail())
	accountTab.ExpiresLabel.SetText(account.GetExpiresAt())
	accountTab.LoginFrame.SetSensitive(false)
	accountTab.OAuthFrame.SetSensitive(false)
	accountTab.LogoutButton.SetSensitive(true)
	accountTab.StatusLabel.SetText("Logged In")

	return nil
}

// UpdateSessionStatus updates the information on the 'Session' tab. This
// function is intended to be run as a goroutine and will update the
// information once every second.
func (app Application) UpdateSessionStatus() {
	for {
		if app.Client != nil {
			status, err := app.Client.Status()
			sessionTab := app.Window.SessionTab

			if err == nil && status.GetState() == "Connected" {
				sessionTab.StatusLabel.SetText(status.GetState())
				sessionTab.ServerLabel.SetText(status.GetHostname())
				sessionTab.CountryLabel.SetText(status.GetCountry())
				sessionTab.CityLabel.SetText(status.GetCity())
				sessionTab.ServerIPLabel.SetText(status.GetIp())
				sessionTab.TechnologyLabel.SetText(status.GetTechnology().
					String())
				sessionTab.ProtocolLabel.SetText(status.GetProtocol().String())
				sessionTab.BytesReceivedLabel.SetText(util.FormatBytes(status.
					GetDownload()))
				sessionTab.BytesSentLabel.SetText(util.FormatBytes(status.
					GetUpload()))
				sessionTab.UptimeLabel.SetText(util.FormatDuration(status.
					GetUptime()))
			} else {
				sessionTab.StatusLabel.SetText("Disconnected")
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// Connect connects to the server specified by the given tag.
func (app Application) Connect(tag string) error {
	infoBar := app.Window.InfoBar

	if app.Client == nil {
		errMsg := "you are not connected to the NordVPN daemon"
		util.LogError(errMsg, nil)
		infoBar.DisplayMessage(errMsg, gtk.MESSAGE_ERROR)
		infoBar.Button.SetLabel("Retry")
		infoBar.Button.Connect("clicked", app.ConnectToDaemon)

		return errors.New(errMsg)
	}

	client, err := app.Client.Connect(&pb.ConnectRequest{
		ServerTag: tag,
		Protocol:  pb.ProtocolEnum_UDP,
		Obfuscate: false,
		CyberSec:  app.Config.CyberSecEnabled,
		Dns:       nil,
		WhiteList: nil,
	})

	infoBar.Button.Connect("clicked", infoBar.HideMessage)
	infoBar.Button.SetLabel("Dismiss")
	if err != nil {
		util.LogError("Lost connection to NordVPN daemon", err)
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Lost connection to NordVPN daemon: "+err.Error(),
			gtk.MESSAGE_ERROR)
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
				err = errors.New("you are not logged in")
			} else {
				util.LogError("Lost connection to NordVPN daemon", err)
				infoBar.Button.SetLabel("Dismiss")
				infoBar.Button.Connect("clicked", infoBar.HideMessage)
				infoBar.DisplayMessage(
					"Lost connection to NordVPN daemon: "+err.Error(),
					gtk.MESSAGE_ERROR)
			}
			return err
		}

		if count == 2 {
			infoBar.DisplayMessage("Connected to "+msg.GetMessages()[1],
				gtk.MESSAGE_INFO)
			util.LogInfo("Connected to " + msg.GetMessages()[1])
		}
	}

	_ = app.UpdateConnectionStatus()

	return nil
}
