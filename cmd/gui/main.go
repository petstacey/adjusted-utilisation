package main

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/gocarina/gocsv"
	tl "github.com/pso-dev/utilisation/pkg/pso/time_listing"
	ute "github.com/pso-dev/utilisation/pkg/pso/utilization"
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
	d.Resize(fyne.NewSize(500, 600))
	d.Show()
	w.Show()
}

func isCSV(path string) bool {
	return filepath.Ext(path) == ".csv"
}

func main() {
	a := app.New()
	w := a.NewWindow("Adjusted Utilisation")

	w2 := fyne.CurrentApp().NewWindow("dialog")

	wkUte := widget.NewLabel("Not selected")
	qtdUte := widget.NewLabel("Not selected")
	wkTL := widget.NewLabel("Not selected")
	qtdTL := widget.NewLabel("Not selected")
	errorLbl := widget.NewLabel("")

	btn := widget.NewButton("Submit", func() {
		errorLbl.SetText("")
		if !isCSV(wkUte.Text) {
			wkUte.SetText("Not a recognised (UTF-8) CSV file")
			errorLbl.SetText("One or more files are not of the correct type!")
		}
		if !isCSV(qtdUte.Text) {
			qtdUte.SetText("Not a recognised (UTF-8) CSV file")
			errorLbl.SetText("One or more files are not of the correct type!")
		}
		if !isCSV(wkTL.Text) {
			wkTL.SetText("Not a recognised (UTF-8) CSV file")
			errorLbl.SetText("One or more files are not of the correct type!")
		}
		if !isCSV(qtdTL.Text) {
			qtdTL.SetText("Not a recognised (UTF-8) CSV file")
			errorLbl.SetText("one or more files are not of the correct type!")
		}
		errorLbl.SetText("Generating...")
		err := generateReports(wkUte.Text, qtdUte.Text, wkTL.Text, qtdTL.Text)
		if err != nil {
			errorLbl.SetText(err.Error())
		} else {
			errorLbl.SetText("Done!")
		}
	})

	w.SetContent(container.NewVBox(
		widget.NewButton("Weekly Utilisation", func() {
			chooseDirectory(w2, wkUte)
		}),
		wkUte,
		widget.NewButton("QTD Utilisation", func() {
			chooseDirectory(w2, qtdUte)
		}),
		qtdUte,
		widget.NewButton("Weekly Time Listing", func() {
			chooseDirectory(w2, wkTL)
		}),
		wkTL,
		widget.NewButton("QTD Time Listing", func() {
			chooseDirectory(w2, qtdTL)
		}),
		qtdTL,
		btn,
		errorLbl,
	))
	w.Resize(fyne.NewSize(500, 400))
	w2.Resize(fyne.NewSize(550, 650))
	w.SetMaster()
	w.ShowAndRun()
}

func generateReports(wkUte, qtdUte, wkTL, qtdTL string) error {
	wkUteRep := ute.NewUtilizationReport(wkUte)
	wkTLRep := tl.NewTimeListingReport(wkTL)

	if err := wkUteRep.ReadUtilization(); err != nil {
		return err
	}

	if err := wkTLRep.ReadTimeListing(); err != nil {
		return err
	}

	weekReport, err := ute.GenerateAdjustedUtilization(wkUteRep, wkTLRep)
	if err != nil {
		return err
	}

	qtdUteRep := ute.NewUtilizationReport(qtdUte)
	qtdTLRep := tl.NewTimeListingReport(qtdTL)

	if err := qtdUteRep.ReadUtilization(); err != nil {
		return err
	}

	if err := qtdTLRep.ReadTimeListing(); err != nil {
		return err
	}

	qtdReport, err := ute.GenerateAdjustedUtilization(qtdUteRep, qtdTLRep)
	if err != nil {
		return err
	}

	if err := SaveReport(weekReport, "WeekAdjusted.csv"); err != nil {
		return err
	}

	if err := SaveReport(qtdReport, "QTDAdjusted.csv"); err != nil {
		return err
	}

	return nil
}

func SaveReport(report *ute.AdjustedUtilizationReport, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	err = gocsv.MarshalFile(&report.Rows, file)
	if err != nil {
		return err
	}
	return nil
}
