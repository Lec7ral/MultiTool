package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Lec7ral/MultiTool/tools"
)

// CreateAppLayout construye y devuelve el layout principal de la aplicación.
func CreateAppLayout(myWindow fyne.Window) fyne.CanvasObject {
	toolRegistry := tools.NewToolRegistry()
	tools.RegisterDefaultTools(toolRegistry)

	// Agrupar herramientas por categoría
	categories := make(map[string][]tools.Tool)
	categoryOrder := []string{"System", "Files", "Text", "Network"} // Orden deseado
	for _, tool := range toolRegistry.GetAll() {
		categories[tool.GetCategory()] = append(categories[tool.GetCategory()], tool)
	}

	// Mapa de iconos para las categorías
	categoryIcons := map[string]fyne.Resource{
		"System":  theme.SettingsIcon(),
		"Files":   theme.FolderIcon(),
		"Text":    theme.DocumentIcon(),
		"Network": theme.ComputerIcon(),
	}

	// --- Pestañas de Categorías (Nivel Superior) ---
	categoryTabs := container.NewAppTabs()

	for _, categoryName := range categoryOrder {
		if toolsInCat, ok := categories[categoryName]; ok {

			// --- Contenido de la Herramienta (Panel Derecho) ---
			toolContent := container.NewMax()

			// --- Pestañas de Herramientas (Panel Izquierdo) ---
			toolTabs := container.NewAppTabs()
			toolTabs.SetTabLocation(container.TabLocationLeading)

			// Mapa para asociar cada TabItem con su herramienta correspondiente
			tabToToolMap := make(map[*container.TabItem]tools.Tool)

			for _, tool := range toolsInCat {
				tabItem := container.NewTabItemWithIcon(tool.GetName(), tool.GetIcon(), container.NewWithoutLayout())
				toolTabs.Append(tabItem)
				tabToToolMap[tabItem] = tool
			}

			toolTabs.OnSelected = func(selectedTab *container.TabItem) {
				if tool, ok := tabToToolMap[selectedTab]; ok {
					toolContent.Objects = []fyne.CanvasObject{tool.GetUI()}
					toolContent.Refresh()
				}
			}

			// Cargar la primera herramienta de la categoría por defecto
			if len(toolTabs.Items) > 0 {
				toolTabs.SelectIndex(0)
				if firstTool, ok := tabToToolMap[toolTabs.Items[0]]; ok {
					toolContent.Objects = []fyne.CanvasObject{firstTool.GetUI()}
					toolContent.Refresh()
				}
			}

			layout := container.NewBorder(nil, nil, toolTabs, nil, toolContent)
			categoryTabs.Append(container.NewTabItemWithIcon(categoryName, categoryIcons[categoryName], layout))
		}
	}

	// --- Barra de Estado Inferior ---
	projectURL, _ := url.Parse("https://github.com/Lec7ral/MultiTool")
	aboutButton := widget.NewButton("About", func() {
		aboutContent := container.NewVBox(
			widget.NewLabelWithStyle("MultiTool v1.0.0", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Developed by Lec7ral", fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewHyperlinkWithStyle("Project on GitHub", projectURL, fyne.TextAlignCenter, fyne.TextStyle{}),
		)
		dialog.ShowCustom("About", "Close", aboutContent, myWindow)
	})

	// Usar un layout vacío en lugar de nil para el Border
	statusBar := container.NewBorder(container.NewWithoutLayout(), container.NewWithoutLayout(), container.NewWithoutLayout(), aboutButton, container.NewWithoutLayout())

	// --- Layout Principal Final ---
	return container.NewBorder(nil, statusBar, nil, nil, categoryTabs)
}
