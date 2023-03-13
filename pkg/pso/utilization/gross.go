package utilization

import (
	"os"

	"github.com/gocarina/gocsv"
)

type Utilization struct {
	Department                     string  `json:"geoDepartment" csv:"GEO / Department"`
	Workgroup                      string  `json:"workgroup" csv:"Workgroup Name"`
	ResourceName                   string  `json:"resourceName" csv:"Resource Name"`
	Email                          string  `json:"email" csv:"Email Address"`
	Region                         string  `json:"region" csv:"Region"`
	SubRegion                      string  `json:"subRegion" csv:"Sub Region"`
	Practice                       string  `json:"practice" csv:"BU / Practice"`
	Function                       string  `json:"function" csv:"Primary Function"`
	ResourceType                   string  `json:"resourceType" csv:"Resource Type"`
	AvailableHoursGross            float64 `json:"availableHoursGross" csv:"Available Hours (Gross)"`
	BillableHours                  float64 `json:"billableHours" csv:"Billable Hours"`
	BillableUtilizationGross       string  `json:"billableUtilizationGross" csv:"Billable Utilization (Gross)"`
	ProductiveHours                float64 `json:"productiveHours" csv:"Productive Hours (Gross)"`
	ProductiveUtilizationGross     string  `json:"productiveUtilizationGross" csv:"Productive Utilization (Gross)"`
	InternalFundedHours            float64 `json:"internalFundedHours" csv:"Internal Funded Hours"`
	InternalFundedUtilizationGross string  `json:"InternalFundedUtilizationGross" csv:"Internal Funded Utilization (Gross)"`
	ChargeableHours                float64 `json:"chargeableHours" csv:"Chargeable Hours"`
	ChargeableUtilizationGross     string  `json:"chargeableUtilizationGross" csv:"Chargeable Utilization (Gross)"`
	LOAHours                       float64 `json:"loaHours" csv:"LOA Hours"`
	UnapprovedBillableHours        float64 `json:"unapprovedBillableHours" csv:"Unapproved Billable Hours"`
	UnapprovedProductiveHours      float64 `json:"unapprovedProductiveHours" csv:"Unapproved Productive Hours"`
	Manager                        string  `json:"manager" csv:"Manager"`
}

type UtilizationReport struct {
	Filename string
	Rows     []*Utilization
}

func NewUtilizationReport(filename string) *UtilizationReport {
	return &UtilizationReport{
		Filename: filename,
		Rows:     []*Utilization{},
	}
}

func (u *UtilizationReport) ReadUtilization() error {
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

func (u *UtilizationReport) GetNames() []string {
	names := []string{}
	for _, row := range u.Rows {
		names = append(names, row.ResourceName)
	}
	return names
}

func (u *UtilizationReport) GetForName(name string) *Utilization {
	for _, row := range u.Rows {
		if row.ResourceName == name {
			return row
		}
	}
	return nil
}
