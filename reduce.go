package main

import (
	"./guesser"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
)

func main() {
	img, _, err := image.Decode(os.Stdin)
	if err != nil {
		panic(err.Error())
	}
	png.Encode(os.Stdout, guesser.Reduce(img))
}
