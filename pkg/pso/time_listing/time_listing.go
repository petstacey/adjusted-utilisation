package time_listing

import (
	"errors"
	"os"

	"github.com/gocarina/gocsv"
)

type TimeListing struct {
	Company             string  `json:"company" csv:"Company"`
	Engagement          string  `json:"engagement" csv:"Engagement"`
	Project             string  `json:"project" csv:"project"`
	TaskName            string  `json:"taskName" csv:"Task Name"`
	CombinedName        string  `json:"combinedName" csv:"Combined Name"`
	GlobalWorkgroup     string  `json:"globalWorkgroup" csv:"Global Workgroup"`
	Workgroup           string  `json:"workgroup" csv:"Workgroup"`
	ResourceName        string  `json:"resourceName" csv:"Resource"`
	TimeEntryDate       string  `json:"timeEntryDate" csv:"Time Entry Date"`
	CreatedOn           string  `json:"createdOn" csv:"Created On"`
	TimeEntryStatus     string  `json:"timeEntryStatus" csv:"Time Entry Status"`
	TimeEntryType       string  `json:"timeEntryType" csv:"Time Entry Type"`
	RequestProject      string  `json:"requestProject" csv:"Request Project"`
	StartTime           string  `json:"startTime" csv:"Start Time"`
	EndTime             string  `json:"endTime" csv:"End Time"`
	FixedFee            string  `json:"fixedFee" csv:"Fixed Fee"`
	FixedFeeDeliverable string  `json:"fixedFeeDeliverable" csv:"Fixed Fee Deliverable"`
	Recognized          string  `json:"recognized" csv:"Recognized"`
	Prepaid             string  `json:"prepaid" csv:"prepaid"`
	Billable            string  `json:"billable" csv:"Billable"`
	Capital             string  `json:"capital" csv:"Capital"`
	Utilized            string  `json:"utilized" csv:"Utilized"`
	RegularHours        float64 `json:"regularHours" csv:"Regular Hours"`
	OvertimeHours       float64 `json:"otHours" csv:"OT Hours"`
	TimeCode1           string  `json:"timeCode1" csv:"Time Code 1"`
	TimeCode2           string  `json:"timeCode2" csv:"Time Code 2"`
	TimeCode3           string  `json:"timeCode3" csv:"Time Code 3"`
	Description         string  `json:"description" csv:"Description"`
	WorkLocation        string  `json:"workLocation" csv:"Work Location"`
	WorkCode            string  `json:"workCode" csv:"Work Code"`
	Locked              string  `json:"locked" csv:"Locked"`
	Interfaced          string  `json:"interfaced" csv:"Interfaced"`
	CurrentStatus       string  `json:"currentStatus" csv:"Current Status"`
}

type TimeListingReport struct {
	Filename string
	Rows     []*TimeListing
}

func NewTimeListingReport(filename string) *TimeListingReport {
	return &TimeListingReport{
		Filename: filename,
		Rows:     []*TimeListing{},
	}
}

func (u *TimeListingReport) ReadTimeListing() error {
	file, err := os.OpenFile(u.Filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = gocsv.UnmarshalFile(file, &u.Rows); err != nil {
		return err
	}
	return nil
}

func (u *TimeListingReport) TimeListingByResource(names []string) (map[string][]TimeListing, error) {
	if len(u.Rows) == 0 {
		return nil, errors.New("cannot generate report from no rows")
	}
	result := make(map[string][]TimeListing)
	for _, name := range names {
		result[name] = u.TimeListingForResource(name)
	}
	return result, nil
}

func (u *TimeListingReport) TimeListingForResource(name string) []TimeListing {
	result := []TimeListing{}
	for _, listing := range u.Rows {
		if listing.ResourceName == name {
			result = append(result, *listing)
		}
	}
	return result
}

func (u *TimeListingReport) NonBillableProjectHours(times []TimeListing) float64 {
	var hours float64
	for _, row := range times {
		if row.TimeEntryType == "project" && row.Billable == "No" {
			hours += row.RegularHours
			hours += row.OvertimeHours
		}
	}
	return hours
}

func (u *TimeListingReport) HoursForField(field string, times []TimeListing) float64 {
	var hours float64
	for _, row := range times {
		if row.TaskName == field {
			hours += row.RegularHours
			hours += row.OvertimeHours
		}
	}
	return hours
}
