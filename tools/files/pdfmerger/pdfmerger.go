package pdfmerger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// --- Custom Sized Entry Widget ---
type sizedEntry struct {
	widget.Entry
	minWidth float32
}

func newSizedEntry(width float32) *sizedEntry {
	e := &sizedEntry{minWidth: width}
	e.ExtendBaseWidget(e)
	return e
}

func (e *sizedEntry) MinSize() fyne.Size {
	min := e.Entry.MinSize()
	if min.Width < e.minWidth {
		min.Width = e.minWidth
	}
	return min
}

// pdfFileItem holds data for a single file in the merge list.
type pdfFileItem struct {
	Path      string
	PageRange string
	PageCount int
}

// --- Tool Definition ---
type PDFMergerTool struct {
	pdfFiles []pdfFileItem
	fileList *widget.List
	icon     fyne.Resource // Cache del icono
}

func New() *PDFMergerTool {
	t := &PDFMergerTool{
		pdfFiles: make([]pdfFileItem, 0),
	}
	return t
}

func (t *PDFMergerTool) GetName() string {
	return "PDF Merger"
}

func (t *PDFMergerTool) GetDescription() string {
	return "Combine and reorder PDFs with page selection"
}

func (t *PDFMergerTool) GetCategory() string {
	return "Files"
}

func (t *PDFMergerTool) GetIcon() fyne.Resource {
	// Cargar el icono solo una vez y cachearlo.
	if t.icon == nil {
		resource, err := fyne.LoadResourceFromPath("assets/pdf.svg")
		if err != nil {
			fyne.LogError("Failed to load pdf icon", err)
			return nil
		}
		t.icon = resource
	}
	return t.icon
}

// OnFilesDropped is called by the app layout when files are dropped.
func (t *PDFMergerTool) OnFilesDropped(files []string) {
	for _, p := range files {
		path := p
		// On Windows, file URIs from Fyne can have a leading slash.
		// We remove it to ensure compatibility with file system operations.
		if len(path) > 2 && path[0] == '/' && path[2] == ':' {
			path = path[1:]
		}

		if filepath.Ext(path) == ".pdf" {
			count, err := api.PageCountFile(filepath.FromSlash(path))
			if err != nil {
				fyne.LogError("Failed to count pages for "+path, err)
			}
			t.pdfFiles = append(t.pdfFiles, pdfFileItem{Path: path, PageCount: count})
		}
	}
	if t.fileList != nil {
		t.fileList.Refresh()
	}
}

