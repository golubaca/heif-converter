package main

import (
	"log"

	"github.com/MaestroError/go-libheif"
)

func main() {
	err := libheif.HeifToJpeg("C:/Users/Golub/Downloads/IMG_3161.HEIC", "C:/Users/Golub/Downloads/IMG_3161.HEIC.jpg", 80)
	if err != nil {
		log.Fatal(err)
	}

	err = libheif.HeifToPng("C:/Users/Golub/Downloads/IMG_3161.HEIC", "C:/Users/Golub/Downloads/IMG_3161.HEIC.png")
	if err != nil {
		log.Fatal(err)
	}
}
