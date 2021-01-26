package Section

import "github.com/Projects/RidingTrainingSystem/internal/pkg/entity"

type SectionRepo interface {
	GetIDOfSectionsOfRound(RoundID uint) (*[]uint, error)
	DeleteSectionByID(SectionID uint) error
	GetSectionByID(SectionID uint) (*entity.Section, error)
	GetSectionsOfRound(SectionID uint) (*[]entity.Section, error)
	DeleteScheduleDataUsingSectionID(SectionID uint) error
	DeleteSectionsOfRound(RoundID uint) error
	DeleteLecturesOfRound(RoundID uint) error
	// GetRoundByIDForSchedule(RoundID uint) (*entity.Round, error)
}
