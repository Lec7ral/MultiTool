package main

import (
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Lec7ral/MultiTool/ui"
)

// myApp y myWindow son globales para ser accesibles desde múltiples funciones.
var (
	myApp    fyne.App
	myWindow fyne.Window
)

func main() {
	// 1. Inicializar la aplicación.
	myApp = app.NewWithID("com.lec7ral.multitool")

	// 2. Aplicar el tema.
	ui.ApplyTheme(myApp)

	// 3. Instalar la bandeja del sistema desde el principio. Esto es crucial para que
	//    la aplicación no se cierre cuando la última ventana se cierre.
	ui.InstallSystray(myApp, createAndShowMainWindow)

	// 4. Crear y mostrar la ventana principal por primera vez.
	createAndShowMainWindow()

	// 5. Ejecutar el bucle principal de la aplicación.
	myApp.Run()
}

// createAndShowMainWindow encapsula toda la lógica para crear y configurar la ventana.
func createAndShowMainWindow() {
	// Si la ventana ya existe y es visible, simplemente la enfocamos y salimos.
	if myWindow != nil {
		myWindow.RequestFocus()
		return
	}

	// Creamos una nueva ventana.
	w := myApp.NewWindow("Toolbox - Windows Utilities")
	myWindow = w // La asignamos a nuestra variable global.

	w.Resize(fyne.NewSize(1050, 600))

	// ¡¡NO LLAMAR A SetMaster()!!
	// Al no tener una ventana maestra, el ciclo de vida de la app no está atado a esta ventana.
	// Fyne no cerrará la app si la bandeja del sistema está activa.

	// Construimos el layout principal y obtenemos la función para configurar los callbacks.
	mainLayout, setupCallbacks := ui.CreateAppLayout()
	w.SetContent(mainLayout)

	// Configuramos los callbacks (como OnDropped) para esta ventana específica.
	setupCallbacks(w)

	// Interceptamos el cierre de la ventana.
	w.SetCloseIntercept(func() {
		w.Close() // Ahora esto es seguro. Cierra la ventana pero no la app.
	})

	// Cuando la ventana se cierre (después de w.Close()), limpiamos nuestra
	// referencia a ella y le pedimos al recolector de basura que se ejecute.
	w.SetOnClosed(func() {
		myWindow = nil // Eliminamos la referencia.
		runtime.GC()   // Sugerimos una recolección de basura.
	})

	w.Show()
}
