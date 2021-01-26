package Payment

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

type PaymentService interface {
	CreatePayment(payment *entity.Payment) *entity.Payment
	PaymentsOfRound(Roundid, BranchNumber uint) *[]entity.Payment
	PaymentsOfStudent(StudentID uint) *[]entity.Payment
	PaymentsOfSecretary(SecretaryID, limit, offset uint) *[]entity.Payment
	StudentPaidAmount(StudentID uint) float64
	RoundPaidAmount(Roundid uint) float64
	GetSinglePaymentReciptData(PaymentID uint) *entity.SinglePaymentDataStructure
	GetBranchPaymentReport(BranchID uint, date *etc.Date) *[]entity.Payment
	GetDailyPaymentReport(date etc.Date, BranchID uint) *entity.DailyPaymentResult
	DPR(date etc.Date, BranchID int)
	//
	// PaymentsOfEachSecretariesForMonth(BranchID, Year, Month uint) map[*entity.Admin]float64
	// PaymentsOfRoundForTheMonth(BranchID, Year, Month uint) map[*entity.Round]float64
	// PaymentsPfCategoryForTheMonth(BranchID, Year, Month uint) map[*entity.Category]float64
	// TotalPaymentofMonth(BranchID, Year, Month uint) float64
}

// 84049930  Rambo
