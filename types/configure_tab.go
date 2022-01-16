package types

import (
	"github.com/adamdb5/opennord/pb"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"main/util"
	"strings"
)

type ConfigureTab struct {
	AutoConnectSwitch      *gtk.Switch
	AutoConnectServerEntry *gtk.Entry
	AutoConnectButton      *gtk.Button
	CyberSecSwitch         *gtk.Switch
	DNSEntry               *gtk.Entry
	DnsButton              *gtk.Button
	FirewallSwitch         *gtk.Switch
	IPv6Switch             *gtk.Switch
	KillSwitchSwitch       *gtk.Switch
	NotifySwitch           *gtk.Switch
	ObfuscationSwitch      *gtk.Switch
	ProtocolComboText      *gtk.ComboBoxText
	TechnologyComboText    *gtk.ComboBoxText
	SaveButton             *gtk.Button
}

func BuildConfigureTab(builder *gtk.Builder) *ConfigureTab {
	return &ConfigureTab{
		AutoConnectSwitch: util.BuilderGetSwitch(builder,
			"configure_autoconnect_switch"),
		AutoConnectServerEntry: util.BuilderGetEntry(builder,
			"configure_autoconnect_entry"),
		AutoConnectButton: util.BuilderGetButton(builder,
			"configure_autoconnect_button"),
		CyberSecSwitch: util.BuilderGetSwitch(builder,
			"configure_cybersec_switch"),
		DNSEntry:  util.BuilderGetEntry(builder, "configure_dns_entry"),
		DnsButton: util.BuilderGetButton(builder, "configure_dns_button"),
		FirewallSwitch: util.BuilderGetSwitch(builder,
			"configure_firewall_switch"),
		IPv6Switch: util.BuilderGetSwitch(builder,
			"configure_ipv6_switch"),
		KillSwitchSwitch: util.BuilderGetSwitch(builder,
			"configure_kill_switch_switch"),
		NotifySwitch: util.BuilderGetSwitch(builder,
			"configure_notify_switch"),
		ObfuscationSwitch: util.BuilderGetSwitch(builder,
			"configure_obfuscation_switch"),
		ProtocolComboText: util.BuilderGetComboBoxText(builder,
			"configure_protocol_combo_text"),
		TechnologyComboText: util.BuilderGetComboBoxText(builder,
			"configure_technology_combo_text"),
		SaveButton: util.BuilderGetButton(builder, "configure_save_button"),
	}
}

func AutoConnectClicked(app *Application) error {
	configureTab := app.Window.ConfigureTab
	serverTag, _ := configureTab.AutoConnectServerEntry.GetText()
	protocol := pb.ProtocolEnum_UDP
	dnsText, _ := configureTab.DNSEntry.GetText()
	dns := strings.Split(dnsText, ",")

	if configureTab.ProtocolComboText.GetActiveText() == "TCP" {
		protocol = pb.ProtocolEnum_TCP
	}

	_, err := app.Client.SetAutoConnect(&pb.SetAutoConnectRequest{
		ServerTag:   serverTag,
		Protocol:    protocol,
		CyberSec:    configureTab.CyberSecSwitch.GetActive(),
		Obfuscate:   configureTab.ObfuscationSwitch.GetActive(),
		AutoConnect: configureTab.AutoConnectSwitch.GetActive(),
		Dns:         dns,
		Whitelist:   nil,
	})

	if err != nil {
		util.LogError("Unable to set Auto-connect configuration", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to set Auto-connect configuration: "+
			err.Error(), gtk.MESSAGE_ERROR)
		return err
	}

	return nil
}

func ConfigureSaveClicked(app *Application) error {
	configureTab := app.Window.ConfigureTab
	serverTag, _ := configureTab.AutoConnectServerEntry.GetText()
	dnsText, _ := configureTab.DNSEntry.GetText()
	dns := strings.Split(dnsText, ",")

	app.Config.AutoConnectEnabled = configureTab.AutoConnectSwitch.GetActive()
	app.Config.AutoConnectServerTag = serverTag
	app.Config.DNSServers = dns
	app.Config.CyberSecEnabled = configureTab.CyberSecSwitch.GetActive()
	app.Config.FirewallEnabled = configureTab.FirewallSwitch.GetActive()
	app.Config.IPv6Enabled = configureTab.IPv6Switch.GetActive()
	app.Config.KillSwitchEnabled = configureTab.KillSwitchSwitch.GetActive()
	app.Config.NotificationsEnabled = configureTab.NotifySwitch.GetActive()
	app.Config.ObfuscationEnabled = configureTab.ObfuscationSwitch.GetActive()
	app.Config.Protocol = configureTab.ProtocolComboText.GetActiveText()
	app.Config.Technology = configureTab.TechnologyComboText.GetActiveText()

	return SaveConfig(app)
}

