package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	appWindow := myApp.NewWindow("Tomate")

	countdown_label := widget.NewLabel("TO BE SET")

	appWindow.SetContent(
		container.New(layout.NewCenterLayout(),
			container.New(layout.NewVBoxLayout(),
				countdown_label,
				widget.NewButton("Start", func() {
				}),
			),
		),
	)

	appWindow.ShowAndRun()
}
