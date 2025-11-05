package tools

import (
	"github.com/Lec7ral/MultiTool/tools/network/networkswitcher"
)

// NewNetworkSwitcherTool crea una instancia de la herramienta NetworkSwitcher.
func NewNetworkSwitcherTool() Tool {
	return networkswitcher.New()
}