// --- Main UI ---
func (t *PDFMergerTool) GetUI(window fyne.Window) fyne.CanvasObject {
	var selectedIndex int = -1

	statusLabel := widget.NewLabel("Arrastra y suelta archivos o usa 'Añadir PDFs'. Para seleccionar páginas, usa rangos (ej: 2-5), números sueltos (ej: 8), rangos abiertos (ej: 12-) o exclusiones (ej: !10).")

	// --- File List with Page Range ---
	t.fileList = widget.NewList(
		func() int { return len(t.pdfFiles) },
		func() fyne.CanvasObject {
			pageEntry := newSizedEntry(150)
			pageEntry.SetPlaceHolder("e.g., 1-5, !3")
			return container.NewBorder(nil, nil, nil, pageEntry, widget.NewLabel("template"))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			label := c.Objects[0].(*widget.Label)
			labelText := filepath.Base(t.pdfFiles[i].Path)
			if t.pdfFiles[i].PageCount > 0 {
				labelText = fmt.Sprintf("%s (%d pages)", labelText, t.pdfFiles[i].PageCount)
			}
			label.SetText(labelText)

			entry := c.Objects[1].(*sizedEntry)
			entry.SetText(t.pdfFiles[i].PageRange)
			entry.OnChanged = func(s string) {
				t.pdfFiles[i].PageRange = s
			}
		},
	)
	t.fileList.OnSelected = func(id widget.ListItemID) { selectedIndex = id }

	// --- Action Buttons (Right Panel) ---
	addBtn := widget.NewButton("Add PDFs...", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			path := reader.URI().Path()
			// On Windows, file URIs from Fyne can have a leading slash.
			// We remove it to ensure compatibility with file system operations.
			if len(path) > 2 && path[0] == '/' && path[2] == ':' {
				path = path[1:]
			}
			count, err := api.PageCountFile(filepath.FromSlash(path))
			if err != nil {
				fyne.LogError("Failed to count pages for "+path, err)
			}
			t.pdfFiles = append(t.pdfFiles, pdfFileItem{Path: path, PageCount: count})
			t.fileList.Refresh()
		}, window)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
		fileDialog.Show()
	})

	removeBtn := widget.NewButton("Remove", func() {
		if selectedIndex < 0 || selectedIndex >= len(t.pdfFiles) {
			return
		}
		t.pdfFiles = append(t.pdfFiles[:selectedIndex], t.pdfFiles[selectedIndex+1:]...)
		selectedIndex = -1
		t.fileList.UnselectAll()
		t.fileList.Refresh()
	})

	moveUpBtn := widget.NewButton("Move Up", func() {
		if selectedIndex <= 0 {
			return
		}
		t.pdfFiles[selectedIndex], t.pdfFiles[selectedIndex-1] = t.pdfFiles[selectedIndex-1], t.pdfFiles[selectedIndex]
		t.fileList.Select(selectedIndex - 1)
	})

	moveDownBtn := widget.NewButton("Move Down", func() {
		if selectedIndex < 0 || selectedIndex >= len(t.pdfFiles)-1 {
			return
		}
		t.pdfFiles[selectedIndex], t.pdfFiles[selectedIndex+1] = t.pdfFiles[selectedIndex+1], t.pdfFiles[selectedIndex]
		t.fileList.Select(selectedIndex + 1)
	})

	actionButtons := container.NewVBox(addBtn, removeBtn, moveUpBtn, moveDownBtn)

	// --- Output & Merge (Bottom Panel) ---
	outputEntry := widget.NewEntry()
	outputEntry.Disable()

	saveAsBtn := widget.NewButton("Save As...", func() {
		fileDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			path := writer.URI().Path()
			if len(path) > 2 && path[0] == '/' && path[2] == ':' {
				path = path[1:]
			}
			outputEntry.SetText(path)
		}, window)
		fileDialog.SetFileName("merged.pdf")
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
		fileDialog.Show()
	})

	mergeBtn := widget.NewButton("Merge PDFs", func() {
		if len(t.pdfFiles) < 1 {
			statusLabel.SetText("Error: Please add at least one PDF file.")
			return
		}
		if outputEntry.Text == "" {
			statusLabel.SetText("Error: Please select an output file location.")
			return
		}
		statusLabel.SetText("Merging...")
		if err := mergePDFs(t.pdfFiles, outputEntry.Text); err != nil {
			statusLabel.SetText("Error: " + err.Error())
		} else {
			statusLabel.SetText("Success! PDFs merged into " + filepath.Base(outputEntry.Text))
		}
	})

	outputArea := container.NewBorder(nil, nil, nil, saveAsBtn, outputEntry)
	bottomPanel := container.NewVBox(outputArea, mergeBtn, statusLabel)

	// --- Final Layout ---
	listContainer := container.NewBorder(nil, nil, nil, actionButtons, t.fileList)
	return container.NewBorder(nil, bottomPanel, nil, nil, listContainer)
}

// --- Backend Logic ---
func mergePDFs(files []pdfFileItem, outFile string) error {
	if len(files) == 0 {
		return errors.New("no files to merge")
	}

	filePaths := make([]string, 0, len(files))
	tempDirs := make([]string, 0)

	defer func() {
		for _, d := range tempDirs {
			os.RemoveAll(d)
		}
	}()

	for _, f := range files {
		pageRange := strings.TrimSpace(f.PageRange)

		if pageRange != "" {
			pageSelectionRaw := strings.Split(pageRange, ",")
			pageSelection := make([]string, 0, len(pageSelectionRaw))
			for _, s := range pageSelectionRaw {
				trimmed := strings.TrimSpace(s)
				if trimmed != "" {
					pageSelection = append(pageSelection, trimmed)
				}
			}

			if len(pageSelection) == 0 {
				filePaths = append(filePaths, f.Path)
				continue
			}

			tempDir, err := os.MkdirTemp("", "pdfmerger-")
			if err != nil {
				return fmt.Errorf("failed to create temp dir: %w", err)
			}
			tempDirs = append(tempDirs, tempDir)

			sourcePath := filepath.FromSlash(f.Path)

			if err := api.ExtractPagesFile(sourcePath, tempDir, pageSelection, nil); err != nil {
				return fmt.Errorf("failed to extract pages from '%s' (pages: %s): %w", filepath.Base(sourcePath), pageRange, err)
			}

			dirEntries, err := os.ReadDir(tempDir)
			if err != nil {
				return fmt.Errorf("failed to read temp dir %s: %w", tempDir, err)
			}
			if len(dirEntries) == 0 {
				return fmt.Errorf("page extraction produced no files for '%s' (pages: %s)", filepath.Base(sourcePath), pageRange)
			}

			for _, entry := range dirEntries {
				extractedFilePath := filepath.Join(tempDir, entry.Name())
				filePaths = append(filePaths, extractedFilePath)
			}

		} else {
			filePaths = append(filePaths, f.Path)
		}
	}

	if err := api.MergeCreateFile(filePaths, outFile, false, nil); err != nil {
		return fmt.Errorf("failed to merge pdfs: %w", err)
	}

	return nil
}
