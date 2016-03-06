package main

import (
	"./guesser"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func main() {
	img, _, err := image.Decode(os.Stdin)
	if err != nil {
		panic(err.Error())
	}
	println(guesser.Guess(img).String())
}
