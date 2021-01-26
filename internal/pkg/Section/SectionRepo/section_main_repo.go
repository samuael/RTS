package SectionRepo

import (
	"fmt"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Section"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/jinzhu/gorm"
)

// SectionRepo struct
type SectionRepo struct {
	DB *gorm.DB
}

// NewSectionRepo function
func NewSectionRepo(db *gorm.DB) Section.SectionRepo {
	return &SectionRepo{
		DB: db,
	}
}

// GetIDOfSectionsOfRound  method
func (sectionrepo *SectionRepo) GetIDOfSectionsOfRound(RoundID uint) (*[]uint, error) {
	sectionIDS := &[]uint{}
	fmt.Println("The Round ID ", RoundID)
	sections := &[]entity.Section{}
	newErra := sectionrepo.DB.Table("sections").Select("id").Where("round_refer=?", RoundID).Find(sections).Error
	for l := 0; l < len(*sections); l++ {
		*sectionIDS = append(*sectionIDS, (*sections)[l].ID)
	}
	return sectionIDS, newErra
}

// DeleteSectionByID method to return errro
func (sectionrepo *SectionRepo) DeleteSectionByID(SectionID uint) error {
	newError := sectionrepo.DB.Table("sections").Delete("id=?", SectionID).Error
	if newError != nil {
		return newError
	}
	return newError
}

// GetSectionByID method
func (sectionrepo *SectionRepo) GetSectionByID(SectionID uint) (*entity.Section, error) {
	section := &entity.Section{}
	newError := sectionrepo.DB.Where("id=?", SectionID).First(section).Error
	// Populating the section Datas
	fmt.Println("Populating The Sections ")
	newError = sectionrepo.DB.Model(section).Related(&section.ClassDates, "ClassDates").Error
	newError = sectionrepo.DB.Model(section).Related(&section.TrainingDates, "TrainingDates").Error
	newError = sectionrepo.DB.Model(section).Related(&section.Trainings, "Trainings").Error
	newError = sectionrepo.DB.Model(section).Related(&section.Lectures, "Lectures").Error
	newError = sectionrepo.DB.Model(section).Related(&section.Students, "Students").Error
	newError = sectionrepo.DB.Model(section).Related(&section.Round, "Round").Error
	if newError != nil {
		fmt.Println("  The Error While Populating section Is ", newError.Error())
	}
	return section, newError
}

// GetSectionsOfRound method
func (sectionrepo *SectionRepo) GetSectionsOfRound(RoundID uint) (*[]entity.Section, error) {
	sections := &[]entity.Section{}
	newErro := sectionrepo.DB.Where("round_refer=?", RoundID).Find(sections).Error
	for i := 0; i < len(*sections); i++ {
		section := &(*sections)[i]
		section.Room = &entity.Room{}
		newErro = sectionrepo.DB.Model(section).Related(&section.ClassDates, "ClassDates").Error
		newErro = sectionrepo.DB.Model(section).Related(&section.TrainingDates, "TrainingDates").Error
		newErro = sectionrepo.DB.Model(section).Related(&section.Trainings, "Trainings").Error
		newErro = sectionrepo.DB.Model(section).Related(&section.Lectures, "Lectures").Error
		newErro = sectionrepo.DB.Model(section).Related(&section.Students, "Students").Error
		newErro = sectionrepo.DB.Model(section).Related(&section.Round, "Round").Error
		newErro = sectionrepo.DB.Model(section).Related(section.Room, "RoomRefer").Error
		if newErro == nil && len(section.Lectures) > 0 {
			for i := 0; i < len(section.Lectures); i++ {
				lecture := &section.Lectures[i]
				newErro = sectionrepo.DB.Model(lecture).Related(&lecture.Course, "CourseRefer").Error
				newErro = sectionrepo.DB.Model(lecture).Related(&lecture.StartDate, "StartDateRefer").Error
				newErro = sectionrepo.DB.Model(lecture).Related(&lecture.EndDate, "EndDateRefer").Error
				newErro = sectionrepo.DB.Model(lecture).Related(&lecture.Teacher, "TeacherRefer").Error
			}

			for i := 0; i < len(section.Trainings); i++ {
				training := &section.Trainings[i]
				newErro = sectionrepo.DB.Model(training).Related(&training.StartDate, "StartDate").Error
				newErro = sectionrepo.DB.Model(training).Related(&training.EndDate, "EndDate").Error
				newErro = sectionrepo.DB.Model(training).Related(&training.Trainer, "Trainer").Error
				newErro = sectionrepo.DB.Model(training).Related(&training.Students, "Students").Error
			}
		}
	}
	return sections, newErro
}

