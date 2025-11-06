# MultiTool

MultiTool es una aplicación de escritorio versátil escrita en Go que reúne una colección de herramientas prácticas para simplificar tus tareas diarias. La aplicación está diseñada con una interfaz de usuario limpia e intuitiva gracias al framework Fyne.

## Características

Actualmente, MultiTool incluye las siguientes herramientas:

*   **Fusión de PDFs:** Combina múltiples archivos PDF en uno solo. Permite reordenar los archivos, y seleccionar páginas específicas o rangos de páginas de cada PDF antes de unirlos.
*   **Cambiador de Red:** (Descripción de la herramienta de cambio de red)
*   **Gestor de Perfiles:** (Descripción del gestor de perfiles)

## Cómo Empezar

### Prerrequisitos

*   Tener instalado Go (versión 1.16 o superior).
*   Tener un compilador de C/C++ (como GCC) instalado para las dependencias de Fyne.

### Instalación y Ejecución

1.  Clona el repositorio:
    ```sh
    git clone https://github.com/Lec7ral/MultiTool.git
    ```
2.  Navega al directorio del proyecto:
    ```sh
    cd MultiTool
    ```
3.  Ejecuta la aplicación:
    ```sh
    go run .
    ```
4.  Compila la aplicación:
    ```sh
    fyne package -os windows -icon assets/icon.ico -release --app-id com.Lec7ral.multitool
    ```
## Guía de Uso

### Fusión de PDFs

1.  **Añadir Archivos:** Puedes añadir archivos PDF a la lista de dos maneras:
    *   **Arrastrar y Soltar:** Simplemente arrastra los archivos PDF desde tu explorador de archivos y suéltalos en cualquier parte de la ventana de la aplicación.
    *   **Botón 'Añadir PDFs...':** Haz clic en este botón para abrir un diálogo de selección de archivos.

2.  **Selección de Páginas:** Al lado de cada archivo en la lista, encontrarás un campo de texto para especificar qué páginas quieres incluir. Si lo dejas en blanco, se incluirá el PDF completo. La sintaxis es muy flexible:
    *   **Rangos:** `2-5` (incluye las páginas de la 2 a la 5).
    *   **Números Sueltos:** `8` (incluye solo la página 8). Puedes combinarlo con rangos: `2-5, 8`.
    *   **Rangos Abiertos:** `12-` (incluye desde la página 12 hasta el final).
    *   **Exclusiones:** `!10` (incluye todas las páginas excepto la 10). Puedes combinarlo: `1-15, !10, !12`.

3.  **Ordenar Archivos:** Usa los botones `Mover Arriba` y `Mover Abajo` para cambiar el orden en que los archivos serán fusionados.

4.  **Fusionar:**
    *   Haz clic en `Guardar Como...` para elegir la ubicación y el nombre del archivo PDF resultante.
    *   Haz clic en `Fusionar PDFs` para iniciar el proceso. Un mensaje en la barra de estado te informará del resultado.
