package Section

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type SectionService interface {
	GetIDOfSectionsOfRound(RoundID uint) *[]uint
	DeleteSectionByID(SectionID uint) bool
	GetSectionByID(SectionID uint) *entity.Section
	GetSectionsOfRound(RoundID uint) *[]entity.Section
	DeleteScheduleDataUsingSectionID(RoundID uint) bool
	DeleteSectionsOfRound(RoundID uint) bool
	DeleteLecturesOfRound(RoundID uint) bool
}
