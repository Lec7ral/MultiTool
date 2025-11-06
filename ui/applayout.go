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
func CreateAppLayout() (fyne.CanvasObject, func(fyne.Window)) {
	toolRegistry := tools.NewToolRegistry()
	tools.RegisterDefaultTools(toolRegistry)

	// Agrupar descriptores de herramientas por categoría
	categories := make(map[string][]tools.ToolDescriptor)
	categoryOrder := []string{"System", "Files", "Text", "Network"} // Orden deseado
	for _, descriptor := range toolRegistry.GetAllDescriptors() {
		categories[descriptor.Category] = append(categories[descriptor.Category], descriptor)
	}

	categoryIcons := map[string]fyne.Resource{
		"System":  theme.SettingsIcon(),
		"Files":   theme.FolderIcon(),
		"Text":    theme.DocumentIcon(),
		"Network": theme.ComputerIcon(),
	}

	// --- Pestañas de Categorías (Nivel Superior) ---
	categoryTabs := container.NewAppTabs()

	for _, categoryName := range categoryOrder {
		if descriptorsInCat, ok := categories[categoryName]; ok {

			// --- Contenido de la Herramienta (Panel Derecho) ---
			toolContent := container.NewMax()

			// --- Pestañas de Herramientas (Panel Izquierdo) ---
			toolTabs := container.NewAppTabs()
			toolTabs.SetTabLocation(container.TabLocationLeading)

			// Mapa para asociar cada TabItem con su descriptor de herramienta
			tabToDescriptorMap := make(map[*container.TabItem]tools.ToolDescriptor)

			for _, descriptor := range descriptorsInCat {
				// El contenido inicial de la pestaña está vacío. La herramienta no se crea aquí.
				tabItem := container.NewTabItemWithIcon(descriptor.Name, descriptor.Icon, container.NewWithoutLayout())
				toolTabs.Append(tabItem)
				tabToDescriptorMap[tabItem] = descriptor
			}

			toolTabs.OnSelected = func(selectedTab *container.TabItem) {
				if selectedTab == nil {
					return
				}
				if descriptor, ok := tabToDescriptorMap[selectedTab]; ok {
					// Obtenemos la herramienta (se crea aquí si es la primera vez).
					tool := toolRegistry.Get(descriptor.Name)
					if tool != nil {
						toolContent.Objects = []fyne.CanvasObject{tool.GetUI(nil)}
						toolContent.Refresh()
					}
				}
			}

			// Cargar la primera herramienta de la categoría por defecto.
			if len(toolTabs.Items) > 0 {
				// Seleccionamos la primera pestaña.
				toolTabs.SelectIndex(0)

				// Y cargamos su contenido manualmente para asegurar que la UI inicial aparezca,
				// ya que SelectIndex() no siempre dispara OnSelected() al inicio.
				firstTab := toolTabs.Items[0]
				if descriptor, ok := tabToDescriptorMap[firstTab]; ok {
					tool := toolRegistry.Get(descriptor.Name)
					if tool != nil {
						toolContent.Objects = []fyne.CanvasObject{tool.GetUI(nil)}
						toolContent.Refresh()
					}
				}
			} else {
				// Si no hay herramientas en la categoría, limpiamos el contenido.
				toolContent.Objects = nil
				toolContent.Refresh()
			}

			layout := container.NewBorder(nil, nil, toolTabs, nil, toolContent)
			categoryTabs.Append(container.NewTabItemWithIcon(categoryName, categoryIcons[categoryName], layout))
		}
	}

	// --- Lógica de Arrastrar y Soltar (Drag and Drop) ---
	setupWindowCallbacks := func(w fyne.Window) {
		w.SetOnDropped(func(p fyne.Position, uris []fyne.URI) {
			if categoryTabs.Selected().Text == "Files" {
				// Obtenemos la instancia de PDF Merger solo cuando se necesita.
				pdfMergerInstance := toolRegistry.Get("PDF Merger")
				if dropper, ok := pdfMergerInstance.(tools.FileDropper); ok {
					var filePaths []string
					for _, u := range uris {
						filePaths = append(filePaths, u.Path())
					}
					dropper.OnFilesDropped(filePaths)
				}
			}
		})
	}

	// --- Barra de Estado Inferior ---
	projectURL, _ := url.Parse("https://github.com/Lec7ral/MultiTool")
	aboutButton := widget.NewButton("About", func() {
		if fyne.CurrentApp().Driver().AllWindows() != nil && len(fyne.CurrentApp().Driver().AllWindows()) > 0 {
			activeWindow := fyne.CurrentApp().Driver().AllWindows()[0]
			aboutContent := container.NewVBox(
				widget.NewLabelWithStyle("MultiTool v1.0.0", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				widget.NewLabelWithStyle("Developed by Lec7ral", fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewHyperlinkWithStyle("Project on GitHub", projectURL, fyne.TextAlignCenter, fyne.TextStyle{}),
			)
			dialog.ShowCustom("About", "Close", aboutContent, activeWindow)
		}
	})

	statusBar := container.NewBorder(nil, nil, nil, aboutButton, container.NewWithoutLayout())

	// --- Layout Principal Final ---
	mainLayout := container.NewBorder(nil, statusBar, nil, nil, categoryTabs)

	return mainLayout, setupWindowCallbacks
}
