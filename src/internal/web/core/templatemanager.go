package core

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type TemplateManager struct {
	templates map[string]*template.Template
	root      string
	layouts   string
	partials  string
	pages     string
}

func NewTemplateManager(root string) *TemplateManager {
	mgr := &TemplateManager{
		templates: make(map[string]*template.Template),
		root:      root,
		layouts:   filepath.Join(root, "layouts"),
		partials:  filepath.Join(root, "partials"),
		pages:     filepath.Join(root, "pages"),
	}
	mgr.loadTemplates()
	return mgr
}

func (tm *TemplateManager) loadTemplates() {
	layouts, err := filepath.Glob(filepath.Join(tm.layouts, "*.html"))
	if err != nil {
		log.Fatalf("Failed to read layouts: %v", err)
	}
	partials, err := filepath.Glob(filepath.Join(tm.partials, "*.html"))
	if err != nil {
		log.Fatalf("Failed to read partials: %v", err)
	}

	// Walk pages directory recursively
	err = filepath.Walk(tm.pages, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		relPath, err := filepath.Rel(tm.pages, path)
		if err != nil {
			return err
		}

		files := append([]string{}, layouts...)
		files = append(files, partials...)
		files = append(files, path)

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Printf("Failed to parse template %s: %v", relPath, err)
			return nil
		}
		relPath = filepath.ToSlash(relPath) // for Windows compatibility
		tm.templates[relPath] = tmpl
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk pages: %v", err)
	}
}

func (tm *TemplateManager) Render(w http.ResponseWriter, name string, data any, layout string) {
	tmpl, ok := tm.templates[name]
	if !ok {
		http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
		return
	}
	err := tmpl.ExecuteTemplate(w, layout, data)
	if err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
	}
}
