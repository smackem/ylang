package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
	sourcePath := flag.String("source", "", "the source image path")
	targetPath := flag.String("target", "", "the target image path")
	flag.Parse()

	if *sourcePath == "" {
		flag.Usage()
		return
	}

	img, err := loadImage(*sourcePath)
	if err != nil {
		log.Fatalf("Error loading image %s: %s", *sourcePath, err.Error())
	}

	if err = saveImage(img, *targetPath); err != nil {
		log.Fatalf("Error saving image %s: %s", *targetPath, err.Error())
	}

	log.Printf("Saved image to %s as png", *targetPath)
}
