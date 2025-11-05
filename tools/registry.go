package tools

import "fyne.io/fyne/v2"

type Tool interface {
	GetName() string
	GetDescription() string
	GetCategory() string
	GetIcon() fyne.Resource
	GetUI() fyne.CanvasObject
}

type ToolRegistry struct {
	tools map[string]Tool
	order []string
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
		order: make([]string, 0),
	}
}

func (tr *ToolRegistry) Register(name string, tool Tool) {
	if _, exists := tr.tools[name]; !exists {
		tr.order = append(tr.order, name)
	}
	tr.tools[name] = tool
}

func (tr *ToolRegistry) Get(name string) Tool {
	return tr.tools[name]
}

func (tr *ToolRegistry) GetAll() []Tool {
	result := make([]Tool, 0, len(tr.order))
	for _, name := range tr.order {
		result = append(result, tr.tools[name])
	}
	return result
}

func RegisterDefaultTools(registry *ToolRegistry) {
	// Solo registrar la herramienta principal
	registry.Register("Network Switcher", NewNetworkSwitcherTool())
}
