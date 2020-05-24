package oled

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strings"

	"golang.org/x/exp/io/i2c"
)

// I2cOpener allows to open a real I2C screen
type I2cOpener struct {
}

type i2cScreen struct {
	dev *i2c.Device
}

func (o *I2cOpener) open() (Screen, error) {
	dev, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x3c)
	if err != nil {
		return nil, err
	}

	screen := &i2cScreen{dev: dev}
	if err := screen.init(); err != nil {
		dev.Close()
		return nil, err
	}
	return screen, nil
}

func (s *i2cScreen) init() error {
	if err := s.dev.Write([]byte{0x00, 0xAE}); err != nil {
		return fmt.Errorf("Failed to turn off: %v", err)
	}
	if err := s.dev.Write([]byte{0x00, 0xA1}); err != nil {
		return fmt.Errorf("Failed to rotate: %v", err)
	}
	if err := s.dev.Write([]byte{0x00, 0xC8}); err != nil {
		return fmt.Errorf("Failed to flip: %v", err)
	}
	if err := s.dev.Write([]byte{0x00, 0x40}); err != nil {
		return fmt.Errorf("Failed to set offset: %v", err)
	}
	if err := s.Clear(); err != nil {
		return fmt.Errorf("Failed to clean: %v", err)
	}
	if err := s.dev.Write([]byte{0x00, 0xAF}); err != nil {
		return fmt.Errorf("Failed to turn on: %v", err)
	}
	return nil
}

func (s *i2cScreen) Close() error {
	return s.dev.Close()
}

func (s *i2cScreen) Clear() error {
	const width = 132
	emptyLine := make([]byte, width+1)
	for i := range emptyLine {
		emptyLine[i] = 0x00
	}
	emptyLine[0] = 0x40
	for i := 0; i < 8; i++ {
		if err := s.dev.Write([]byte{0x00, 0xB0 + byte(i&0x7), 0x02, 0x10}); err != nil {
			return err
		}
		if err := s.dev.Write(emptyLine); err != nil {
			return err
		}
	}
	return nil
}

func (s *i2cScreen) Print(line int, offset int, message string) error {
	if err := s.dev.Write([]byte{0x00, 0xB0 | byte(line&0x7), byte((offset + 2) & 0x07), 0x10 | byte(((offset+2)>>4)&0x07)}); err != nil {
		return fmt.Errorf("Failed to set page and offset: %v", err)
	}
	if len(message) > 21 {
		message = message[:21]
	}
	for _, ch := range strings.ToUpper(message) {
		var letter []byte
		if ch < ' ' || int(ch-' ') >= len(font) {
			letter = font[len(font)-1]
		} else {
			letter = font[ch-' ']
		}
		if err := s.dev.Write(append(append([]byte{0x40}, letter...), 0x00)); err != nil {
			return fmt.Errorf("Failed to print %c: %v", ch, err)
		}
	}
	return nil
}

func (s *i2cScreen) DisplaySignalLevel(line int, offset int, level int) error {
	if err := s.dev.Write([]byte{0x00, 0xB0 | byte(line&0x7), byte((offset + 2) & 0x07), 0x10 | byte(((offset+2)>>4)&0x07)}); err != nil {
		return fmt.Errorf("Failed to set page and offset: %v", err)
	}
	if level >= len(signalLevels) {
		level = len(signalLevels) - 1
	}
	if level < 0 {
		level = 0
	}
	if err := s.dev.Write(append([]byte{0x40}, signalLevels[level]...)); err != nil {
		return fmt.Errorf("Failed to display signal level: %v", err)
	}
	return nil
}

func (s *i2cScreen) DisplayImage(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return err
	}
	rect := img.Bounds()
	if rect.Dx() != 128 || rect.Dy() != 64 {
		return fmt.Errorf("Image should have size 128x64")
	}
	var i byte
	for y := rect.Min.Y; y < rect.Max.Y; y += 8 {
		if err := s.dev.Write([]byte{0x00, 0xB0 + i, 0x02, 0x10}); err != nil {
			return fmt.Errorf("Failed to set page: %v", err)
		}
		var ts []byte
		for x := rect.Min.X; x < rect.Max.X; x++ {
			var t byte
			for yy := 7; yy >= 0; yy-- {
				c := color.GrayModel.Convert(img.At(x, y+yy)).(color.Gray)
				t <<= 1
				if c.Y < 0x80 {
					t |= 1
				}
			}
			ts = append(ts, t)
		}
		if err := s.dev.Write(append(append([]byte{0x40}, ts...), 0x00)); err != nil {
			return fmt.Errorf("Failed to output: %v", err)
		}
		i++
	}
	return nil
}
