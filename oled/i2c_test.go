package oled

import (
	"fmt"
	"testing"
)

func TestI2C(t *testing.T) {
	dev, err := Open(&I2cOpener{})
	if err != nil {
		fmt.Printf("Failed to open: %v", err)
		t.FailNow()
	}
	if dev == nil {
		fmt.Print("Open returned nil")
		t.FailNow()
	}
	defer dev.Close()
	dev.Print(0, 0, "QUICK BROWN FOX JUMPS")
	dev.Print(1, 0, "OVER THE LAZY DOG")
	dev.Print(2, 0, "!\"#$%&'()*+,-./\\[]^")
	dev.Print(3, 0, "0123456789:;<=>?_`@")
}
