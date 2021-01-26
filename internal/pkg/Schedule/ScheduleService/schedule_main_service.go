package ScheduleService

import (
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Schedule"
)

// ScheduleService struct
type ScheduleService struct {
	ScheduleRepo Schedule.ScheduleRepo
}

// NewScheduleService function
func NewScheduleService(Schedulerepo Schedule.ScheduleRepo) Schedule.ScheduleService {
	return &ScheduleService{
		ScheduleRepo: Schedulerepo,
	}
}
