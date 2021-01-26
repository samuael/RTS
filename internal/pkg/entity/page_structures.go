package entity

import (
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/Projects/RidingTrainingSystem/pkg/form"
)

// ControllPageStructure representing Controll page info
type ControllPageStructure struct {
	Branch     Branch
	Admin      Admin
	Categories []Category
	Rooms      []Room
	Host       string
}

// AdminRegistrationsPageStructure struct
type AdminRegistrationsPageStructure struct {
	Host       string
	Branch     Branch
	Admin      Admin
	CSRF       string
	VErrors    form.ValidationErrors
	Categories []Category
}

// Four04 struct
type Four04 struct {
	Admin      *Admin
	Branch     *Branch
	Path       string
	StatusCode int
	Statustext string
}

// ScheduleDataStructure struct
type ScheduleDataStructure struct {
	Message  string
	Host     string
	Branch   Branch
	Round    Round
	Sections []Section
	Success  bool
}

// SinglePaymentDataStructure struct
type SinglePaymentDataStructure struct {
	Success bool
	Message string
	Host    string
	Branch  Branch
	Payment Payment
}

// ActiveLecturesListDataStructure struct
type ActiveLecturesListDataStructure struct {
	Success  bool
	Message  string
	Lectures []Lecture
}

// AddingVehicleToFieldMan struct
type AddingVehicleToFieldMan struct {
	FieldmanID uint `json:"trainer_id"`
	VehicleID  uint `json:"vehicle_id"`
}

// MonthlyPaymentResult struct
type MonthlyPaymentResult struct {
	Year                  uint
	Month                 uint
	SecretaryToPaidAmount map[*Admin]float64
	RoundToPaidAmount     map[*Round]float64
	CategoryToPaidAmount  map[*Category]float64
	TotalPaidAmount       float64
}

// DailyPaymentResult struct
type DailyPaymentResult struct {
	Date   etc.Date
	Branch Branch
	Today  etc.Date
	// Year                  uint
	// Month                 uint
	SecretaryToPaidAmount map[*Admin]float64
	RoundToPaidAmount     map[*Round]float64
	CategoryToPaidAmount  map[*Category]float64
	TotalPaidAmount       float64
}

// PaymentPageData struct
type PaymentPageData struct {
	Success           bool
	Message           string
	SecretaryPayments map[*Admin]float64
	TotalPayment      float64
	Branch            Branch
}
