package BreakDates

import etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"

type BreakDateRepo interface {
	CreateBreakDate(date *etc.Date) (*etc.Date, error)
	DeleteBreakDate(ID uint) error
	GetBreakDates(BranchID uint, date *etc.Date) (*[]etc.Date, error)
}
