package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/Lec7ral/MultiTool/tools/network/networkswitcher"
	"github.com/Lec7ral/MultiTool/tools/profiles"
)

// InstallSystray configura e instala la bandeja del sistema y su menú.
func InstallSystray(app fyne.App, showWindow func()) {
	if desk, ok := app.(desktop.App); ok {
		// Función para construir/reconstruir el menú
		buildMenu := func() {
			iconResource, err := fyne.LoadResourceFromPath("icon.png")
			if err != nil {
				fyne.LogError("Failed to load systray icon", err)
				desk.SetSystemTrayIcon(theme.FyneLogo())
			} else {
				desk.SetSystemTrayIcon(iconResource)
			}

			menu := fyne.NewMenu("Toolbox",
				fyne.NewMenuItem("Open", showWindow),
			)

			// Cargar los perfiles de red desde el inicio.
			loadedProfiles, err := profiles.LoadProfiles()
			if err == nil && len(loadedProfiles) > 0 {
				profileSubMenu := fyne.NewMenu("")
				for _, p := range loadedProfiles {
					profile := p
					item := fyne.NewMenuItem(profile.Name, func() {
						go func() {
							if err := networkswitcher.ApplyProfile(profile); err != nil {
								app.SendNotification(&fyne.Notification{Title: "Toolbox", Content: "Failed to apply profile " + profile.Name})
							} else {
								app.SendNotification(&fyne.Notification{Title: "Toolbox", Content: "Profile '" + profile.Name + "' applied."})
							}
						}()
					})
					profileSubMenu.Items = append(profileSubMenu.Items, item)
				}

				modeMenuItem := fyne.NewMenuItem("Mode", nil)
				modeMenuItem.ChildMenu = profileSubMenu

				menu.Items = append(menu.Items, fyne.NewMenuItemSeparator())
				menu.Items = append(menu.Items, modeMenuItem)
			}

			menu.Items = append(menu.Items, fyne.NewMenuItemSeparator())
			menu.Items = append(menu.Items, fyne.NewMenuItem("Quit", func() {
				app.Quit()
			}))

			desk.SetSystemTrayMenu(menu)
		}

		buildMenu() // Construir el menú inicial

		// Pasar la función de reconstrucción a la herramienta de red
		networkswitcher.SetSystrayCallback(buildMenu)
	}
}
