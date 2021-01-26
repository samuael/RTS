package Payment

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

type PaymentRepo interface {
	CreatePayment(payment *entity.Payment) (*entity.Payment, error)
	PaymentsOfRound(Roundid, BranchNumber uint) (*[]entity.Payment, error)
	PaymentsOfStudent(StudentID uint) (*[]entity.Payment, error)
	PaymentsOfSecretary(SecretaryID, limit, offset uint) (*[]entity.Payment, error)
	StudentPaidAmount(StudentID uint) (float64, error)
	RoundPaidAmount(Roundid uint) (float64, error)
	GetPaymentByID(PaymentID uint) (*entity.Payment, error)
	GetBranchPaymentReport(BranchID uint, date *etc.Date) (*[]entity.Payment, error)
	GetDailyAdminsPaymentReport(date etc.Date, BranchID uint) (map[*entity.Admin]float64, error)
	// GetDailyRoundsPaymentReport(date etc.Date, BranchID uint) (map[*entity.Round]float64, error)
	// GetDailyCategoryPaymentReport(date etc.Date, BranchID uint) (map[*entity.Category]float64, error)

	// PaymentsOfEachSecretariesForMonth(BranchID, Year, Month uint) (map[*entity.Admin]float64, error)
	// PaymentsOfRoundForTheMonth(BranchID, Year, Month uint) (map[*entity.Round]float64, error)
	// PaymentsPfCategoryForTheMonth(BranchID, Year, Month uint) (map[*entity.Category]float64, error)
	// TotalPaymentofMonth(BranchID, Year, Month uint) (float64, error)
}
