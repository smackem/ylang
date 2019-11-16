package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/smackem/ylang/internal/emitter"
	"github.com/smackem/ylang/internal/program"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	sourceImgPath := flag.String("image", "", "the source image path")
	sourceCodePath := flag.String("code", "", "the path of the source code file")
	targetImgPath := flag.String("out", "", "the target image path")
	jsOutputPath := flag.String("js", "", "the javascript output path")
	flag.Parse()

	if *sourceImgPath == "" {
		serverMain()
		return
	}

	sourceFile, err := os.Open(*sourceImgPath)
	if err != nil {
		log.Fatalf("Could not load %s: %s", *sourceImgPath, err.Error())
	}
	defer func() { _ = sourceFile.Close() }()

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

	if jsOutputPath != nil {
		js := emitter.EmitJS(prog)
		if err := ioutil.WriteFile(*jsOutputPath, []byte(js), 0644); err != nil {
			log.Fatalf("error writing javascript to '%s': %s", *jsOutputPath, err)
		}
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

func serverMain() {
	go httpMain()
	go listenerMain()
	fmt.Printf("Running on port %d (HTTP) and %s (gRPC). Press Enter to quit...", httpPort, listenerPort)
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
