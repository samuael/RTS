package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/BreakDates"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Lecture"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Room"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Section"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Schedule"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/Projects/RidingTrainingSystem/pkg/Helper"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

// ScheduleHandler strict
type ScheduleHandler struct {
	ScheduleService  Schedule.ScheduleService
	SessionHandler   *session.Cookiehandler
	RoundService     Round.RoundService
	RoomService      Room.RoomService
	SectionService   Section.SectionService
	LectureService   Lecture.LectureService
	TrainerService   Trainer.TrainerService
	TeacherService   Teacher.TeacherService
	BreakDateService BreakDates.BreakDateService
}

// NewScheduleHandler function
func NewScheduleHandler(
	scheduleservice Schedule.ScheduleService,
	sessionser *session.Cookiehandler,
	RoundService Round.RoundService,
	RoomService Room.RoomService,
	SectionService Section.SectionService,
	// LectureService Lecture.LectureService,
	FieldSession Trainer.TrainerService,
	TeacherService Teacher.TeacherService,
	BreakDateService BreakDates.BreakDateService,
) *ScheduleHandler {
	return &ScheduleHandler{
		ScheduleService: scheduleservice,
		SessionHandler:  sessionser,
		RoundService:    RoundService,
		RoomService:     RoomService,
		SectionService:  SectionService,
		// LectureService:  LectureService,
		TrainerService:   FieldSession,
		TeacherService:   TeacherService,
		BreakDateService: BreakDateService,
	}
}