func DNSButtonClicked(app *Application) error {
	configureTab := app.Window.ConfigureTab
	dnsText, _ := configureTab.DNSEntry.GetText()
	dns := strings.Split(dnsText, ",")

	err := app.Client.SetDns(&pb.SetDNSRequest{
		Dns:      dns,
		CyberSec: true,
	})

	log.Printf("here")

	if err != nil {
		util.LogError("Unable to set DNS", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to set DNS: "+
			err.Error(), gtk.MESSAGE_ERROR)
		return err
	}

	return nil
}

func CyberSecSwitchToggled(app *Application) error {
	err := app.Client.SetCyberSec(app.Window.ConfigureTab.CyberSecSwitch.
		GetActive())

	if err != nil {
		util.LogError("Unable to set CyberSec", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to set CyberSec: "+
			err.Error(), gtk.MESSAGE_ERROR)
		return err
	}

	app.Config.CyberSecEnabled = app.Window.ConfigureTab.CyberSecSwitch.
		GetActive()

	return nil
}

func FirewallSwitchToggled(app *Application) error {
	err := app.Client.SetFirewall(app.Window.ConfigureTab.FirewallSwitch.
		GetActive())

	if err != nil {
		util.LogError("Unable to set Firewall", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to set Firewall: "+
			err.Error(), gtk.MESSAGE_ERROR)
		return err
	}

	app.Config.FirewallEnabled = app.Window.ConfigureTab.FirewallSwitch.
		GetActive()

	return nil
}

func IPv6SwitchToggled(app *Application) error {
	err := app.Client.SetIpv6(app.Window.ConfigureTab.IPv6Switch.
		GetActive())

	if err != nil {
		util.LogError("Unable to set IPv6", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to set IPv6: "+
			err.Error(), gtk.MESSAGE_ERROR)
		return err
	}

	app.Config.IPv6Enabled = app.Window.ConfigureTab.IPv6Switch.
		GetActive()

	return nil
}

func KillSwitchSwitchToggled(app *Application) error {
	err := app.Client.SetKillSwitch(app.Window.ConfigureTab.KillSwitchSwitch.
		GetActive())

	if err != nil {
		util.LogError("Unable to set Kill Switch", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to set Kill Switch: "+
			err.Error(), gtk.MESSAGE_ERROR)
		return err
	}

	app.Config.KillSwitchEnabled = app.Window.ConfigureTab.KillSwitchSwitch.
		GetActive()

	return nil
}

func NotificationsSwitchToggled(app *Application) error {
	err := app.Client.SetNotify(app.Window.ConfigureTab.NotifySwitch.
		GetActive())

	if err != nil {
		util.LogError("Unable to set Notifications", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to set Notifications: "+
			err.Error(), gtk.MESSAGE_ERROR)
		return err
	}

	app.Config.NotificationsEnabled = app.Window.ConfigureTab.NotifySwitch.
		GetActive()

	return nil
}

func ObfuscationSwitchToggled(app *Application) error {
	err := app.Client.SetObfuscate(app.Window.ConfigureTab.ObfuscationSwitch.
		GetActive())

	if err != nil {
		util.LogError("Unable to set Obfuscation", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to set Obfuscation: "+
			err.Error(), gtk.MESSAGE_ERROR)
		return err
	}

	app.Config.ObfuscationEnabled = app.Window.ConfigureTab.ObfuscationSwitch.
		GetActive()

	return nil
}

func ProtocolComboTextChanged(app *Application) error {
	configureTab := app.Window.ConfigureTab
	protocol := pb.ProtocolEnum_UDP

	if configureTab.ProtocolComboText.GetActiveText() == "TCP" {
		protocol = pb.ProtocolEnum_TCP
	}

	err := app.Client.SetProtocol(protocol)
	if err != nil {
		util.LogError("Unable to set Protocol", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to set Protocol: "+
			err.Error(), gtk.MESSAGE_ERROR)
		return err
	}

	app.Config.Protocol = configureTab.ProtocolComboText.GetActiveText()

	return nil
}

func TechnologyComboTextChanged(app *Application) error {
	configureTab := app.Window.ConfigureTab
	technology := pb.TechnologyEnum_OPENVPN

	if configureTab.TechnologyComboText.GetActiveText() == "NORDLYNX" {
		technology = pb.TechnologyEnum_NORDLYNX
	}

	err := app.Client.SetTechnology(technology)
	if err != nil {
		util.LogError("Unable to set Technology", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to set Technology: "+
			err.Error(), gtk.MESSAGE_ERROR)
		return err
	}

	app.Config.Technology = configureTab.TechnologyComboText.GetActiveText()

	return nil
}
