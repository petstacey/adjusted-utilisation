package utilization

import (
	"errors"
	"fmt"

	tl "github.com/pso-dev/utilisation/pkg/pso/time_listing"
)

type AdjustedUtilization struct {
	ResourceName             string  `json:"resourceName" csv:"Resource Name"`
	Workgroup                string  `json:"Workgroup" csv:"Workgroup"`
	ResourceType             string  `json:"resourceType" csv:"Resource Type"`
	AvailableHours           float64 `json:"availableHours" csv:"Available Hours"`
	BillableHours            float64 `json:"billableHours" csv:"Billable Hours"`
	NonBillableProjectHours  float64 `json:"nonBillableProjectHours" csv:"Non-Billable Project Hours"`
	PaidTimeOff              float64 `json:"pto" csv:"PTO"`
	PublicHoliday            float64 `json:"publicHoliday" csv:"PublicHoliday"`
	SickTime                 float64 `json:"sickTime" csv:"Sick Time"`
	BereavementLeave         float64 `json:"bereavement" csv:"Bereavement Leave"`
	TimeOffWithoutPay        float64 `json:"timeOffWithoutPay" csv:"Time-off Without Pay"`
	SalesSupport             float64 `json:"salesSupport" csv:"Sales Support"`
	Mentoring                float64 `json:"mentoring" csv:"Mentoring"`
	Travel                   float64 `json:"travel" csv:"Travel"`
	Training                 float64 `json:"trainingAndDevelopment" csv:"Training & Development"`
	PSOPMO                   float64 `json:"psoPmo" csv:"PSO-PMO"`
	Administration           float64 `json:"administration" csv:"Administration"`
	BillableUtilizationGross float64 `json:"billableUtilizationGross" csv:"Billable Utilization (Gross)"`
	AdjustedUtilization      float64 `json:"adjustedUtilization" csv:"Adjusted Utilization"`
}

type AdjustedUtilizationReport struct {
	Rows []*AdjustedUtilization
}

func GenerateAdjustedUtilization(ute *UtilizationReport, tl *tl.TimeListingReport) (*AdjustedUtilizationReport, error) {
	if ute == nil || tl == nil {
		return nil, errors.New("Utilization Report or Time Listing Report cannot be nil")
	}

	if len(ute.Rows) == 0 || len(tl.Rows) == 0 {
		return nil, errors.New("Utilization Report or Time Listing Report cannot be 0 length")
	}

	filteredTimes, err := tl.TimeListingByResource(ute.GetNames())
	if err != nil {
		return nil, err
	}

	adjustedReport := AdjustedUtilizationReport{
		Rows: []*AdjustedUtilization{},
	}

	for name, times := range filteredTimes {
		if len(times) == 0 {
			return nil, fmt.Errorf("times not recorded for %s", name)
		}
		var adjusted AdjustedUtilization
		gross := ute.GetForName(name)
		adjusted.ResourceName = name
		adjusted.Workgroup = times[0].Workgroup
		adjusted.ResourceType = gross.ResourceType
		adjusted.AvailableHours = gross.AvailableHoursGross
		adjusted.BillableHours = gross.BillableHours
		adjusted.NonBillableProjectHours = tl.NonBillableProjectHours(times)
		adjusted.PaidTimeOff = PermanentOnly("PTO", gross, tl, times)
		adjusted.PublicHoliday = PermanentOnly("Public Holiday", gross, tl, times)
		adjusted.SickTime = PermanentOnly("Sick Time", gross, tl, times)
		adjusted.BereavementLeave = PermanentOnly("Bereavement", gross, tl, times)
		adjusted.TimeOffWithoutPay = tl.HoursForField("Time off without pay", times)
		adjusted.SalesSupport = PermanentOnly("Sales Support", gross, tl, times)
		adjusted.Mentoring = tl.HoursForField("Mentoring", times)
		adjusted.Travel = tl.HoursForField("Travel", times)
		adjusted.Training = tl.HoursForField("Training and development", times)
		adjusted.PSOPMO = tl.HoursForField("PSO - PMO", times)
		adjusted.Administration = tl.HoursForField("Administration", times)
		adjusted.BillableUtilizationGross = BillableUtilizationGross(gross.BillableHours, gross.AvailableHoursGross)
		if gross.AvailableHoursGross <= 0 {
			adjusted.AdjustedUtilization = 0.0
		} else {
			adjusted.AdjustedUtilization = (adjusted.BillableHours + adjusted.NonBillableProjectHours + adjusted.PaidTimeOff + adjusted.PublicHoliday + adjusted.SickTime + adjusted.BereavementLeave + adjusted.TimeOffWithoutPay + adjusted.SalesSupport + adjusted.Mentoring + adjusted.Travel) / gross.AvailableHoursGross
		}
		adjustedReport.Rows = append(adjustedReport.Rows, &adjusted)
	}

	return &adjustedReport, nil
}

func PermanentOnly(field string, gross *Utilization, tl *tl.TimeListingReport, times []tl.TimeListing) float64 {
	if gross.ResourceType != "Full time" {
		return 0.0
	}
	return tl.HoursForField(field, times)
}

func BillableUtilizationGross(billableHours, availableHours float64) float64 {
	if availableHours <= 0 {
		return 0.0
	}
	return billableHours / availableHours
}