// GenerateScheduleForRound method  the method is get request
func (schedulehandler *ScheduleHandler) GenerateScheduleForRound(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success bool
		Message string
		Round   entity.Round
	}{
		Success: false,
	}
	roundidstring := request.FormValue("round_id")
	mode := request.FormValue("mode")
	startday := request.FormValue("day")
	startmonth := request.FormValue("month")
	startyear := request.FormValue("year")
	day, era := strconv.Atoi(startday)
	month, era := strconv.Atoi(startmonth)
	year, era := strconv.Atoi(startyear)
	roundid, erra := strconv.Atoi(roundidstring)
	trainingPerDay, trainingPerDayError := strconv.Atoi(request.FormValue("tpd"))
	if trainingPerDayError != nil || trainingPerDay < 1 || trainingPerDay > 2 {
		ds.Message = "Invalid Training Per Date value"
		ds.Success = false
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
	}
	now := etc.NewDate(0)
	if era != nil || erra != nil || day <= 0 || day > 30 || month <= 0 || month >= 13 || year > now.Year {
		// Error Specific Niggoye
		ds.Message = "Invalid Input Values "
		ds.Success = false
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
	}
	startdate := &etc.Date{
		Year:  year,
		Month: month,
		Day:   day,
	}
	today := etc.NewDate(0)
	if !startdate.IsFuture(today) {
		ds.Message = "Invalid Date Start Date "
		ds.Success = false
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	startdate = startdate.FulFill()
	if strings.ToLower(mode) == entity.HALF {
		round, message := schedulehandler.GenerateScheduleHalfMode(uint(roundid), startdate, uint(trainingPerDay))
		ds.Message = message
		if round == nil {
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		round = schedulehandler.RoundService.SaveRound(round)
		ds.Success = true
		ds.Round = *round
		ds.Message = "Succsfukky Generate Half mode "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
	} else if strings.ToLower(mode) == entity.FULL { //GenerateScheduleFullMode
		round, message := schedulehandler.GenerateScheduleFullMode(uint(roundid), startdate, uint(trainingPerDay))
		ds.Message = message
		if round == nil {
			jsonReturn, _ := json.Marshal(ds)
			response.Write(jsonReturn)
			return
		}
		round = schedulehandler.RoundService.SaveRound(round)
		ds.Success = true
		ds.Round = *round
		ds.Message = "Succsfukky Generate Full mode "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
	}
}

// GenerateScheduleHalfMode method
func (schedulehandler *ScheduleHandler) GenerateScheduleHalfMode(RoundID uint, date *etc.Date, trainingPerDay uint) (*entity.Round, string) {
	round := schedulehandler.RoundService.GetRoundByIDForSchedule(RoundID)
	message := ""
	date = date.Modify()
	fmt.Println(
		"Learning :", round.Learning,
		"OnRegistration : ", round.OnRegistration,
	)
	rooms := schedulehandler.RoomService.GetRoomsOfABranch(round.Branchnumber)
	fmt.Println("The Lengtth of Dates The Room Is Reserved Are ", len((*rooms)[0].ReservedDates))
	if round.Learning {
		message = "The Round Has A Schedule"
		return round, message
	}
	if round == nil {
		return nil, "No Rond By Corresponding ID "
	}
	if rooms == nil {
		message = "There Is No Room To Assign"
		return nil, message
	}
	var teachers *[]entity.Teacher
	if len(round.Teachers) > 0 {
		teachers = &round.Teachers
	} else {
		// Trainers of the Category
		teachers = (schedulehandler.TeacherService.TeachersOfBranchByID(round.Branchnumber))
		if &teachers == nil {
			message = " No Teachers For Corresponding Branch "
			return nil, message
		}
	}
	var trainers *[]entity.FieldAssistant
	if len(round.Trainers) > 0 {
		trainers = &round.Trainers
	} else {
		// Trainers of the Category
		trainers = (schedulehandler.TrainerService.TrainersOfCategoryID(round.CategoryRefer))
		if trainers == nil {
			message = "Not Trainer of Corresponding Category "
			return nil, message
		}
	}

	breakdates := schedulehandler.BreakDateService.GetBreakDates(round.Branchnumber, etc.NewDate(0))
	courses := round.Courses // []entity.Course
	var lectures []entity.Lecture
	coursesDuration := 0
	for _, course := range courses {
		coursesDuration += int(course.Duration)
	}
	var passed bool
	freeRooms := []*entity.Room{}
	for d := 0; d < len(*rooms); d++ {
		rom := &(*rooms)[d]
		checkDate := *date
		checkDate.Hour = 2
		checkDate.Minute = 30
		checkDate = *(checkDate.FulFill())
		passed = true // the room is elegible if the it passes this days gize looping and Comparing
		for Helper.IsDateIn(checkDate, *breakdates) || checkDate.Name == etc.Ehud {
			checkDate = *(checkDate.NextDay())
		}
		for i := 0; i < (coursesDuration / 3); i++ {
			for _, day := range rom.ReservedDates {
				if day.Year == checkDate.Year && day.Month == checkDate.Month && day.Day == checkDate.Day && day.SubName == checkDate.SubName {
					passed = false
					break
				}
			}
			if passed == false {
				break
			} else {
				checkDate = *(checkDate.NextTime())
			}
		}
		if passed {
			freeRooms = append(freeRooms, rom)
		}
	}
	if len(freeRooms) == 0 {
		return nil, "Error While Selecting Rooms Please Select Date where All The Rooms Are Free"
	}
	rooma := (Helper.OrderRooms(freeRooms))
	selectedRooms := Helper.SelectedRoomsHalfMode(rooma, uint(len(round.Students))) // Type []*entity.Room
	totalCapacity := 0
	for _, ros := range selectedRooms {
		totalCapacity += int(ros.Capacity)
	}
	// log.Println("The Selected Rooms are  , ", len(*selectedRooms), " And The Total Capacity is : ", totalCapacity, "\n\nWhile the first one Has Length Of ", (*selectedRooms)[0].Capacity)
	sections := &[]entity.Section{} // Creating sections for the Round
	initIndex := 0
	for i := 0; i < len((selectedRooms)); i++ {
		roome := (selectedRooms)[i]
		capacity := roome.Capacity
		dirsha := int(math.Round(float64(len(round.Students)) * float64(float64(capacity)/float64(totalCapacity))))
		section := entity.Section{
			Categoryid:  round.CategoryRefer,
			Room:        roome, // The Variable Room is Pointer to A room not a room Value
			RoomRefer:   roome.ID,
			Sectionname: "Section " + strconv.Itoa(i+1),
			RoundRefer:  round.ID,
			// SectionCourseToDurations: []*entity.CourseToDuration{},
		}
		section.Students = round.Students[initIndex : initIndex+dirsha]
		*sections = append((*sections), section)
		initIndex += (dirsha)
	}
	// Generating Lectures and FieldSessions for each Section
	beginDate := *date
	beginDate.Hour = 2
	beginDate.Minute = 15
	// Generating Lectures
	coursetoduration := map[uint]uint{} //  entity.CourseToDuration{
	for i := 0; i < len(courses); i++ {
		for {
			course := courses[i]
			elapsed := coursetoduration[course.ID]
			if elapsed < course.Duration {
				remain := course.Duration - elapsed
				lectu := entity.Lecture{
					Course:      course,
					Branchid:    round.Branchnumber,
					CourseRefer: course.ID,
					Roundid:     round.ID,
				}
				if remain > entity.MaxLectureDuration {
					lectu.Duration = uint(entity.MaxLectureDuration)
					coursetoduration[course.ID] += uint(entity.MaxLectureDuration)
				} else {
					lectu.Duration = remain
					coursetoduration[course.ID] += uint(remain)
				}
				lectures = append(lectures, lectu)
				round.Lectures = append(round.Lectures, lectu)
			} else {
				break
			}
		}
	}
	// Generating Schedule For Each Section
	//In this Case it Will Be Once a Day
	for i := 0; i < len(*sections); i++ {
		Dato := *date
		newDate := &Dato
		newDate.Hour = 2
		newDate.Minute = 15
		classes := 0
		// duration := 0
		var teacher *entity.Teacher
		section := &((*sections)[i])
		room := section.Room /// Here it is giving the Pointer to the room variable
		// if len(room.ReservedDates) > 0 {
		// 	newDate = *(newDate.NextTime())
		// }
		freeDatesoTheRoom := []etc.Date{
			// This free dates can also represent dates that are free the Room of the Section Could Have
		}
		// This counter counts the Nummber of Counts where the variable Next Time is Called
		// InHalf Mode the Number of times the nextDay calling has to be after two times calling of nextTime
		oneRemain := false
		for _, lecture := range lectures {
			newLecture := lecture
			if classes == 0 {
				duration := lecture.Course.Duration
				classes = int(duration / entity.MaxLectureDuration) // Classes represents the Number of lectures Could A lecture Could HAve
				if duration%entity.MaxLectureDuration > 0 {
					classes++
				}
				// Generaate Free Dates of room to the Lecture of Specific Course to be Thought
				for k := 0; k < int(classes); k++ {
					// newDate = *(newDate.Modify())
					if oneRemain {
						newDate.RoundID = round.ID
						room.ReservedDates = append(room.ReservedDates, *newDate)
						freeDatesoTheRoom = append(freeDatesoTheRoom, *newDate)
						oneRemain = false
						newDate = newDate.NextDay()
						continue
					}
					for Helper.IsReservedHalf(*newDate, room.ReservedDates) || Helper.IsDateIn(*newDate, *breakdates) || newDate.Name == etc.Ehud {
						if newDate.Name == etc.Ehud || Helper.IsDateIn(*newDate, *breakdates) {
							newDate.SubName = etc.KESEAT
							newDate.Shift = etc.ShiftD
						}
						newDate = (newDate.NextTimeHalf())
					}
					newDate.RoundID = round.ID
					if ((newDate.Shift != etc.ShiftD) || (newDate.Shift != etc.ShiftA)) && (classes-k > 1) {
						room.ReservedDates = append(room.ReservedDates, *newDate)
						freeDatesoTheRoom = append(freeDatesoTheRoom, *newDate)

						newDate = newDate.NextTime()
						newDate.RoundID = round.ID
						room.ReservedDates = append(room.ReservedDates, *newDate)
						freeDatesoTheRoom = append(freeDatesoTheRoom, *newDate)
						k++
						newDate = newDate.NextDay()
					} else {
						newDate.RoundID = round.ID
						room.ReservedDates = append(room.ReservedDates, *newDate)
						freeDatesoTheRoom = append(freeDatesoTheRoom, *newDate)
						if Helper.IsReservedHalf(*newDate, room.ReservedDates) {
							oneRemain = true
						}
					}
				}
				teacher = Helper.GetFreeTeacherWithDate(teachers, freeDatesoTheRoom...)
				if &teacher == nil {
					return nil, "No Free Teacher Found For the Created Lecture "
				}
			}
			newDates := (Helper.GetFirstDate(freeDatesoTheRoom))
			if newDates != nil {
				newDate = newDates
			} else {
				return round, "Error While Getting Free Date for Lecture"
			}
			freeDatesoTheRoom = Helper.RemoveFirstDate(freeDatesoTheRoom)
			if &newDate == nil {
				return nil, " Error While Generating Training Date"
			}
			startDate, endDate := Helper.ConfigureLectureHour(*newDate, lecture.Duration)
			newLecture.StartDate = startDate
			section.ClassDates = append(section.ClassDates, startDate)
			section.Room.ReservedDates = append(section.Room.ReservedDates, startDate)
			newLecture.EndDate = endDate
			newLecture.Passed = false
			newLecture.Teacher = *teacher
			startDate.RoundID = round.ID
			teacher.BusyDates = append(teacher.BusyDates, startDate)
			newLecture.TeacherRefer = teacher.ID
			newLecture.SectionRefer = section.ID
			section.Lectures = append(section.Lectures, newLecture)
			// fmt.Println("lecture :  ", lecture.StartDate, lecture.EndDate, lecture.Course.Title)
			classes--
		}
		// file.Write([]byte("\n\n"))
		(*sections)[i] = *section
	}
	newDate := *date
	// For length of hours Gize le temariwochu of the Section Field Session Yimedeblachewal
	theDate := *date
	newDate = theDate
	for m := 0; m < int(round.TrainingDuration); m++ {
		// Setting Field Session For Each Seciiiition
		for i := 0; i < len(*sections); i++ {
			section := &((*sections)[i])
			studentQuantity := len(section.Students)
			for {
				if Helper.IsReservedHalf(newDate, section.ClassDates) || Helper.IsReservedHalf(newDate, section.TrainingDates) || Helper.IsDateIn(newDate, *breakdates) || newDate.Name == etc.Ehud {
					newDate = *(newDate.NextTimeHalf())
				} else {
					quantity := studentQuantity / 5
					if (studentQuantity % 5) > 2 {
						quantity++
					}
					unreservedtrainers := Helper.FreeTrainersOfDate(newDate, trainers, uint(quantity))
					// There has to be a trainer in that time else training time Generation eill be Failed
					if unreservedtrainers == nil {
						return nil,
							fmt.Sprintf(" There is No Trainers for the training Time To Be Held At %d/%d/%d %s %d:%d",
								newDate.Day, newDate.Month, newDate.Year, newDate.SubName, newDate.Hour, newDate.Minute)
						//  Case Where the Schedule Generation Will Have Error or Other Choice Will Be Used Here
					}
					beggining := 0
					// fmt.Println("The Length Of UnReserved Traines is ", len(*unreservedtrainers))
					for j := 0; j < len(*unreservedtrainers); j++ {
						traier := (*unreservedtrainers)[j]
						starttime, endtime := Helper.ConfigureTrainingaHour(newDate, 5)
						fieldSession := entity.FieldSession{
							StartDate:     starttime,
							EndDate:       endtime,
							Trainer:       traier,
							Passed:        false,
							RoundRefer:    round.ID,
							Round:         *round,
							FieldmanRefer: traier.ID,
						}

						if (len(section.Students)-beggining)/5 == 0 && (len(section.Students)-beggining)%5 > 0 {
							fieldSession.Students = section.Students[beggining:]
						} else {
							fieldSession.Students = section.Students[beggining : beggining+5]
						}
						beggining += 5
						starttime.RoundID = round.ID
						traier.BusyDates = append(traier.BusyDates, starttime)
						section.Trainings = append(section.Trainings, fieldSession)
						section.TrainingDates = append(section.TrainingDates, newDate)
					}
					// if trainingPerDay == 1 {
					// 	newDate = *newDate.NextDay()
					// }
					break
				}
			}
		}
		theDate = *(theDate.NextDay())
	}
	round.Trainers = *trainers
	round.Teachers = *teachers
	round.Sections = *sections
	round.Learning = true
	round.OnRegistration = false
	return round, "Succesfully Created A Round Half Mode"
}

// GenerateScheduleFullMode method
func (schedulehandler *ScheduleHandler) GenerateScheduleFullMode(RoundID uint, date *etc.Date, trainingPerDay uint) (*entity.Round, string) {
	round := schedulehandler.RoundService.GetRoundByID(RoundID)
	message := ""
	date = date.Modify()
	fmt.Println(
		"Learning :", round.Learning,
		"OnRegistration : ", round.OnRegistration,
	)
	if round.Learning {
		message = "The Round Has A Schedule"
		return round, message
	}
	if round == nil {
		return nil, "No Rond By Corresponding ID "
	}
	rooms := schedulehandler.RoomService.GetRoomsOfABranch(round.Branchnumber)
	if rooms == nil {
		message = "There Is No Room To Assign"
		return nil, message
	}
	var teachers *[]entity.Teacher
	if len(round.Teachers) > 0 {
		teachers = &round.Teachers
	} else {
		// Trainers of the Category
		teachers = schedulehandler.TeacherService.TeachersOfBranchByID(round.Branchnumber)
		if &teachers == nil {
			message = " No Teachers For Corresponding Branch "
			return nil, message
		}
	}
	var trainers *[]entity.FieldAssistant
	if len(round.Trainers) > 0 {
		trainers = &round.Trainers
	} else {
		// Trainers of the Category
		trainers = (schedulehandler.TrainerService.TrainersOfCategoryID(round.CategoryRefer))
		if trainers == nil {
			message = "Not Trainer of Corresponding Category "
			return nil, message
		}
	}
	breakdates := schedulehandler.BreakDateService.GetBreakDates(round.Branchnumber, etc.NewDate(0))
	// Selecting the Break Dates the Branch Have
	// Those Break dates are Not To Be InCluded In the Schedule Class And Training Dates
	courses := round.Courses // []entity.Course
	// fieldtime := round.TrainingDuration // uint
	lectures := round.Lectures // []entity.Lecture
	coursesDuration := 0
	for _, course := range courses {
		coursesDuration += int(course.Duration)
	}
	var passed bool
	freeRooms := []*entity.Room{}
	for d := 0; d < len(*rooms); d++ {
		rom := &((*rooms)[d])
		checkDate := *date
		checkDate.Hour = 2
		checkDate.Minute = 30
		checkDate = *(checkDate.FulFill())
		for Helper.IsDateIn(checkDate, *breakdates) || checkDate.Name == etc.Ehud {
			checkDate = *(checkDate.NextDay())
		}
		passed = true // the room is elegible if the it passes this days gize looping and Comparing
		for i := 0; i < (coursesDuration / 3); i++ {
			for _, day := range rom.ReservedDates {
				if day.Year == checkDate.Year && day.Month == checkDate.Month && day.Day == checkDate.Day && day.SubName == checkDate.SubName {
					passed = false
					break
				}
			}
			if passed == false {
				break
			} else {
				checkDate = *(checkDate.NextTime())
			}
		}
		if passed {
			freeRooms = append(freeRooms, rom)
		}
	}
	rooma := (Helper.OrderRooms(freeRooms))
	selectedRooms := Helper.SelectedRoomsFullMode(rooma, uint(len(round.Students))) // Type []*entity.Room
	totalCapacity := 0
	for _, ros := range selectedRooms {
		totalCapacity += int(ros.Capacity)
	}
	// log.Println("The Selected Rooms are  , ", len(*selectedRooms), " And The Total Capacity is : ", totalCapacity, "\n\nWhile the first one Has Length Of ", (*selectedRooms)[0].Capacity)
	sections := &[]entity.Section{} // Creating sections for the Round
	initIndex := 0
	for i := 0; i < len((selectedRooms)); i++ {
		roome := (selectedRooms)[i]
		capacity := roome.Capacity
		dirsha := int(math.Round(float64(len(round.Students)) * float64(float64(capacity)/float64(totalCapacity))))
		section := entity.Section{
			Categoryid:  round.CategoryRefer,
			Room:        roome, // The Variable Room is Pointer to A room not a room Value
			RoomRefer:   roome.ID,
			Sectionname: "Section " + strconv.Itoa(i+1),
			RoundRefer:  round.ID,
			// SectionCourseToDurations: []*entity.CourseToDuration{},
		}
		section.Students = round.Students[initIndex : initIndex+dirsha]
		*sections = append((*sections), section)
		initIndex += (dirsha)
	}

	// Generating Lectures and FieldSessions for each Section
	beginDate := *date
	beginDate.Hour = 2
	beginDate.Minute = 0
	// endDate := beginDate
	// endDate.Hour += 3
	// Generating Lectures
	coursetoduration := map[uint]uint{} //  entity.CourseToDuration{
	for i := 0; i < len(courses); i++ {
		for {
			course := courses[i]
			elapsed := coursetoduration[course.ID]
			if elapsed < course.Duration {
				remain := course.Duration - elapsed
				lectu := entity.Lecture{
					Course:      course,
					Branchid:    round.Branchnumber,
					CourseRefer: course.ID,
					Roundid:     round.ID,
				}
				if remain > entity.MaxLectureDuration {
					lectu.Duration = uint(entity.MaxLectureDuration)
					coursetoduration[course.ID] += uint(entity.MaxLectureDuration)
				} else {
					lectu.Duration = remain
					coursetoduration[course.ID] += uint(remain)
				}
				lectures = append(lectures, lectu)
				round.Lectures = append(round.Lectures, lectu)
			} else {
				break
			}
		}
	}
	// Generating Schedule For Each Section
	//In this Case it Will Be Once a Day
	for i := 0; i < len(*sections); i++ {
		newDate := *date
		newDate.Hour = 2
		newDate.Minute = 15
		classes := 0
		// duration := 0
		// hasUsedFromTheComingLecture := false
		var teacher *entity.Teacher
		section := &((*sections)[i])
		room := section.Room /// Here it is giving the Pointer to the room variable
		freeDatesoTheRoom := []etc.Date{
			// This free dates can also represent dates that are free the Room of the Section Could Have
		}
		lengthOfLectures := len(lectures)
		for k := 0; k < lengthOfLectures; k++ {
			lecture := lectures[k]
			newLecture := lecture
			if classes == 0 {
				duration := lecture.Course.Duration
				if duration == 0 {
					continue
				}
				classes = int(duration / entity.MaxLectureDuration) // Classes represents the Number of lectures a single Teachers Could HAve
				if duration%entity.MaxLectureDuration > 0 {
					classes++
				}
				// Generaate Free Dates of room to the Lecture of Specific Course to be Thought
				for k := 0; k < int(classes); k++ {
					for Helper.IsReserved(newDate, room.ReservedDates) ||
						Helper.IsDateIn(newDate, *breakdates) ||
						newDate.Name == etc.Ehud {
						if newDate.Name == etc.Ehud {
							newDate.Shift = 3
						}
						newDate = *(newDate.NextTime())
					}
					room.ReservedDates = append(room.ReservedDates, newDate)
					newDate.RoundID = round.ID
					freeDatesoTheRoom = append(freeDatesoTheRoom, newDate)
					newDate = *(newDate.NextTime())
				}
				// room.ReservedDates = append(room.ReservedDates, freeDatesoTheRoom...)
				teacher = Helper.GetFreeTeacherWithDate(teachers, freeDatesoTheRoom...)
				// pointer to the Teacher where the Teacher is in The Teachers List
				if teacher == nil {
					// Error Schedule
					return nil, " No Free Teacher Found For the Created Lectures  "
				}
			}
			newDates := (Helper.GetFirstDate(freeDatesoTheRoom))
			if newDates != nil {
				newDate = *newDates
			} else {
				return nil, "error while Selecting dates from free Dates of the Room "
			}
			freeDatesoTheRoom = Helper.RemoveFirstDate(freeDatesoTheRoom)
			if &newDate == nil {
				return nil, " Error While Generating Training Date "
			}
			startDate, endDate := Helper.ConfigureLectureHour(newDate, lecture.Duration)
			newLecture.StartDate = startDate
			section.ClassDates = append(section.ClassDates, startDate)
			newLecture.EndDate = endDate
			newLecture.Passed = false
			newLecture.Teacher = *teacher
			newLecture.SectionRefer = section.ID
			startDate.RoundID = round.ID
			teacher.BusyDates = append(teacher.BusyDates, startDate)
			// What should I Use to delete the Busy Dates Of Teacher for the Roud
			newLecture.TeacherRefer = teacher.ID
			newLecture.SectionRefer = section.ID
			section.Lectures = append(section.Lectures, newLecture)
			classes--
		}
		(*sections)[i] = *section
	}
	// Training Date assignment Algorithm
	newDate := *date
	theDate := *date
	theDate = *theDate.FulFill()
	for m := 0; m < int(round.TrainingDuration); m++ {
		// Setting Field Session For Each Section
		newDate = theDate
		for i := 0; i < len(*sections); i++ {
			section := &((*sections)[i])
			studentQuantity := len(section.Students)
			for {
				if Helper.IsReservedShift(newDate, section.ClassDates) ||
					Helper.IsReservedShift(newDate, section.TrainingDates) ||
					Helper.IsDateIn(newDate, *breakdates) ||
					newDate.Name == etc.Ehud {

					newDate = *(newDate.NextTime())

				} else {
					quantity := studentQuantity / 5
					if (studentQuantity % 5) > 2 {
						quantity++
					}
					unreservedtrainers := Helper.FreeTrainersOfDate(newDate, trainers, uint(quantity))
					if unreservedtrainers == nil && trainingPerDay == 1 {
						// If The Training in the Day is only once a day Use the Next Shift to Find a trainer
						// else if the training Per Day is 2 twice a day there has to be a trainer
						// in this shift to teach students
						newDate = *(newDate.NextShift())
						unreservedtrainers = Helper.FreeTrainersOfDate(newDate, trainers, uint(quantity))
					}
					if unreservedtrainers == nil {
						return nil, fmt.Sprintf(" There is No Trainers for the Section In This Day  %d/%d/%d", newDate.Day, newDate.Month, newDate.Year)
						//  Case Where the Schedule Generation Will Have Error or Other Choice Will Be Used Here
					}
					beggining := 0
					// fmt.Println("The Length Of UnReserved Traines is ", len(*unreservedtrainers))
					for j := 0; j < len(*unreservedtrainers); j++ {
						traier := &(*unreservedtrainers)[j]
						starttime, endtime := Helper.ConfigureTrainingaHour(newDate, 5)
						fieldSession := entity.FieldSession{
							StartDate:     starttime,
							Trainer:       *traier,
							Passed:        false,
							RoundRefer:    round.ID,
							Round:         *round,
							FieldmanRefer: traier.ID,
						}
						if (len(section.Students)-beggining)/5 == 0 && (len(section.Students)-beggining)%5 > 0 {
							fieldSession.Students = section.Students[beggining:]
						} else {
							fieldSession.Students = section.Students[beggining : beggining+5]
						}
						endtime.Hour = starttime.Hour + len(fieldSession.Students)
						fieldSession.EndDate = endtime
						beggining += 5
						section.Trainings = append(section.Trainings, fieldSession)
						starttime.RoundID = round.ID
						traier.BusyDates = append(traier.BusyDates, starttime)
						section.TrainingDates = append(section.TrainingDates, newDate)
					}
					break
				}
			}
		}
		theDate = *(theDate.NextDay())
	}
	round.Teachers = *teachers
	round.Trainers = *trainers
	round.Sections = *sections
	round.Learning = true
	round.OnRegistration = false
	return round, "Succesfully Created A Round Full Mode"
}

// ReCreateSchedule method to Create a schedule Agein
// func (schedulehandler *ScheduleHandler) ReCreateSchedule(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Add("Content-Type", "application/json")
// }

// DeleteSchedule method to delete the Schedule of a roudn from the Database
func (schedulehandler *ScheduleHandler) DeleteSchedule(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	roundidString := request.FormValue("round_id")
	RoundID, _ := strconv.Atoi(roundidString)
	fmt.Println("The Round ID ", RoundID)
	ds := struct {
		Success bool
		Message string
	}{
		Success: false,
	}
	sectionsID := schedulehandler.SectionService.GetIDOfSectionsOfRound(uint(RoundID))
	fmt.Println("Sections OF A Round : ", *sectionsID)
	for _, val := range *sectionsID {
		fmt.Println("Section ID : ", val)
	}
	successWhileDeletingUsingRoundID := schedulehandler.RoundService.DeleteScheduleDataUsingRoundID(uint(RoundID))
	if !successWhileDeletingUsingRoundID {
		// This means there was errors while deleting datas related  to Round
		ds.Message = "Error While  Deleting Datas Using the Round ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	var errorWhileDeletingUsingSection bool
	for _, id := range *sectionsID {
		errorWhileDeletingUsingSection = schedulehandler.SectionService.DeleteScheduleDataUsingSectionID(uint(id))
	}
	fmt.Println(len(*sectionsID))
	if !errorWhileDeletingUsingSection {
		// This means Error Has Happened
		ds.Message = "Error While  Deleting Datas Using the Section ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	// Now UpDating the Upper Level Datas
	// Deleting the Section
	successWhileDeletingSections := schedulehandler.SectionService.DeleteSectionsOfRound(uint(RoundID))
	if !successWhileDeletingSections {
		// There Happens Some thing while Deleting the Section of Round
		ds.Message = "Error While  Deleting Sections of Round ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	// Deleting the Section
	successWhileDeletingLectures := schedulehandler.SectionService.DeleteLecturesOfRound(uint(RoundID))
	if !successWhileDeletingLectures {
		// There Happens Some thing while Deleting the Section of Round
		ds.Message = "Error While  Deleting Lectures Using the Round ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	updateSuccess := schedulehandler.RoundService.UpdateToRegistration(uint(RoundID))
	if !updateSuccess {
		// There Happens Some thing while Deleting the Section of Round
		ds.Message = " Updating Round Not Succesfull.. "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	ds.Success = true
	ds.Message = "Succesfully Deleted"
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// AcceptLecture method to meka a schedule as Seen
// Method GET
// Variable lecture_id
// Authorization TEACHER ONLY
func (schedulehandler *ScheduleHandler) AcceptLecture(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	reqses := schedulehandler.SessionHandler.GetSession(request)
	if reqses == nil {
		han := http.NotFoundHandler()
		han.ServeHTTP(response, request)
		return
	}
	response.Header().Add("Content-Type", "application/json")
	ds := struct {
		Success   bool
		Message   string
		LectureID uint
		Lecture   entity.Lecture
	}{
		Success: false,
	}
	lectureID, era := strconv.Atoi(request.FormValue("lecture_id"))
	if era != nil {
		ds.Message = "Invalid Request "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	lecture := schedulehandler.LectureService.GetLectureByID(uint(lectureID))
	if lecture == nil {
		ds.Message = "There IS No Lecture By this ID "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if lecture.TeacherRefer != reqses.ID {
		ds.Message = "You are Not Allowed to Pass this Lecture "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if lecture.Passed {
		ds.Message = "The Lecture is alread Passed "
		ds.LectureID = uint(lectureID)
		ds.Lecture = *lecture
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	success := schedulehandler.LectureService.PassLecture(lecture.ID)
	if success {
		ds.Message = " Passing the Lecture Succesfull "
		ds.Success = true
	} else {
		ds.Message = "Error while Passing the Lecture "
		ds.Success = false
	}
	lecture.Passed = true
	ds.Lecture = *lecture
	jsonReturn, _ := json.Marshal(ds)
	response.Write(jsonReturn)
}

// DropLecture method to drop the lecture
// Method GET
// Variable lectuer_id
func (schedulehandler *ScheduleHandler) DropLecture(response http.ResponseWriter, request *http.Request) {
	reqses := schedulehandler.SessionHandler.GetSession(request)
	if reqses != nil {
		handler := http.NotFoundHandler()
		handler.ServeHTTP(response, request)
		return
	}
	// lectureID
}
