//package ScheduleRepo
package ScheduleRepo

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Schedule"
	"github.com/jinzhu/gorm"
)

// ScheduleRepo struct
type ScheduleRepo struct {
	DB *gorm.DB
}

// NewScheduleRepo function
func NewScheduleRepo(db *gorm.DB) Schedule.ScheduleRepo {
	return &ScheduleRepo{
		DB: db,
	}
}

// ScheduleOfARound method
// func (shedulerepo *ScheduleRepo) ScheduleOfARound(RoundID uint)  {}
