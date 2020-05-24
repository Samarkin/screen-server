package main

import (
	"log"
	"os"

	oled ".."
)

func main() {
	var filename string
	if len(os.Args) < 2 {
		filename = "../../img/logo.png"
	} else {
		filename = os.Args[1]
	}
	dev, err := oled.OpenI2C()
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	err = dev.DisplayImage(filename)
	if err != nil {
		log.Fatal(err)
	}
}
