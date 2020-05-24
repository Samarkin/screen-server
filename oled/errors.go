package oled

type oledError string

const (
	// ErrorScreenClosed means an attempt to communicate with a closed OLED screen
	ErrorScreenClosed = oledError("Screen is closed")
)

func (e oledError) Error() string {
	return string(e)
}
