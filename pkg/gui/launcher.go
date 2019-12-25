package gui

import (
	"fyne.io/fyne"
	fyneapp "fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	app "github.com/fernandreu/spotlight/pkg"
	"github.com/fernandreu/spotlight/pkg/img"
)

func LaunchGui(options app.Options) {
	result := app.ProcessFiles(options.Source, options.Destination)

	application := fyneapp.New()

	picsContainer := fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithColumns(3))
	fillContainer(picsContainer, result.NewFiles)
	picsContainer.AddObject(layout.NewSpacer())

	destBox := widget.NewHBox(
		widget.NewLabel("Destination:"),
		widget.NewLabel(result.Destination))
	quitButton := widget.NewButton("Quit", func() {
		application.Quit()
	})

	w := application.NewWindow("Windows Spotlight Fetcher")
	w.SetContent(fyne.NewContainerWithLayout(
		layout.NewBorderLayout(destBox, quitButton, nil, nil),
		destBox,
		quitButton,
		widget.NewGroupWithScroller("New Pictures", picsContainer)))

	w.ShowAndRun()
}

func fillContainer(container *fyne.Container, images []img.ImageFile) {

	for _, image := range images {
		item := &canvas.Image{
			File:     image.FullPath(),
			FillMode: canvas.ImageFillContain,
		}
		item.SetMinSize(fyne.Size{
			Width:  150,
			Height: 100,
		})
		container.AddObject(item)
	}
}
