package tools

import (
	"github.com/Lec7ral/MultiTool/tools/files/pdfmerger"
	"github.com/Lec7ral/MultiTool/tools/network/networkswitcher"
)

// NewNetworkSwitcherTool crea una instancia de la herramienta NetworkSwitcher.
func NewNetworkSwitcherTool() Tool {
	return networkswitcher.New()
}

// NewPDFMergerTool crea una instancia de la herramienta PDFMerger.
func NewPDFMergerTool() Tool {
	return pdfmerger.New()
}
