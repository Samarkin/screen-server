# samarkin/screen-server/oled

A simple go library to work with an SH1106 128x64 OLED screen.

Includes a custom 5x7 font.

## Usage

1. Ensure the display is properly connected over I2C.

2. Import the module:
```go
import "github.com/samarkin/screen-server/oled"
```

3. Profit:
```go
scr, err := oled.Open(&oled.I2cOpener{})
if err != nil {
    log.Fatalf("Failed to open screen: %v", err)
}
defer scr.Close()
scr.Print(0, 0, "Hello, world!")
```