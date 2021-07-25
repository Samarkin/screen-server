package oled

import (
	"io"
	"log"
)

// MockOpener allows to open a screen object that does not perform any real connection
// Can be used for testing
type MockOpener struct {
}

type mockScreen struct {
	open bool
}

func (o *MockOpener) open() (Screen, error) {
	screen := &mockScreen{}
	screen.open = true
	log.Printf("Mock screen opened")
	return screen, nil
}

func (o *mockScreen) Clear() error {
	if !o.open {
		return ErrorScreenClosed
	}
	log.Printf("Mock screen cleared")
	return nil
}

func (o *mockScreen) Close() error {
	if !o.open {
		log.Printf("Attempt to close an already closed screen")
	}
	log.Printf("Mock screen closed")
	return nil
}

func (o *mockScreen) Print(line int, offset int, message string) error {
	if !o.open {
		return ErrorScreenClosed
	}
	log.Printf("Mock screen is now displaying message \"%s\" at line %d, offset %d", message, line, offset)
	return nil
}

func (o *mockScreen) DisplaySignalLevel(line int, offset int, level int) error {
	if !o.open {
		return ErrorScreenClosed
	}
	log.Printf("Mock screen is now displaying signal level %d at line %d, offset %d", level, line, offset)
	return nil
}

func (o *mockScreen) DisplayImageFile(filepath string) error {
	if !o.open {
		return ErrorScreenClosed
	}
	log.Printf("Mock screen is now displaying image \"%s\"", filepath)
	return nil
}

func (o *mockScreen) DisplayImage(reader io.Reader) error {
	if !o.open {
		return ErrorScreenClosed
	}
	log.Printf("Mock screen is now displaying image from the provided reader")
	return nil
}
