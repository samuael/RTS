package BreakDateService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/BreakDates"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

// BreakDateService struct
type BreakDateService struct {
	BreakDateRepo BreakDates.BreakDateRepo
}

// NewBreakDateService function
func NewBreakDateService(breakRepo BreakDates.BreakDateRepo) BreakDates.BreakDateService {
	return &BreakDateService{
		BreakDateRepo: breakRepo,
	}
}

// CreateBreakDate function
func (bdservice *BreakDateService) CreateBreakDate(date *etc.Date) *etc.Date {
	newDate, erra := bdservice.BreakDateRepo.CreateBreakDate(date)
	if erra != nil {
		return nil
	}
	return newDate
}

// DeleteBreakDate method
func (bdservice *BreakDateService) DeleteBreakDate(DateID uint) bool {
	erra := bdservice.BreakDateRepo.DeleteBreakDate(DateID)
	if erra != nil {
		return false
	}
	return true
}

// GetBreakDates method
func (bdservice *BreakDateService) GetBreakDates(BranchID uint, date *etc.Date) *[]etc.Date {
	dates, erra := bdservice.BreakDateRepo.GetBreakDates(BranchID, date)
	if erra != nil {
		return nil
	}
	return dates
}
