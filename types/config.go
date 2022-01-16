package types

import (
	"encoding/json"
	"io/ioutil"
	"main/util"
	"os"
	"path/filepath"
)

type Config struct {
	Connect              *Connect
	AutoConnectEnabled   bool
	AutoConnectServerTag string
	CyberSecEnabled      bool
	DNSServers           []string
	FirewallEnabled      bool
	IPv6Enabled          bool
	KillSwitchEnabled    bool
	NotificationsEnabled bool
	ObfuscationEnabled   bool
	Protocol             string
	Technology           string
	WhiteList            *WhiteList
}

type Connect struct {
	Country string
	City    string
	Group   string
	Server  string
}

type WhiteList struct {
	Subnets  []string
	UDPPorts []uint32
	TCPPorts []uint32
}

func LoadConfig() *Config {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		util.LogWarning("Unable to determine user config directory", err)
		return NewConfig()
	}

	userConfigPath := filepath.Join(userConfigDir, ConfigDir, ConfigFile)
	configFile, err := os.Open(userConfigPath)
	if err != nil {
		util.LogWarning("Unable to open config file", err)
		return NewConfig()
	}

	var config Config
	bytes, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(bytes, &config)
	util.LogInfo("Loaded user config file")

	configFile.Close()

	return &config
}

func NewConfig() *Config {
	return &Config{
		Connect: &Connect{
			Country: "",
			City:    "",
			Group:   "",
			Server:  "",
		},
		AutoConnectEnabled:   false,
		AutoConnectServerTag: "",
		CyberSecEnabled:      false,
		DNSServers:           nil,
		FirewallEnabled:      false,
		IPv6Enabled:          false,
		KillSwitchEnabled:    false,
		NotificationsEnabled: false,
		ObfuscationEnabled:   false,
		Protocol:             "",
		Technology:           "",
		WhiteList: &WhiteList{
			Subnets:  []string{},
			UDPPorts: []uint32{},
			TCPPorts: []uint32{},
		},
	}
}

func SaveConfig(app *Application) error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		util.LogError("Unable to determine user config directory", err)
		return err
	}
	userConfigPath := filepath.Join(userConfigDir, ConfigDir)
	err = os.MkdirAll(userConfigPath, 0o700)
	if err != nil {
		util.LogError("Unable to create directory", err)
	}

	configFile, err := os.OpenFile(filepath.Join(userConfigPath, ConfigFile),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0700)
	if err != nil {
		util.LogError("err", err)
	}
	bytes, _ := json.MarshalIndent(app.Config, "", "  ")
	configFile.Write(bytes)
	configFile.Close()
	return nil
}
