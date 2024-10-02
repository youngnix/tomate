package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	appWindow := myApp.NewWindow("Tomate")

	appWindow.SetContent(widget.NewLabel("Hello, World!"))
	appWindow.ShowAndRun()
}
