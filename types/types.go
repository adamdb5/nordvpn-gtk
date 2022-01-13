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
	// Root
	Window *gtk.Window

	// Global
	InfoBar       *gtk.InfoBar
	InfoBarLabel  *gtk.Label
	InfoBarButton *gtk.Button

	// Connect Tab
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

	// Session Tab
	SessionStatusLabel        *gtk.Label
	SessionServerLabel        *gtk.Label
	SessionCountryLabel       *gtk.Label
	SessionCityLabel          *gtk.Label
	SessionServerIPLabel      *gtk.Label
	SessionTechnologyLabel    *gtk.Label
	SessionProtocolLabel      *gtk.Label
	SessionBytesReceivedLabel *gtk.Label
	SessionBytesSentLabel     *gtk.Label
	SessionUptimeLabel        *gtk.Label

	// Account Tab
	AccountStatusLabel       *gtk.Label
	AccountEmailLabel        *gtk.Label
	AccountExpiresLabel      *gtk.Label
	AccountLogoutButton      *gtk.Button
	AccountRefreshButton     *gtk.Button
	AccountEmailEntry        *gtk.Entry
	AccountPasswordEntry     *gtk.Entry
	AccountLoginButton       *gtk.Button
	AccountOAuthButton       *gtk.Button
	AccountOAuthURLEntry     *gtk.Entry
	AccountOpenBrowserButton *gtk.Button
	AccountLoginFrame        *gtk.Frame
	AccountOAuthFrame        *gtk.Frame

	// About Tab
	AboutNameLabel        *gtk.Label
	AboutVersionLabel     *gtk.Label
	AboutDescriptionLabel *gtk.Label
	AboutWebsiteLabel     *gtk.Label
	AboutCopyrightLabel   *gtk.Label
	AboutLicenseLabel     *gtk.Label
}