// DeleteScheduleDataUsingSectionID metod to delete all the datas related to the schedule of the round id passed as a parameter
func (sectionrepo *SectionRepo) DeleteScheduleDataUsingSectionID(SectionID uint) error {
	fmt.Println("Inside This Funcition ..... ", SectionID)
	// For Trainings
	sectionTrainingDatesIDS := &[]uint{}

	sectionTrainingsIDS := &[]uint{}
	row, newError := sectionrepo.DB.Debug().Table("section_training_dates").Select("date_id").Where("section_id=?", SectionID).Rows()
	if newError != nil {
		fmt.Println(newError.Error())
		return newError
	}
	var era error

	for era == nil && row.Next() {
		val := 0
		era = row.Scan(&val)

		if era != nil {
			fmt.Println(era.Error())
			break
		}
		*sectionTrainingDatesIDS = append(*sectionTrainingDatesIDS, uint(val))
	}
	rows, newError := sectionrepo.DB.Debug().Table("section_trainings").Select("field_session_id").Where("section_id=?", SectionID).Rows()
	for era == nil && rows.Next() {
		val := 0
		era = rows.Scan(&val)
		if era != nil {
			fmt.Println(era.Error())
			break
		}
		*sectionTrainingsIDS = append(*sectionTrainingsIDS, uint(val))
	}
	// Now the IDS of training are fetched
	// Now the IDS of training dates are fetched
	newError = sectionrepo.DB.Debug().Table("dates").Delete(etc.Date{}, "id in (?)", *sectionTrainingDatesIDS).Error
	newError = sectionrepo.DB.Debug().Table("field_sessions").Delete(entity.FieldSession{}, "id in (?)", *sectionTrainingsIDS).Error

	// For Lectures Meaning Face To Face Lectures
	sectionLectureDateIDS := &[]uint{}
	sectionLectureIDS := &[]uint{}

	// Now I'm gonna delete the Data related to that (sections Lectures ID list ) and (Sections Lectures Date IDS List )
	sectinonclassdateidsrow, newError := sectionrepo.DB.Debug().Table("section_class_dates").Select("date_id").Where("section_id=?", SectionID).Rows()
	sectionlecturesidrows, newError := sectionrepo.DB.Debug().Table("section_lectures").Select("lecture_id").Where("section_id=?", SectionID).Rows()

	for sectinonclassdateidsrow.Next() {
		val := 0
		era = sectinonclassdateidsrow.Scan(&val)
		if era != nil {
			fmt.Println(era.Error())
			break
		}
		*sectionLectureDateIDS = append(*sectionLectureDateIDS, uint(val))
	}
	for sectionlecturesidrows.Next() {
		val := 0
		era = sectionlecturesidrows.Scan(&val)
		if era != nil {
			fmt.Println(era.Error())
			break
		}
		*sectionLectureIDS = append(*sectionLectureIDS, uint(val))
	}
	fmt.Println(*sectionLectureIDS)
	// Now I Am gonna Delete the Lectures and the Dates related to the Lectures
	newError = sectionrepo.DB.Debug().Table("dates").Where("id in (?)", *sectionLectureDateIDS).Delete(nil).Error
	newError = sectionrepo.DB.Debug().Table("lectures").Where("id in (?)", *sectionLectureIDS).Delete(nil).Error
	newError = sectionrepo.DB.Debug().Table("section_students").Where("section_id=?", SectionID).Delete(nil).Error
	if newError != nil {
		fmt.Println(newError.Error())
	}
	return newError
}

// DeleteSectionsOfRound method for deleting the Data of t he sections
func (sectionrepo *SectionRepo) DeleteSectionsOfRound(RoundID uint) error {
	newError := sectionrepo.DB.Table("sections").Where("round_refer=?", RoundID).Delete(nil).Error
	return newError
}

// DeleteLecturesOfRound method to delete the
func (sectionrepo *SectionRepo) DeleteLecturesOfRound(RoundID uint) error {
	newError := sectionrepo.DB.Table("lectures").Where("roundid=?", RoundID).Delete(nil).Error
	return newError
}
