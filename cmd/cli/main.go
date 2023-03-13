package main

import (
	"fmt"
	"os"

	tl "github.com/pso-dev/utilisation/pkg/pso/time_listing"
	ute "github.com/pso-dev/utilisation/pkg/pso/utilization"
)

func main() {
	u := ute.NewUtilizationReport("UtilizationGrossReport.csv")
	if err := u.ReadUtilization(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	names := u.GetNames()
	fmt.Println(len(names))

	tListing := tl.NewTimeListingReport("Time Listing Report.csv")
	if err := tListing.ReadTimeListing(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(len(tListing.Rows))

	report, err := ute.GenerateAdjustedUtilization(u, tListing)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, row := range report.Rows {
		fmt.Println(row)
	}
}
