package BreakDates

import (
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

type BreakDateService interface {
	CreateBreakDate(date *etc.Date) *etc.Date
	DeleteBreakDate(ID uint) bool
	GetBreakDates(BranchID uint, date *etc.Date) *[]etc.Date
}
