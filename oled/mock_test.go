package oled

import (
	"fmt"
	"testing"
)

func TestMock(t *testing.T) {
	dev, err := Open(&MockOpener{})
	if err != nil {
		fmt.Printf("Failed to open: %v", err)
		t.FailNow()
	}
	err = dev.Close()
	if err != nil {
		fmt.Printf("Failed to close: %v", err)
		t.FailNow()
	}
}
