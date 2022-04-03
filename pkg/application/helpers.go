package application

import (
	"os"

	"fyne.io/fyne/v2/widget"
)

func quitButton(label string) *widget.Button {
	return widget.NewButton(label, func() {
		os.Exit(0)
	})
}
