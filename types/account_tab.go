package types

import (
	"github.com/adamdb5/opennord/pb"
	"github.com/gotk3/gotk3/gtk"
	"main/util"
	"os/exec"
)

// AccountTab contains the GTK components for the 'Account' GTKNotebook page.
type AccountTab struct {
	StatusLabel       *gtk.Label
	EmailLabel        *gtk.Label
	ExpiresLabel      *gtk.Label
	LogoutButton      *gtk.Button
	RefreshButton     *gtk.Button
	EmailEntry        *gtk.Entry
	PasswordEntry     *gtk.Entry
	LoginButton       *gtk.Button
	OAuthButton       *gtk.Button
	OAuthURLEntry     *gtk.Entry
	OpenBrowserButton *gtk.Button
	LoginFrame        *gtk.Frame
	OAuthFrame        *gtk.Frame
}

// BuildAccountTab constructs the GTKNotebook page for the 'Account' tab from
// the provided builder.
func BuildAccountTab(builder *gtk.Builder) *AccountTab {
	return &AccountTab{
		StatusLabel: util.BuilderGetLabel(builder,
			"account_status_label"),
		EmailLabel: util.BuilderGetLabel(builder,
			"account_email_label"),
		ExpiresLabel: util.BuilderGetLabel(builder,
			"account_expires_label"),
		LogoutButton: util.BuilderGetButton(builder,
			"account_logout_button"),
		RefreshButton: util.BuilderGetButton(builder,
			"account_refresh_button"),
		EmailEntry: util.BuilderGetEntry(builder,
			"account_email_entry"),
		PasswordEntry: util.BuilderGetEntry(builder,
			"account_password_entry"),
		LoginButton: util.BuilderGetButton(builder,
			"account_login_button"),
		OAuthButton: util.BuilderGetButton(builder,
			"account_oauth_button"),
		OAuthURLEntry: util.BuilderGetEntry(builder,
			"account_oauth_url_entry"),
		OpenBrowserButton: util.BuilderGetButton(builder,
			"account_open_browser_button"),
		LoginFrame: util.BuilderGetFrame(builder,
			"account_login_frame"),
		OAuthFrame: util.BuilderGetFrame(builder,
			"account_oauth_frame"),
	}
}

// AccountRefreshClicked is invoked whenever the 'Refresh' button on the
// 'Account' tab is clicked.
func AccountRefreshClicked(app *Application) error {
	isLoggedIn, err := app.Client.IsLoggedIn()

	if err != nil {
		util.LogError("Lost connection to NordVPN daemon", err)
		infoBar := app.Window.InfoBar
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage(
			"Lost connection to NordVPN daemon", gtk.MESSAGE_ERROR)
		return err
	}

	if isLoggedIn.GetIsLoggedIn() {
		_ = app.UpdateAccountInformation()
	}

	return nil
}

// GenerateOAuthClicked is invoked whenever the 'Generate OAuth Token' button
// on the 'Account' tab is clicked.
func GenerateOAuthClicked(app *Application) error {
	isLoggedIn, err := app.Client.IsLoggedIn()
	infoBar := app.Window.InfoBar

	if err != nil {
		util.LogError("Lost connection to NordVPN daemon", err)
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Lost connection to NordVPN daemon",
			gtk.MESSAGE_ERROR)
		return err
	}

	// This should only happen if the user already logged in using the CLI
	if isLoggedIn.GetIsLoggedIn() {
		_ = app.UpdateAccountInformation()
		return nil
	}

	oauth, err := app.Client.LoginOAuth2()
	if err != nil {
		util.LogError("Unable to get OAuth token", err)
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to get OAuth token", gtk.MESSAGE_ERROR)
		return err
	}

	accountTab := app.Window.AccountTab
	accountTab.OAuthURLEntry.SetText(oauth.GetUrl())
	accountTab.OpenBrowserButton.Connect("clicked", func() {
		err = exec.Command("xdg-open", oauth.GetUrl()).Start()
		if err != nil {
			util.LogError("Unable to open URL", err)
			infoBar.Button.SetLabel("Dismiss")
			infoBar.Button.Connect("clicked", infoBar.HideMessage)
			infoBar.DisplayMessage("Unable to open URL", gtk.MESSAGE_ERROR)
		}
	})

	return nil
}

// LoginClicked is invoked whenever the 'Login' button on the 'Account' tab is
// clicked.
func LoginClicked(app *Application) error {
	isLoggedIn, err := app.Client.IsLoggedIn()
	infoBar := app.Window.InfoBar

	if err != nil {
		util.LogError("Lost connection to NordVPN daemon", err)
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Lost connection to NordVPN daemon",
			gtk.MESSAGE_ERROR)
		return err
	}

	// This should only happen if the user already logged in using the CLI
	if isLoggedIn.GetIsLoggedIn() {
		_ = app.UpdateAccountInformation()
		return nil
	}

	username, _ := app.Window.AccountTab.EmailEntry.GetText()
	password, _ := app.Window.AccountTab.PasswordEntry.GetText()
	err = app.Client.Login(&pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		util.LogError("Unable to log in", err)
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to log in", gtk.MESSAGE_ERROR)
		return err
	}

	_ = app.UpdateAccountInformation()
	util.LogInfo("Logged in using email" + username)
	infoBar.Button.SetLabel("Dismiss")
	infoBar.Button.Connect("clicked", infoBar.HideMessage)
	infoBar.DisplayMessage("Logged in using email "+username, gtk.MESSAGE_INFO)
	return nil
}

// LogoutClicked is invoked whenever the 'Logout' button on the 'Account' tab is
// clicked.
func LogoutClicked(app *Application) error {
	isLoggedIn, err := app.Client.IsLoggedIn()
	infoBar := app.Window.InfoBar

	if err != nil {
		util.LogError("Lost connection to NordVPN", err)
		infoBar.Button.SetLabel("Reconnect")
		infoBar.Button.Connect("clicked", func() { _ = app.ConnectToDaemon() })
		infoBar.DisplayMessage("Lost connection to NordVPN daemon",
			gtk.MESSAGE_ERROR)
		return err
	}

	// This should only happen if the user already logged out using the CLI
	if !isLoggedIn.GetIsLoggedIn() {
		_ = app.UpdateAccountInformation()
		return nil
	}

	err = app.Client.Logout()
	if err != nil {
		util.LogError("Unable to log out", err)
		infoBar.Button.SetLabel("Dismiss")
		infoBar.Button.Connect("clicked", infoBar.HideMessage)
		infoBar.DisplayMessage("Unable to log out", gtk.MESSAGE_ERROR)
		return err
	}

	_ = app.UpdateAccountInformation()
	util.LogInfo("Logged out")
	infoBar.Button.SetLabel("Dismiss")
	infoBar.Button.Connect("clicked", infoBar.HideMessage)
	infoBar.DisplayMessage("Successfully logged out.", gtk.MESSAGE_INFO)
	return nil
}
