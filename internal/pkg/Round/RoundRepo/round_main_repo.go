package RoundRepo

import (
	"fmt"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/jinzhu/gorm"
)

// RoundRepo struct
type RoundRepo struct {
	DB *gorm.DB
}

// NewRoundRepo function
func NewRoundRepo(db *gorm.DB) Round.RoundRepo {
	return &RoundRepo{
		DB: db,
	}
}

// GetRoundByID method returning a round from gorm repo
func (roundrepo *RoundRepo) GetRoundByID(ID uint) (*entity.Round, error) {
	newRound := &entity.Round{}
	theError := roundrepo.DB.Table("rounds").Where("id=?", ID).First(newRound).Error
	// Here Comes the Population of the Datas of the Round Niggoye
	roundrepo.DB.Model(newRound).Related(&newRound.Category, "CategoryRefer")
	roundrepo.DB.Model(newRound).Related(&newRound.Courses, "Courses")
	roundrepo.DB.Model(newRound).Related(&newRound.Lectures, "Lectures")
	roundrepo.DB.Model(newRound).Related(&newRound.Sections, "Sections")
	roundrepo.DB.Model(newRound).Related(&newRound.Trainers, "Trainers")
	roundrepo.DB.Model(newRound).Related(&newRound.Teachers, "Teachers")
	roundrepo.DB.Model(newRound).Related(&newRound.Students, "Students")
	roundrepo.DB.Model(newRound).Related(&newRound.CreatedBY, "AdminRefer")
	return newRound, theError
}

// SaveRound emthod to save a round
func (roundrepo *RoundRepo) SaveRound(round *entity.Round) (*entity.Round, error) {
	newEror := roundrepo.DB.Save(round).Debug().Error
	newEror = roundrepo.DB.Model(&entity.Round{}).Where("id=?", round.ID).Updates(map[string]interface{}{
		"max_students":      round.MaxStudents,
		"studentscount":     len(round.Students),
		"learning":          round.Learning,
		"on_registration":   round.OnRegistration,
		"training_duration": round.TrainingDuration,
		"cost":              round.Cost,
	}).Error

	return round, newEror
}

// CreateRound method
func (roundrepo *RoundRepo) CreateRound(round *entity.Round) (*entity.Round, error) {
	newError := roundrepo.DB.Create(round).Error
	return round, newError
}

// GetRounds (BranchID uint) *[]entity.Round
func (roundrepo *RoundRepo) GetRounds(BranchID uint) (*[]entity.Round, error) {
	rounds := &[]entity.Round{}
	newError := roundrepo.DB.Where("branchnumber=?", BranchID).Find(rounds).Error
	for index := 0; index < len(*rounds); index++ {
		round := &((*rounds)[index])
		// Populating each round
		roundrepo.DB.Model(round).Related(&round.Courses, "Courses")
		roundrepo.DB.Model(round).Related(&round.Students, "Students")
		roundrepo.DB.Model(round).Related(&round.Teachers, "Teachers")
		roundrepo.DB.Model(round).Related(&round.Trainers, "Trainers")
		roundrepo.DB.Model(round).Related(&round.Sections, "Sections")
		roundrepo.DB.Model(round).Related(&round.Lectures, "Lectures")
		roundrepo.DB.Model(round).Related(&round.CreatedBY, "AdminRefer")
		round.AdminRefer = round.CreatedBY.ID
		round.Studentscount = uint(len(round.Students))
		(*rounds)[index] = *round
	}
	return rounds, newError
}

// DeleteRoundLectires By

// DeleteFieldDates   (RoundID uint) error
func (roundrepo *RoundRepo) DeleteFieldDates(RoundID uint) error {
	datesID := &[]uint{}
	newError := roundrepo.DB.Table("field_dates").Select("date_id").
		Joins("left join dates on dates.id=field_dates.date_id").
		Where("round_id=?", RoundID).
		Scan(datesID).
		Error
	newError = roundrepo.DB.Table("dates").Delete(*datesID).Error
	fmt.Println(datesID)
	newError = roundrepo.DB.Table("field_dates").Where("date_id IN (?)", *datesID).Delete(nil).Error
	// fmt.Println(newError.Error())
	return newError
}

