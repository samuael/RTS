package BreakDateRepo

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/BreakDates"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/jinzhu/gorm"
)

// BreakDateRepo struct s
type BreakDateRepo struct {
	DB *gorm.DB
}

// NewBreakDateRepo function
func NewBreakDateRepo(db *gorm.DB) BreakDates.BreakDateRepo {
	return &BreakDateRepo{
		DB: db,
	}
}

// CreateBreakDate method
func (bdrepo *BreakDateRepo) CreateBreakDate(date *etc.Date) (*etc.Date, error) {
	date.IsBreakDate = true
	newError := bdrepo.DB.Save(date).Error
	return date, newError
}

// DeleteBreakDate method
func (bdrepo *BreakDateRepo) DeleteBreakDate(DateID uint) error {
	newErrro := bdrepo.DB.Table("dates").Where("id=?", DateID).Delete(&etc.Date{}, DateID).Error
	return newErrro
}

// GetBreakDates method
func (bdrepo *BreakDateRepo) GetBreakDates(BranchID uint, date *etc.Date) (*[]etc.Date, error) {
	breakdates := &[]etc.Date{}
	newErro := bdrepo.DB.Where("branch_id=? and is_break_date=?  and year >= ?", BranchID, true, date.Year).Find(breakdates).Error
	return breakdates, newErro
}
