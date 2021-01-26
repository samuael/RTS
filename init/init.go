package main

import (
	"fmt"
	"log"

	db "github.com/Projects/RidingTrainingSystem/internal/pkg/Db"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/jinzhu/gorm"
)

func main() {
	initialize()
}

var student *entity.Student

func initialize() {
	autoMigrate()
	// PopulateDataBase()
}

var dbs *gorm.DB

func autoMigrate() {
	dbs, erro := db.InitializPostgres()
	if erro != nil {
		log.Fatal(erro)
	}
	defer dbs.Close()
	dbs.AutoMigrate(
		entity.Payment{}, //
		entity.CourseToDuration{},
		entity.Address{}, //
		entity.Student{}, //
		entity.Teacher{},
		entity.FieldAssistant{},
		entity.Admin{},
		entity.Course{},
		entity.Category{},
		entity.Branch{},
		entity.Resource{},
		entity.Room{},
		entity.Section{},
		entity.Lecture{},
		entity.Question{},
		entity.AskedQuetion{},
		entity.FieldSession{},
		entity.Round{},
		entity.ActiveRound{},
		entity.Vehicle{},
		etc.Date{},
		entity.SectionTraining{},
		entity.SectionTrainingDate{},
		entity.SectionLecture{},
		entity.SectionStudent{},
		entity.SectionClassDate{},
		entity.Information{},
		entity.FieldDate{},
		// entity.RoomToDate{},
		//SectionTraining
		//SectionTrainingDate
		// SectionLecture
		// SectionStudent
		// SectionClassDate
		// RoomToDate
	)
	err := dbs.GetErrors()
	if err != nil {
		log.Println(err)
	}
}

// PopulateDataBase  function
func PopulateDataBase() {
	dbs, erro := db.InitializPostgres()
	if erro != nil {
		log.Fatal(erro)
	}
	defer dbs.Close()
	fmt.Println("Populating the branch ")
	branch := entity.Branch{
		Email:         "samuaeladnew@gmail.com",
		City:          "Assosa",
		BranchAcronym: "SHA",
		Logourl:       "img/logo.png",
		Name:          "Shambel Drivir Training Institute ",
		Country:       "Ethiopia ",
		Address: entity.Address{
			Kebele: "04",
			Region: "Benishangul Gumz",
			Woreda: "Assosa",
			Zone:   "assosa ",
		},
		LicenceGivenDate:      *etc.NewDate(0),
		BranchFullnameAmharic: " Shambel Drivers Training System",
		BranchFullnameEnglish: " Shambel Drivers Training System",
		// LicenceDate:           "10/10/2005",
		LicenceNumber: "78787",
		Moto:          " Alamachin Himumann Merdat New ",
	}
	admin := entity.Admin{
		// gorm.Model: gorm.Model{ID: 5},
		Branch:    branch,
		Firstname: "Samuael",
		Lastname:  "Adnew ",
		Lang:      "en",
		Password:  "samuaelfirst",
		Phone:     "+251992078204",
		Uploader: entity.Uploader{
			Username: "samuael/1/2012",
			Imageurl: "img/logo.png",
		},
		Role: entity.SUPERADMIN,
	}
	admin.ID = 1
	admin.Branch.Createdby = admin.Username
	new := dbs.Save(&admin).Error
	if new != nil {
		fmt.Println(new.Error())
	}
	admino := &entity.Admin{}
	brancho := &entity.Branch{}
	address := &entity.Address{}
	rn := dbs.Last(admino)
	if rn != nil {
		fmt.Println(rn)
	}
	er := dbs.Model(admino).Related(brancho, "BranchRefer").Debug().Error
	if er != nil {
		fmt.Println(er)
	}
	dbs.Model(branch).Related(address, "OwnerID").Debug()
	fmt.Print("\n\n\n\n\n", address, "\n\n\n\n\n")
	branch.Address = *address
	fmt.Println("n\n\n\n\n", "This is the Addresss \n", address, "\n\n")
	erroa := dbs.Model(branch).Related(branch.LicenceGivenDate, "LicenceGivenDate").Debug().Error
	if erroa != nil {
		fmt.Println("Error While Finding t he DAte ", branch.LicenceGivenDate)
		// return
	}
	// branch.LicenceGivenDate = dateet
	admino.Branch = *brancho

	// fmt.Println(admin)
	fmt.Println(admino.ID, admino.Username, admino.Password, admino.Phone, admino.Branch)
	// fmt.Println(brancho)
	category := &entity.Category{
		Branchid:                 1,
		Title:                    "Bajaj",
		Rounds:                   []entity.Round{},
		ImageURL:                 "img/logo.png",
		NumberOfStudentsLearning: 23,
	}
	category.ID = 1
	categorys := &entity.Category{}
	vehicle := &entity.Vehicle{
		BoardNumber: "1234",
		Branch:      branch,
		Category:    *category,
		CategoryID:  category.ID,
		BranchNo:    branch.ID,
		Imageurl:    "breadcrumb-bg.jpg",
	}
	trainer := entity.FieldAssistant{
		Categoty:     *category,
		CategoryID:   category.ID,
		Email:        "samuaeladnew@gmail.com",
		Createdby:    admin.Username,
		Firstname:    "samuael",
		Lastname:     "lastname ",
		GrandName:    "Birhane ",
		Username:     "Trainer/samuael/0",
		Imageurl:     "img/logo.png",
		Lang:         "amh",
		Password:     "1234",
		Phonenumber:  "0992948594",
		BranchNumber: 1,
		// BoardNumber:  "1234",
		Vehicle:   *vehicle,
		VehicleID: vehicle.ID,
	}

	teacher := &entity.Teacher{
		Active:       true,
		BranchNumber: 1,
		BusyDates:    []etc.Date{},
		Createdby:    admin.Username,
		Email:        "teacher@gmail.com",
		Firstname:    "Alemu",
		Lastname:     "Kelemu",
		GrandName:    "Ayana",
		Imageurl:     "img/logo.png",
		Lang:         "amh",
		Password:     "2323",
		Phonenumber:  "0999994447",
		Username:     "Teacher/1/2012",
	}

	trainer.ID = 1
	// vehicle.Trainer = &trainer
	vehicle.TrainerID = trainer.ID
	round1 := entity.Round{
		Branchnumber:   1,
		Category:       *category,
		CategoryRefer:  category.ID,
		CreatedBY:      admin,
		AdminRefer:     admin.ID,
		Cost:           2000,
		OnRegistration: true,
		Year:           2012,
		Roundnumber:    1,
	}

	dbs.Save(&round1)
	dbs.Save(teacher)
	// dbs.Save(&trainer)
	dbs.Save(&vehicle)
	dbs.First(categorys, 1)
	fmt.Println("Category Thing \n", categorys)
}
// drop table payments , course_to_durations , address, students , sessions , langs , teachers , field_assistants  , admins ,cources ,categorys , branchs , resources  , rooms ,  sections , lectures , questions  , asked_quetions , field_sessions , rounds , active_rounds
