package SectionService

import (
	// "fmt"
	"fmt"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Section"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

// SectionService struct
type SectionService struct {
	SectionRepo Section.SectionRepo
}

// NewSectionService function
func NewSectionService(SectionRepo Section.SectionRepo) Section.SectionService {
	return &SectionService{
		SectionRepo: SectionRepo,
	}
}

// GetIDOfSectionsOfRound method
func (sectionser *SectionService) GetIDOfSectionsOfRound(RoundID uint) *[]uint {
	ids, errors := sectionser.SectionRepo.GetIDOfSectionsOfRound(RoundID)
	if errors != nil {
		return nil
	}
	return ids
}

// DeleteSectionByID   (SectionID uint) error
func (sectionser *SectionService) DeleteSectionByID(SectionID uint) bool {
	newError := sectionser.SectionRepo.DeleteSectionByID(SectionID)
	if newError != nil {
		return false
	}
	return true
}

//

// GetSectionByID method returning section Pointer and error
func (sectionser *SectionService) GetSectionByID(SectionID uint) *entity.Section {
	section, errror := sectionser.SectionRepo.GetSectionByID(SectionID)
	if errror != nil {
		return nil
	}
	return section
}

// GetSectionsOfRound  (RoundID uint) *[]entity.Section
func (sectionser *SectionService) GetSectionsOfRound(RoundID uint) *[]entity.Section {
	sections, err := sectionser.SectionRepo.GetSectionsOfRound(RoundID)
	if err != nil {
		return nil
	}
	return sections
}

// DeleteScheduleDataUsingSectionID method for deleting all the datas related to the this rounds schedule
func (sectionser *SectionService) DeleteScheduleDataUsingSectionID(SectionID uint) bool {
	fmt.Println("Inside Section Service ... ")
	newErra := sectionser.SectionRepo.DeleteScheduleDataUsingSectionID(SectionID)
	fmt.Println("Calling Sectionn Reposr DeleteScheduleDataUsingsectionID not Succesful  ... ")

	if newErra != nil {
		return false
	}
	return true
}

// DeleteSectionsOfRound (RoundID uint) bool
func (sectionser *SectionService) DeleteSectionsOfRound(RoundID uint) bool {
	newError := sectionser.SectionRepo.DeleteSectionsOfRound(RoundID)
	if newError != nil {
		return false
	}
	return true
}

// DeleteLecturesOfRound (RoundID uint) bool
func (sectionser *SectionService) DeleteLecturesOfRound(RoundID uint) bool {
	newError := sectionser.SectionRepo.DeleteLecturesOfRound(RoundID)
	if newError != nil {
		fmt.Println(newError.Error())
		return false
	}
	return true
}