// DeleteLectureDates   (RoundID uint) error
func (roundrepo *RoundRepo) DeleteLectureDates(RoundID uint) error {
	teachersBussydatesID := &[]uint{}
	newError := roundrepo.DB.Table("lectures_bussy_date").Select("date_id").
		Joins("left join dates on dates.id=lectures_bussy_date.date_id").
		Where("round_id=?", RoundID).
		Scan(teachersBussydatesID).Error
	newError = roundrepo.DB.Table("dates").Delete(*teachersBussydatesID).Error
	newError = roundrepo.DB.
		Table("lectures_bussy_date").
		Where("date_id in (?)", *teachersBussydatesID).
		Delete(nil).
		Error
	return newError
}

// DeleteFieldSessionAndItsStudents (RoundID uint) error
func (roundrepo *RoundRepo) DeleteFieldSessionAndItsStudents(RoundID uint) error {
	fieldSessions := &[]entity.FieldSession{}
	newError := roundrepo.DB.Table("field_sessions").Where("round_refer=?", RoundID).Select([]string{"id", "start_date_id", "end_date_id"}).Find(fieldSessions).Error
	if newError == nil {
		for _, session := range *fieldSessions {
			newError = roundrepo.DB.Table("field_student").Where("field_session_id=?", session.ID).Delete(nil).Error
			newError = roundrepo.DB.
				Table("field_sessions").
				Where("id=?", session.ID).
				Delete(nil).
				Error
			newError = roundrepo.DB.Table("dates").Delete([]uint{session.StartDateID, session.EndDateID}).Error
		}
	}
	return newError
}

// DeleteLecturesAndItsStudents (RoundID uint) error
func (roundrepo *RoundRepo) DeleteLecturesAndItsStudents(RoundID uint) error {
	lectures := &[]entity.Lecture{}
	newError := roundrepo.DB.Table("lectures").Where("roundid=?", RoundID).Select([]string{"id", "start_date_refer", "end_date_refer"}).Find(lectures).Error
	if newError == nil {
		for _, session := range *lectures {
			newError = roundrepo.DB.Table("lectures").Where("id=?", session.ID).Delete(nil).Error
			newError = roundrepo.DB.Table("dates").Delete([]uint{session.StartDateRefer, session.EndDateRefer}).Error
		}
	}
	return newError
}

// DeleteRoundRelatedRoomsDates (RoundID uint) error
func (roundrepo *RoundRepo) DeleteRoundRelatedRoomsDates(RoundID uint) error {
	roomtoDateids := &[]uint{}
	newError := roundrepo.DB.
		Table("room_to_date").
		Select("date_id").
		Joins("left join dates on dates.id=room_to_date.date_id").
		Where("round_id=?", RoundID).
		Scan(roomtoDateids).Error
	newError = roundrepo.DB.Table("dates").Delete(*roomtoDateids).Error
	newError = roundrepo.DB.
		Table("room_to_date").
		Where("date_id in (?)", *roomtoDateids).
		Delete(nil).Error
	return newError
}

