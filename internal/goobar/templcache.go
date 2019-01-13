package goobar

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type parsedTemplate struct {
	templ       *template.Template
	lastChecked time.Time
	lastModTime time.Time
}

var templateCache = struct {
	sync.Locker
	items map[string]parsedTemplate
}{
	new(sync.Mutex),
	make(map[string]parsedTemplate),
}

// getTemplate returns a pointer to the specified rendered template or an error.
// relativePath is the slash-separated file path of the template, relative to the path
// set with SetViewFolder.
func getTemplate(relativePath string) (templates *template.Template, templName string, err error) {
	now := time.Now()
	relativeDir, templName := path.Split(relativePath)
	absoluteDir := path.Join(viewFolder, relativeDir)
	log.Printf("getTemplate: relativePath=%s, relativeDir=%s, templName=%s, absoluteDir=%s", relativePath, relativeDir, templName, absoluteDir)

	templateCache.Lock()
	defer templateCache.Unlock()
	pt, ok := templateCache.items[relativePath]

	if ok && now.Second() != pt.lastChecked.Second() { // check at most once per second
		outdated, err := hasNewerFile(absoluteDir, pt.lastModTime)
		if err != nil {
			return nil, "", err
		}
		log.Printf("getTemplate: relativePath=%s, found in cache, outdated=%v", relativePath, outdated)
		ok = !outdated
	}

	if !ok {
		templ, lastModTime, err := parseTemplateFiles(absoluteDir)
		if err != nil {
			return nil, "", err
		}
		log.Printf("getTemplate: relativePath=%s, parsed from disk", relativePath)
		pt.templ = templ
		pt.lastModTime = lastModTime
	}

	pt.lastChecked = now
	templateCache.items[relativePath] = pt
	return pt.templ, templName, nil
}

func parseTemplateFiles(absoluteDir string) (templ *template.Template, lastModTime time.Time, err error) {
	var newestModTime time.Time
	files, err := ioutil.ReadDir(filepath.FromSlash(absoluteDir))
	if err != nil {
		return nil, newestModTime, err
	}

	fileNames := []string{}
	for _, file := range files {
		if file.Mode()&os.ModeType == 0 { // regular files only
			modTime := file.ModTime()
			if modTime.After(newestModTime) {
				newestModTime = modTime
			}
			fileNames = append(fileNames, filepath.Join(absoluteDir, file.Name()))
		}
	}

	templ, err = template.ParseFiles(fileNames...)
	return templ, newestModTime, err
}

func hasNewerFile(absoluteDir string, t time.Time) (bool, error) {
	files, err := ioutil.ReadDir(filepath.FromSlash(absoluteDir))
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if file.Mode()&os.ModeType == 0 { // regular files only
			modTime := file.ModTime()
			if modTime.After(t) {
				return true, nil
			}
		}
	}

	return false, nil
}
