package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Lec7ral/MultiTool/ui"
)

func main() {
	// 1. Inicializar la aplicación con un ID único para la persistencia
	myApp := app.NewWithID("com.lec7ral.multitool")
	myWindow := myApp.NewWindow("Toolbox - Windows Utilities")
	myWindow.Resize(fyne.NewSize(1050, 600))
	myWindow.SetMaster()

	// 2. Aplicar el tema personalizado
	ui.ApplyTheme(myApp)

	// 3. Configurar la lógica de cierre y la creación retardada de la bandeja del sistema
	systrayReady := false
	myWindow.SetCloseIntercept(func() {
		myWindow.Hide()
		if !systrayReady {
			ui.InstallSystray(myApp, myWindow)
			systrayReady = true
		}
	})

	// 4. Construir el layout principal de la UI
	mainLayout := ui.CreateAppLayout(myWindow)

	// 5. Establecer el contenido y ejecutar
	myWindow.SetContent(mainLayout)
	myWindow.ShowAndRun()
}