// GetRoundByIDForSchedule (RoundID uint) (*entity.Round , error)
func (roundrepo *RoundRepo) GetRoundByIDForSchedule(RoundID uint) (*entity.Round, error) {
	newRound := &entity.Round{}
	theError := roundrepo.DB.Table("rounds").Where("id=?", RoundID).First(newRound).Error
	// Here Comes the Population of the Datas of the Round Niggoye
	roundrepo.DB.Model(newRound).Related(&newRound.Category, "CategoryRefer")
	roundrepo.DB.Model(newRound).Related(&newRound.Courses, "Courses")
	roundrepo.DB.Model(newRound).Related(&newRound.Lectures, "Lectures")
	roundrepo.DB.Model(newRound).Related(&newRound.Sections, "Sections")
	roundrepo.DB.Model(newRound).Related(&newRound.Trainers, "Trainers")
	roundrepo.DB.Model(newRound).Related(&newRound.Teachers, "Teachers")
	roundrepo.DB.Model(newRound).Related(&newRound.Students, "Students")
	roundrepo.DB.Model(newRound).Related(&newRound.CreatedBY, "AdminRefer")

	// Populating the Round Sections
	for k := 0; k < len(newRound.Sections); k++ {
		roundrepo.DB.Model(&newRound.Sections[k]).Related(&newRound.Sections[k].Trainings, "Trainings")
		roundrepo.DB.Model(&newRound.Sections[k]).Related(&newRound.Sections[k].TrainingDates, "TrainingDates")
		roundrepo.DB.Model(&newRound.Sections[k]).Related(&newRound.Sections[k].Lectures, "Lectures")
		roundrepo.DB.Model(&newRound.Sections[k]).Related(&newRound.Sections[k].Students, "Students")
		roundrepo.DB.Model(&newRound.Sections[k]).Related(&newRound.Sections[k].ClassDates, "ClassDates")
	}
	// Populating the Trainers of the Round
	for k := 0; k < len(newRound.Trainers); k++ {
		theError = roundrepo.DB.Model(newRound.Trainers[k]).Related(&newRound.Trainers[k].BusyDates, "BusyDates").Error
		theError = roundrepo.DB.Model(newRound.Trainers[k]).Related(&newRound.Trainers[k].Categoty, "Categoty").Error
		theError = roundrepo.DB.Model(newRound.Trainers[k]).Related(&newRound.Trainers[k].Vehicle, "Vehicle").Error
	}

	// Populating the Teachers of the Round
	for k := 0; k < len(newRound.Teachers); k++ {
		theError = roundrepo.DB.Model(&newRound.Teachers[k]).Related(&newRound.Trainers[k].BusyDates, "BusyDates").Error
	}

	// Populating the Teachers of the Round
	for k := 0; k < len(newRound.Lectures); k++ {
		theError = roundrepo.DB.Model(&newRound.Lectures[k]).Related(&newRound.Lectures[k].Course, "Course").Error
	}
	return newRound, theError
}

// DeleteSectionsAndRelated methods to delete the Datas related to Section
// func (roundrepo *RoundRepo) DeleteSectionsAndRelated(RoundID uint) error {
// 	newError := roundrepo.DB.Table("sections").Delete("round_refer=?", RoundID).Error
// 	return newError
// }

// UpdateToRegistration  (RoundID uint) error
func (roundrepo *RoundRepo) UpdateToRegistration(RoundID uint) error {
	newErro := roundrepo.
		DB.
		Table("rounds").
		Debug().
		Where("id=?", RoundID).
		Update(map[string]bool{"on_registration": true, "learning": false}).
		Error
	return newErro
}

// GetActiveRoundsOfCategory (CategoryID uint) (*[]entity.Round, error)
func (roundrepo *RoundRepo) GetActiveRoundsOfCategory(CategoryID uint) (*[]entity.Round, error) {
	rounds := &[]entity.Round{}
	newEra := roundrepo.DB.Table("rounds").Find(rounds, "category_refer=? and active=?", CategoryID, true).Error
	return rounds, newEra
}

// GetRoundsOfCategory (CategoryID uint) (*[]entity.Round, error)
func (roundrepo *RoundRepo) GetRoundsOfCategory(CategoryID uint) (*[]entity.Round, error) {
	rounds := &[]entity.Round{}
	newEra := roundrepo.DB.Table("rounds").Find(rounds, "category_refer=?", CategoryID).Error
	return rounds, newEra
}

// IsRoundNumberReseerved method returning integer the number of rows affected Having BranchID , RoundNumber , CategoryID
func (roundrepo *RoundRepo) IsRoundNumberReseerved(BranchID, CategoryID, RoundNumber uint) uint {
	count := roundrepo.DB.Where("branchnumber=?  and roundnumber=? and category_refer=?", BranchID, RoundNumber, CategoryID).RowsAffected
	return uint(count)
}

// UpdateStudentsCount (uint) error
func (roundrepo *RoundRepo) UpdateStudentsCount(RoundID, Studentscount uint) error {
	row := roundrepo.DB.Table("rounds").Where("id=?", RoundID).Select("studentscount").Row()
	studentscount := 0
	era := row.Scan(&Studentscount)
	if era != nil {
		return era
	}
	era = roundrepo.DB.Table("rounds").Where("id=?", RoundID).Updates(map[string]int{"studentscount": studentscount + int(Studentscount)}).Error
	return era
}
