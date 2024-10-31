package terminal

const (
	Black = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

func (o *output) SetColor(colorCode int) {
	o.Printf(csi+"38;5;%d"+"m", colorCode)
}

func (o *output) ResetColor() {
	o.Print(csi + "0m")
}
