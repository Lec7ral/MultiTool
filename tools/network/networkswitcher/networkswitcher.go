package networkswitcher

import (
	"fmt"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Lec7ral/MultiTool/tools/profiles"
)

// --- Callback para actualizar la bandeja del sistema ---
var systrayCallback func()

func SetSystrayCallback(callback func()) {
	systrayCallback = callback
}

// --- Tool Definition ---
type NetworkSwitcherTool struct{}

func New() *NetworkSwitcherTool {
	return &NetworkSwitcherTool{}
}

func (t *NetworkSwitcherTool) GetName() string {
	return "Network Switcher"
}

func (t *NetworkSwitcherTool) GetDescription() string {
	return "Manage and apply network configuration profiles"
}

func (t *NetworkSwitcherTool) GetCategory() string {
	return "Network"
}

func (t *NetworkSwitcherTool) GetIcon() fyne.Resource {
	return nil
}

// --- Main UI ---
func (t *NetworkSwitcherTool) GetUI() fyne.CanvasObject {
	statusLabel := widget.NewLabel("")
	statusLabel.Wrapping = fyne.TextWrapWord

	// --- Profile Selection ---
	loadedProfiles, err := profiles.LoadProfiles()
	if err != nil {
		return widget.NewLabel("Error loading profiles: " + err.Error())
	}

	var selectedProfile profiles.Profile
	profileNames := func() []string {
		names := make([]string, len(loadedProfiles))
		for i, p := range loadedProfiles {
			names[i] = p.Name
		}
		return names
	}()

	profileSelect := widget.NewSelect(profileNames, func(name string) {
		for _, p := range loadedProfiles {
			if p.Name == name {
				selectedProfile = p
				break
			}
		}
	})

	refreshAll := func() {
		loadedProfiles, _ = profiles.LoadProfiles()
		profileNames := func() []string {
			names := make([]string, len(loadedProfiles))
			for i, p := range loadedProfiles {
				names[i] = p.Name
			}
			return names
		}()
		profileSelect.Options = profileNames
		profileSelect.Refresh()

		if systrayCallback != nil {
			systrayCallback()
		}
	}

	if len(loadedProfiles) > 0 {
		profileSelect.SetSelectedIndex(0)
		selectedProfile = loadedProfiles[0]
	}

	// --- Main Buttons ---
	applyBtn := widget.NewButton("Apply Profile", func() {
		if selectedProfile.Name == "" {
			statusLabel.SetText("No profile selected.")
			return
		}
		statusLabel.SetText(fmt.Sprintf("Applying profile '%s'...", selectedProfile.Name))
		if err := ApplyProfile(selectedProfile); err != nil {
			statusLabel.SetText(fmt.Sprintf("Failed to apply profile: %s", err.Error()))
		} else {
			statusLabel.SetText(fmt.Sprintf("Profile '%s' applied successfully.", selectedProfile.Name))
		}
	})

	manageBtn := widget.NewButton("Manage Profiles", func() {
		parentWindow := fyne.CurrentApp().Driver().AllWindows()[0]
		newManagerWindow(parentWindow, refreshAll).Show()
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Select a Profile", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		profileSelect,
		applyBtn,
		widget.NewSeparator(),
		manageBtn,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Status", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		statusLabel,
		widget.NewSeparator(),
		widget.NewLabel("NOTE: This tool requires the application to be run with Administrator privileges."),
	)
}

// --- Profile Manager Window ---
func newManagerWindow(parent fyne.Window, onClosed func()) fyne.Window {
	app := fyne.CurrentApp()
	w := app.NewWindow("Profile Manager")
	w.Resize(fyne.NewSize(600, 400))
	w.CenterOnScreen()

	// --- Data & Form Widgets ---
	loadedProfiles, _ := profiles.LoadProfiles()
	var selectedProfile *profiles.Profile

	nameEntry := widget.NewEntry()
	prioritySelect := widget.NewSelect([]string{"Ethernet", "Wi-Fi"}, nil)
	proxyEnabledCheck := widget.NewCheck("Proxy Enabled", nil)
	proxyServerEntry := widget.NewEntry()

	form := widget.NewForm(
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Priority", prioritySelect),
		widget.NewFormItem("", proxyEnabledCheck),
		widget.NewFormItem("Proxy Server", proxyServerEntry),
	)

	// --- Profile List ---
	profileList := widget.NewList(
		func() int { return len(loadedProfiles) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText(loadedProfiles[i].Name) },
	)

	profileList.OnSelected = func(id widget.ListItemID) {
		selectedProfile = &loadedProfiles[id]
		nameEntry.SetText(selectedProfile.Name)
		prioritySelect.SetSelected(selectedProfile.NetworkPriority)
		proxyEnabledCheck.SetChecked(selectedProfile.ProxyEnabled)
		proxyServerEntry.SetText(selectedProfile.ProxyServer)
	}

	// --- Toolbar Buttons ---
	newBtn := widget.NewButton("New", func() {
		selectedProfile = nil
		profileList.UnselectAll()
		nameEntry.SetText("")
		prioritySelect.ClearSelected()
		proxyEnabledCheck.SetChecked(false)
		proxyServerEntry.SetText("")
	})

	deleteBtn := widget.NewButton("Delete", func() {
		if selectedProfile == nil { return }
		var newProfiles []profiles.Profile
		for _, p := range loadedProfiles {
			if p.Name != selectedProfile.Name {
				newProfiles = append(newProfiles, p)
			}
		}
		loadedProfiles = newProfiles
		profiles.SaveProfiles(loadedProfiles)
		profileList.Refresh()
		newBtn.OnTapped()
	})

	saveBtn := widget.NewButton("Save", func() {
		if selectedProfile != nil { // Update existing
			selectedProfile.Name = nameEntry.Text
			selectedProfile.NetworkPriority = prioritySelect.Selected
			selectedProfile.ProxyEnabled = proxyEnabledCheck.Checked
			selectedProfile.ProxyServer = proxyServerEntry.Text
		} else { // Create new
			newProfile := profiles.Profile{
				Name:            nameEntry.Text,
				NetworkPriority: prioritySelect.Selected,
				ProxyEnabled:    proxyEnabledCheck.Checked,
				ProxyServer:     proxyServerEntry.Text,
			}
			loadedProfiles = append(loadedProfiles, newProfile)
		}
		profiles.SaveProfiles(loadedProfiles)
		profileList.Refresh()
	})

	toolbar := container.NewHBox(newBtn, saveBtn, deleteBtn)
	split := container.NewHSplit(profileList, container.NewVBox(form, toolbar))
	split.Offset = 0.3

	w.SetContent(split)
	w.SetOnClosed(onClosed) // Refresh the main UI when this window closes
	return w
}

// --- Backend Logic ---

// ApplyProfile applies all settings from a given profile.
func ApplyProfile(p profiles.Profile) error {
	if p.NetworkPriority == "Ethernet" {
		if err := SetInterfaceMetric("Ethernet", 10); err != nil { return err }
		if err := SetInterfaceMetric("Wi-Fi", 20); err != nil { return err }
	} else if p.NetworkPriority == "Wi-Fi" {
		if err := SetInterfaceMetric("Wi-Fi", 10); err != nil { return err }
		if err := SetInterfaceMetric("Ethernet", 20); err != nil { return err }
	}

	return SetProxyState(p.ProxyEnabled, p.ProxyServer)
}

// SetInterfaceMetric sets the metric for a network interface.
func SetInterfaceMetric(name string, metric int) error {
	cmd := exec.Command("netsh", "interface", "ipv4", "set", "interface", fmt.Sprintf("interface=%s", name), fmt.Sprintf("metric=%d", metric))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}

// SetProxyState enables or disables the system proxy.
func SetProxyState(enable bool, server string) error {
	regPath := "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings"
	if enable {
		cmdEnable := exec.Command("reg", "add", regPath, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f")
		if _, err := cmdEnable.CombinedOutput(); err != nil { return err }

		cmdServer := exec.Command("reg", "add", regPath, "/v", "ProxyServer", "/t", "REG_SZ", "/d", server, "/f")
		if _, err := cmdServer.CombinedOutput(); err != nil { return err }
	} else {
		cmdDisable := exec.Command("reg", "add", regPath, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f")
		if _, err := cmdDisable.CombinedOutput(); err != nil { return err }
	}
	return nil
}
