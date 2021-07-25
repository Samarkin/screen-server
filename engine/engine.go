package engine

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/samarkin/screen-server/oled"
)

// Engine is a singleton object to manage the screen
type Engine interface {
	Connected() bool
	Clear() error
	GetMessage(line int) string
	DisplayMessage(text string, line int) error
	DisplayTemporaryMessage(text string, line int, timeout time.Duration) error
	DisplayImage(reader io.Reader) error
	DisplayTemporaryImage(reader io.Reader, duration time.Duration) error
	ClearMessage(line int) error
	AppendMessage(text string) error
	Shutdown()
}

var instanceMutex = &sync.Mutex{}
var instance Engine
var initializationError error

var padding = strings.Repeat(" ", 21)
var distantFuture = time.Now().AddDate(10, 0, 0) // 10 years from now
const smallDelay = 10 * time.Millisecond

// GetEngine instantiates a new, or returns an existing Engine instance
// Use Engine.Connected() to see if screen has been connected successfully
func GetEngine() (Engine, error) {
	instanceMutex.Lock()
	defer instanceMutex.Unlock()
	if instance == nil {
		e := &engine{}
		e.mutex = &sync.Mutex{}
		e.scr, initializationError = oled.Open(&oled.I2cOpener{})
		instance = e
	}
	return instance, initializationError
}

type message struct {
	text       string
	expiration time.Time
}

type engine struct {
	mutex      *sync.Mutex
	scr        oled.Screen
	messages   [8]message
	cursorLine int
}

func (e *engine) Connected() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.scr != nil
}

func (e *engine) Clear() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	log.Printf("Clearing screen...")
	for i := range e.messages {
		e.messages[i] = message{"", distantFuture}
	}
	if e.scr == nil {
		return fmt.Errorf("screen not connected")
	}
	return e.scr.Clear()
}

func (e *engine) ClearMessage(line int) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	log.Printf("Clearing message on line %d...", line)
	if line >= 0 && line < 8 {
		e.messages[line] = message{"", distantFuture}
	}
	if e.scr == nil {
		return fmt.Errorf("screen not connected")
	}
	return e.scr.Print(line, 0, padding)
}

func (e *engine) AppendMessage(text string) error {
	e.mutex.Lock()
	cursorLine := e.cursorLine
	e.cursorLine = (cursorLine + 1) & 0x07
	e.mutex.Unlock()
	return e.DisplayMessage(text, cursorLine)
}

func (e *engine) DisplayMessage(text string, line int) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	log.Printf("Displaying message \"%s\" on line %d...", text, line)
	if line >= 0 && line < 8 {
		e.messages[line] = message{text, distantFuture}
	}
	if e.scr == nil {
		return fmt.Errorf("screen not connected")
	}
	return e.scr.Print(line, 0, text+padding)
}

func (e *engine) DisplayTemporaryMessage(text string, line int, duration time.Duration) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	log.Printf("Displaying message \"%s\" on line %d for %s...", text, line, duration)
	if line >= 0 && line < 8 {
		e.messages[line] = message{text, time.Now().Add(duration)}
		go func() {
			time.Sleep(duration + smallDelay)
			e.mutex.Lock()
			defer e.mutex.Unlock()
			if time.Now().After(e.messages[line].expiration) {
				log.Printf("Erasing message on line %d...", line)
				e.messages[line] = message{"", distantFuture}
				if e.scr != nil {
					e.scr.Print(line, 0, padding)
				}
			}
		}()
	}
	if e.scr == nil {
		return fmt.Errorf("screen not connected")
	}
	return e.scr.Print(line, 0, text+padding)
}

func (e *engine) DisplayImage(reader io.Reader) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for i := range e.messages {
		e.messages[i] = message{"<IMAGE>", distantFuture}
	}
	if e.scr == nil {
		return fmt.Errorf("screen not connected")
	}
	return e.scr.DisplayImage(reader)
}

func (e *engine) DisplayTemporaryImage(reader io.Reader, duration time.Duration) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	expiration := time.Now().Add(duration)
	for i := range e.messages {
		e.messages[i] = message{"<IMAGE>", expiration}
	}
	go func() {
		time.Sleep(duration + smallDelay)
		e.mutex.Lock()
		defer e.mutex.Unlock()
		log.Printf("Erasing image")
		now := time.Now()
		for i := range e.messages {
			if now.After(e.messages[i].expiration) {
				e.messages[i] = message{"", distantFuture}
				if e.scr != nil {
					e.scr.Print(i, 0, padding)
				}
			}
		}
	}()
	if e.scr == nil {
		return fmt.Errorf("screen not connected")
	}
	return e.scr.DisplayImage(reader)
}

func (e *engine) GetMessage(line int) string {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if line >= 0 && line < 8 {
		return e.messages[line].text
	}
	return ""
}

func (e *engine) Shutdown() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	log.Printf("Shutting down...")
	if e.scr != nil {
		if err := e.scr.Clear(); err != nil {
			e.scr.Print(0, 0, "Shutting down...")
		}
		e.scr.Close()
		e.scr = nil
	}
}
