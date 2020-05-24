package engine

import (
	"fmt"
	"log"
	"strings"

	"github.com/samarkin/screen-server/oled"
)

// Engine is a singleton object to manage the screen
type Engine interface {
	Connected() bool
	GetLastError() string
	Clear()
	GetMessage(line int) string
	DisplayMessage(text string, line int)
	ClearMessage(line int)
	AppendMessage(text string)
	Shutdown()
}

var instance Engine

// GetEngine instantiates a new, or returns an existing Engine instance
// Use Engine.Connected() to see if screen has been connected successfully
// TODO: Thread-safety
func GetEngine() Engine {
	if instance == nil {
		e := &engine{}
		e.scr, e.lastError = oled.Open(&oled.I2cOpener{})
		instance = e
	}
	return instance
}

type engine struct {
	scr        oled.Screen
	lastError  error
	messages   [8]string
	cursorLine int
}

func (e *engine) Connected() bool {
	return e.scr != nil
}

func (e *engine) GetLastError() string {
	return e.lastError.Error()
}

func (e *engine) Clear() {
	log.Printf("Clearing screen...")
	if e.scr == nil {
		e.lastError = fmt.Errorf("Screen not connected")
		return
	}
	for i := range e.messages {
		e.messages[i] = ""
	}
	e.lastError = e.scr.Clear()
}

func (e *engine) ClearMessage(line int) {
	log.Printf("Clearing message on line %d...", line)
	if e.scr == nil {
		e.lastError = fmt.Errorf("Screen not connected")
		return
	}
	if line >= 0 && line < 8 {
		e.messages[line] = ""
	}
	e.lastError = e.scr.Print(line, 0, strings.Repeat(" ", 21))
}

func (e *engine) AppendMessage(text string) {
	e.DisplayMessage(text, e.cursorLine)
	e.cursorLine = (e.cursorLine + 1) & 0x07
}

func (e *engine) DisplayMessage(text string, line int) {
	log.Printf("Displaying message \"%s\" on line %d...", text, line)
	if e.scr == nil {
		e.lastError = fmt.Errorf("Screen not connected")
		return
	}
	if line >= 0 && line < 8 {
		e.messages[line] = text
	}
	e.lastError = e.scr.Print(line, 0, text)
}

func (e *engine) GetMessage(line int) string {
	if line >= 0 && line < 8 {
		return e.messages[line]
	}
	return ""
}

func (e *engine) Shutdown() {
	log.Printf("Shutting down...")
	if e.scr != nil {
		if err := e.scr.Clear(); err != nil {
			e.scr.Print(0, 0, "Shutting down...")
		}
		e.scr.Close()
	}
}
