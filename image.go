package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

func loadImage(sourcePath string) (*image.NRGBA, error) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("Could not load %s: %s", sourcePath, err.Error())
	}
	defer sourceFile.Close()

	source, encoding, err := image.Decode(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("Error decoding %s: %s", sourcePath, err.Error())
	}

	log.Printf("Image %s decoded as %s", sourcePath, encoding)

	target := image.NewNRGBA(source.Bounds())
	draw.Draw(target, target.Bounds(), source, image.Point{0, 0}, draw.Src)

	return target, nil
}

func saveImage(img image.Image, targetPath string) error {
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("Could not create file %s: %s", targetPath, err)
	}
	defer targetFile.Close()

	if err = png.Encode(targetFile, img); err != nil {
		return fmt.Errorf("Could not encode target image: %s", err)
	}

	return nil
}
