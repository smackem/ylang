package main

import (
	"fmt"
	"github.com/smackem/ylang/internal/emitter"
	"github.com/smackem/ylang/internal/program"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/smackem/ylang/internal/goobar"
)

func initHTTP() error {
	err := os.Mkdir(targetImageDir, os.ModeDir)
	if err != nil && os.IsExist(err) == false {
		return err // ignore "already exists" error
	}
	goobar.SetViewFolder("res/view")
	http.Handle("/", goobar.Get(getIndex))
	http.Handle("/render", goobar.Post(postRender))
	http.Handle("/pub/", http.FileServer(http.Dir("res")))
	return nil
}

const targetImageDir string = "res/pub/temp"

func getIndex(x *goobar.Exchange) goobar.Responder {
	return goobar.View("index.html", nil)
}

func postRender(x *goobar.Exchange) goobar.Responder {
	uri := x.MustGetString("imageUri")
	source := x.MustGetString("sourceCode")

	resp, err := http.Get(uri)
	if err != nil {
		return goobar.Error(500, fmt.Sprintf("error loading '%s'", uri))
	}

	surf, err := loadSurface(resp.Body)
	if err != nil {
		return goobar.Error(500, fmt.Sprintf("error decoding '%s'", uri))
	}

	prog, err := program.Compile(string(source))
	if err != nil {
		return goobar.Error(500, fmt.Sprintf("compilation error: %s", err.Error()))
	}

	err = program.Execute(prog, surf)
	if err != nil {
		return goobar.Error(500, fmt.Sprintf("execution error: %s", err.Error()))
	}

	guid := uuid.New()
	targetFileName := fmt.Sprintf("%s.png", guid.String())
	//targetFilePath := path.Join(targetImageDir, targetFileName)

	//if err = saveImage(surf.target, targetFilePath); err != nil {
	//	return goobar.Error(500, fmt.Sprintf("error saving image %s: %s", targetFilePath, err.Error()))
	//}

	return goobar.JSON(struct {
		ImagePath string `json:"imagepath"`
		JsCode    string `json:"jscode"`
	}{
		ImagePath: path.Join("/pub/temp", targetFileName),
		JsCode:    emitter.EmitJS(prog),
	})
}
