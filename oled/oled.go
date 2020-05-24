package oled

// SignalLevels holds the number of supported signal levels
var SignalLevels int

// Screen contains resources required to work with the OLED screen
type Screen interface {
	// Print displays a string in the specified position of the screen
	Print(line int, offset int, message string) error
	// DisplaySignalLevel displays signal level icon in the specified position of the screen
	DisplaySignalLevel(line int, offset int, level int) error
	// DisplayImage loads image from the specified file and displays it on the screen
	DisplayImage(filepath string) error
	// Clear erases screen RAM contents
	Clear() error
	// Close releases all the resources allocated by this instance of Screen
	Close() error
}

// Opener opens a connection to the screen
type Opener interface {
	open() (Screen, error)
}

// Open is the entry point to start working with an OLED screen
func Open(o Opener) (Screen, error) {
	return o.open()
}
