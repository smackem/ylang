package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/smackem/ylang/internal/lang"
)

const (
	port = 9090
)

func webMain() {
	registerHTTP()

	srv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	fmt.Printf("Running on port %d. Press Ctrl+C to quit...", port)
	err := srv.ListenAndServe()
	log.Print(err)
}

func main() {
	sourceImgPath := flag.String("image", "", "the source image path")
	sourceCodePath := flag.String("code", "", "the path of the source code file")
	targetImgPath := flag.String("out", "", "the target image path")
	flag.Parse()

	if *sourceImgPath == "" {
		flag.Usage()
		return
	}

	surf, err := loadSurface(*sourceImgPath)
	if err != nil {
		log.Fatalf("error loading image from '%s': %s", *sourceImgPath, err.Error())
	}
	src, err := ioutil.ReadFile(*sourceCodePath)
	if err != nil {
		log.Fatalf("error loading source code from '%s': %s", *sourceCodePath, err.Error())
	}

	prog, err := lang.Compile(string(src))
	if err != nil {
		log.Fatalf("compilation error: %s", err.Error())
	}

	err = prog.Execute(surf)
	if err != nil {
		log.Fatalf("execution error: %s", err.Error())
	}

	if err = saveImage(surf.target, *targetImgPath); err != nil {
		log.Fatalf("error saving image %s: %s", *targetImgPath, err.Error())
	}

	log.Printf("Saved image to '%s' as png", *targetImgPath)
}
