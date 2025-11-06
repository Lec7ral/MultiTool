package tools

import "fyne.io/fyne/v2"

// Tool defines the interface for all tools in the application.
type Tool interface {
	GetName() string
	GetDescription() string
	GetCategory() string
	GetIcon() fyne.Resource
	GetUI(fyne.Window) fyne.CanvasObject
}

// FileDropper is an optional interface for tools that can handle dropped files.
type FileDropper interface {
	OnFilesDropped(files []string)
}

// ToolDescriptor contiene la información estática de una herramienta y cómo crearla.
type ToolDescriptor struct {
	Name        string
	Category    string
	Icon        fyne.Resource
	Constructor func() Tool // Función para crear la instancia completa de la herramienta
}

// ToolRegistry gestiona los descriptores de herramientas y un caché de instancias.
type ToolRegistry struct {
	toolDescriptors map[string]ToolDescriptor
	toolInstances   map[string]Tool
	order           []string
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		toolDescriptors: make(map[string]ToolDescriptor),
		toolInstances:   make(map[string]Tool),
		order:           make([]string, 0),
	}
}

// Register registra un descriptor de herramienta.
func (tr *ToolRegistry) Register(descriptor ToolDescriptor) {
	if _, exists := tr.toolDescriptors[descriptor.Name]; !exists {
		tr.order = append(tr.order, descriptor.Name)
	}
	tr.toolDescriptors[descriptor.Name] = descriptor
}

// Get obtiene una instancia de la herramienta, creándola si es necesario (carga perezosa).
func (tr *ToolRegistry) Get(name string) Tool {
	if instance, ok := tr.toolInstances[name]; ok {
		return instance
	}

	if descriptor, ok := tr.toolDescriptors[name]; ok {
		instance := descriptor.Constructor() // Llama a la función constructora.
		tr.toolInstances[name] = instance
		return instance
	}

	return nil
}

// GetAllDescriptors devuelve todos los descriptores de herramientas registrados.
func (tr *ToolRegistry) GetAllDescriptors() []ToolDescriptor {
	result := make([]ToolDescriptor, 0, len(tr.order))
	for _, name := range tr.order {
		result = append(result, tr.toolDescriptors[name])
	}
	return result
}

// RegisterDefaultTools registra los descriptores de las herramientas predeterminadas.
func RegisterDefaultTools(registry *ToolRegistry) {
	// Para obtener la información estática (icono, categoría) sin crear la herramienta,
	// necesitamos una forma de acceder a ella. La solución más limpia es tener
	// una instancia "ligera" o prototipo, o simplemente definirla aquí.

	// Prototipo de PDFMerger para obtener sus metadatos.
	pdfMergerProto := NewPDFMergerTool()
	registry.Register(ToolDescriptor{
		Name:        pdfMergerProto.GetName(),
		Category:    pdfMergerProto.GetCategory(),
		Icon:        pdfMergerProto.GetIcon(),
		Constructor: NewPDFMergerTool,
	})

	// Prototipo de NetworkSwitcher para obtener sus metadatos.
	networkSwitcherProto := NewNetworkSwitcherTool()
	registry.Register(ToolDescriptor{
		Name:        networkSwitcherProto.GetName(),
		Category:    networkSwitcherProto.GetCategory(),
		Icon:        networkSwitcherProto.GetIcon(),
		Constructor: NewNetworkSwitcherTool,
	})
}
