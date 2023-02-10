package main

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func chooseDirectory(w fyne.Window, h *widget.Label) {
	d := dialog.NewFileOpen(func(fi fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if fi.URI() != nil {
			// fmt.Println(fi.URI().Path())
			w.Hide()
		}
		h.SetText(fi.URI().Path())
	}, w)
	d.Resize(fyne.NewSize(500, 650))
	d.Show()
	w.Show()
}

func isCSV(path string) bool {
	return filepath.Ext(path) == ".csv"
}

func main() {
	a := app.New()
	w := a.NewWindow("Adjusted Utilisation")

	// qtd := widget.NewLabel("")
	// wk := widget.NewLabel("")
	// tl := widget.NewLabel("")

	w2 := fyne.CurrentApp().NewWindow("dialog")

	wk := widget.NewLabel("Not selected")
	qtd := widget.NewLabel("Not selected")
	tl := widget.NewLabel("Not selectede")

	btn := widget.NewButton("Submit", func() {
		if !isCSV(wk.Text) {
			wk.SetText("Not a CSV file")
		}
		if !isCSV(qtd.Text) {
			qtd.SetText("Not a CSV file")
		}
		if !isCSV(tl.Text) {
			tl.SetText("Not a CSV file")
		}
		fmt.Println(wk.Text, qtd.Text, tl.Text)
	})

	w.SetContent(container.NewVBox(
		wk,
		widget.NewButton("Weekly Utilisation", func() {
			chooseDirectory(w2, wk)
		}),
		qtd,
		widget.NewButton("QTD Utilisation", func() {
			chooseDirectory(w2, qtd)
		}),
		tl,
		widget.NewButton("QTD Time Listing", func() {
			chooseDirectory(w2, tl)
		}),
		btn,
	))
	w2.Resize(fyne.NewSize(500, 650))
	w.SetMaster()
	w.ShowAndRun()
}
