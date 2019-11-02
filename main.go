package main

import (
	"flag"
	"fmt"
	"github.com/smackem/ylang/internal/program"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	port = 9090
)

func webMain() {
	err := initHTTP()
	if err != nil {
		log.Fatalf("error initializing http server: %s", err.Error())
	}

	srv := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	fmt.Printf("Running on port %d. Press Ctrl+C to quit...", port)
	err = srv.ListenAndServe()
	log.Print(err)
}

func main() {
	sourceImgPath := flag.String("image", "", "the source image path")
	sourceCodePath := flag.String("code", "", "the path of the source code file")
	targetImgPath := flag.String("out", "", "the target image path")
	flag.Parse()

	if *sourceImgPath == "" {
		webMain()
		return
	}

	sourceFile, err := os.Open(*sourceImgPath)
	if err != nil {
		log.Fatalf("Could not load %s: %s", *sourceImgPath, err.Error())
	}
	defer sourceFile.Close()

	surf, err := loadSurface(sourceFile)
	if err != nil {
		log.Fatalf("error loading image from '%s': %s", *sourceImgPath, err.Error())
	}
	src, err := ioutil.ReadFile(*sourceCodePath)
	if err != nil {
		log.Fatalf("error loading source code from '%s': %s", *sourceCodePath, err.Error())
	}

	prog, err := program.Compile(string(src))
	if err != nil {
		log.Fatalf("compilation error: %s", err.Error())
	}

	start := time.Now()
	err = program.Execute(prog, surf)
	if err != nil {
		log.Fatalf("execution error: %s", err.Error())
	}
	log.Printf("execution took %s", time.Since(start))

	if err = saveImage(surf.target, *targetImgPath); err != nil {
		log.Fatalf("error saving image %s: %s", *targetImgPath, err.Error())
	}

	log.Printf("Saved image to '%s' as png", *targetImgPath)
}
