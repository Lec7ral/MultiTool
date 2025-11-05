package profiles

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Profile defines the structure for a network configuration profile.
type Profile struct {
	Name            string `json:"name"`
	NetworkPriority string `json:"networkPriority"` // "Ethernet" or "Wi-Fi"
	ProxyEnabled    bool   `json:"proxyEnabled"`
	ProxyServer     string `json:"proxyServer"`
}

var profilesFilePath string

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback if config dir is not found
		configDir = "."
	}
	appConfigDir := filepath.Join(configDir, "MultiTool")
	os.MkdirAll(appConfigDir, os.ModePerm)
	profilesFilePath = filepath.Join(appConfigDir, "profiles.json")
}

// LoadProfiles reads the profiles from the config file.
// If the file doesn't exist, it creates default profiles.
func LoadProfiles() ([]Profile, error) {
	data, err := os.ReadFile(profilesFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create default profiles and save
			defaultProfiles := []Profile{
				{Name: "Wired", NetworkPriority: "Ethernet", ProxyEnabled: true, ProxyServer: "10.0.0.1:8080"},
				{Name: "Wi-Fi", NetworkPriority: "Wi-Fi", ProxyEnabled: false, ProxyServer: ""},
			}
			if err := SaveProfiles(defaultProfiles); err != nil {
				return nil, err
			}
			return defaultProfiles, nil
		}
		return nil, err
	}

	var profiles []Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

// SaveProfiles writes the given profiles to the config file.
func SaveProfiles(profiles []Profile) error {
	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(profilesFilePath, data, 0644)
}
