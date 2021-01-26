package Helper

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// NUMBERS const numbers
const NUMBERS = "1234567890"

// CHARACTERS const field
const CHARACTERS = "abcdefghijelmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_1234567890"

// GenerateRandomString  function
func GenerateRandomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// OrderRooms function
func OrderRooms(rooms []*entity.Room) []*entity.Room {
	sortedArray := []*entity.Room{}
	for i := len(rooms) - 1; i >= 0; i-- {
		for k := 0; k <= i-1; k++ {
			if rooms[k].Capacity > rooms[k+1].Capacity {
				temp := rooms[k]
				rooms[k] = rooms[k+1]
				rooms[k+1] = temp
			}
		}
	}
	for j := len(rooms) - 1; j >= 0; j-- {
		sortedArray = append(sortedArray, rooms[j])
	}
	return sortedArray
}

// ConfigureLectureHour method
func ConfigureLectureHour(time etc.Date, duration uint) (startdate, enddate etc.Date) {
	time = *time.Modify()
	startdate = time
	enddate = time
	enddate.Hour += int(duration)
	return startdate, enddate
}

// ConfigureTrainingaHour method
func ConfigureTrainingaHour(time etc.Date, NumberOfStudents uint) (starttime, endtime etc.Date) {
	time = *time.Modify()
	starttime = time
	endtime = time
	if time.SubName == etc.TEWAT {
		starttime.Hour = 1
		starttime.Minute = 30
		endtime.Hour = starttime.Hour + int(NumberOfStudents)
		endtime.Minute = starttime.Minute
	} else if time.SubName == etc.KESEAT {
		starttime.Hour = 6
		starttime.Minute = 30
		endtime.Hour = starttime.Hour + int(NumberOfStudents)
		endtime.Minute = starttime.Minute
	}
	return starttime, endtime
}

// SelectedRoomsHalfMode function
func SelectedRoomsHalfMode(rooms []*entity.Room, Capacity uint) []*entity.Room {
	newRooms := []*entity.Room{}
	roomsa := rooms
	for i := len(rooms) - 1; i >= 0; i-- {
		if roomsa[i].Capacity >= Capacity {
			newRooms = append(newRooms, roomsa[i])
			return newRooms
		}
	}
	for i := len(rooms) - 1; i >= 0; i-- {
		for j := i; j >= 0; j-- {
			if (roomsa[i].Capacity + roomsa[j].Capacity) >= Capacity {
				newRooms = append(newRooms, roomsa[i], roomsa[j])
				return newRooms
			}
		}
	}
	var tempo *entity.Room
	tempo = nil
	for i := 0; i < len(rooms); i++ {
		if roomsa[i].Capacity < Capacity {
			if tempo == nil {
				newRooms = append(newRooms, roomsa[i])
				continue
			}
			newRooms = append(newRooms, tempo)
			return newRooms
		} else if roomsa[i].Capacity == Capacity {
			newRooms = append(newRooms, roomsa[i])
			return newRooms
		} else {
			tempo = roomsa[i]
			continue
		}
	}
	return newRooms
}

// SelectedRoomsFullMode function
func SelectedRoomsFullMode(rooms []*entity.Room, Capacity uint) []*entity.Room {
	newRooms := []*entity.Room{}
	roomsa := rooms
	for i := len(rooms) - 1; i >= 0; i-- {
		if roomsa[i].Capacity >= Capacity {
			newRooms = append(newRooms, roomsa[i])
			return newRooms
		}
	}
	for i := len(rooms) - 1; i >= 0; i-- {
		for j := i - 1; j >= 0; j-- {
			if (roomsa[i].Capacity + roomsa[j].Capacity) >= Capacity {
				newRooms = append(newRooms, roomsa[i], roomsa[j])
				return newRooms
			}
		}
	}
	var tempo *entity.Room
	tempo = nil
	newRooms = append(newRooms, roomsa[0])
	Capacity = Capacity - roomsa[0].Capacity
	for i := 1; i < len(rooms); i++ {
		if roomsa[i].Capacity < Capacity {
			if tempo == nil {
				newRooms = append(newRooms, roomsa[i])
				continue
			}
			newRooms = append(newRooms, tempo)
			return newRooms
		} else if roomsa[i].Capacity == Capacity {
			newRooms = append(newRooms, roomsa[i])
			return newRooms
		} else {
			tempo = roomsa[i]
			continue
		}
	}
	return newRooms
}

// IsReserved   function
func IsReserved(date etc.Date, dates []etc.Date) bool {
	date = *date.Modify()
	date = *date.FulFill()
	for _, dato := range dates {
		if date.Year == dato.Year && dato.Month == date.Month && date.Day == dato.Day && date.SubName == dato.SubName && dato.Shift == date.Shift {
			return true
		}
	}
	return false
}

// IsReservedHalf function to check whether the date is reserved at that Half Shift Or Not
func IsReservedHalf(date etc.Date, dates []etc.Date) bool {
	date = *date.Modify()
	date = *date.FulFill()
	for _, dato := range dates {
		if date.Year == dato.Year && dato.Month == date.Month && date.Day == dato.Day && date.SubName == dato.SubName {
			return true
		}
	}
	return false
}

// IsReservedShift function
func IsReservedShift(date etc.Date, dates []etc.Date) bool {
	for _, dato := range dates {
		if date.Year == dato.Year && dato.Month == date.Month && date.Day == dato.Day && date.SubName == dato.SubName {
			return true
		}
	}
	return false
}

// IsReservedFull function
func IsReservedFull(date etc.Date, dates []etc.Date) bool {
	for _, dato := range dates {
		if date.Year == dato.Year && dato.Month == date.Month && date.Day == dato.Day {
			return true
		}
	}
	return false
}

// SelectCourseByID function
func SelectCourseByID(ID uint, courses []entity.Course) *entity.Course {
	for _, course := range courses {
		if course.ID == ID {
			return &course
		}
	}
	return nil
}

// NextTime function
func NextTime(date etc.Date) etc.Date {
	if date.SubName == etc.TEWAT {
		date.Hour = 7
		date.Minute = 20
		date.SubName = etc.KESEAT
	} else {
		date.Day++
		date.SubName = etc.TEWAT
	}
	return date
}

// GetFreeTeacherWithDate function
func GetFreeTeacherWithDate(teachers *[]entity.Teacher, dates ...etc.Date) *entity.Teacher {
	lena := len(*teachers)
	// Shuffling the Teachers Table
	for n := 0; n < lena; n++ {
		randInta := rand.Int63n(int64(lena))
		randIntb := rand.Int63n(int64(lena))
		val := (*teachers)[randInta]
		(*teachers)[randInta] = (*teachers)[randIntb]
		(*teachers)[randIntb] = val
	}
	// Selecting from random students
	for i := 0; i < len(*teachers); i++ {
		teacher := &(*teachers)[i]
		isReserved := false
		for _, newDate := range teacher.BusyDates {
			for _, inputDate := range dates {
				if newDate.Year == inputDate.Year &&
					newDate.Month == inputDate.Month &&
					inputDate.Day == newDate.Day &&
					inputDate.SubName == newDate.SubName &&
					inputDate.Shift == newDate.Shift {
					isReserved = true
					break
				}
			}
		}
		if !isReserved {
			return teacher
		}
	}
	return nil
}

// GetFirstDate function
func GetFirstDate(dates []etc.Date) *etc.Date {
	if len(dates) > 0 {
		return &dates[0]
	}
	return nil
}

// RemoveFirstDate method representing
func RemoveFirstDate(dates []etc.Date) []etc.Date {
	newDates := []etc.Date{}
	for i := 1; i < len(dates); i++ {
		newDates = append(newDates, dates[i])
	}
	return newDates
}

// GetFreeTeacherWithFullDate function
func GetFreeTeacherWithFullDate(teachers []entity.Teacher, dates ...etc.Date) *entity.Teacher {
	for _, teacher := range teachers {
		isReserved := false
		for _, newDate := range teacher.BusyDates {
			for _, inputDate := range dates {
				if newDate.Year == inputDate.Year && newDate.Month == inputDate.Month && inputDate.Day == newDate.Day {
					isReserved = true
					break
				}
			}
		}
		if !isReserved {
			return &teacher
		}
	}
	return nil
}

// FreeTrainersOfDate function
func FreeTrainersOfDate(date etc.Date, trainers *[]entity.FieldAssistant, quantity uint) *[]entity.FieldAssistant {
	freeTrainers := &[]entity.FieldAssistant{}
	count := 0
	for l := 0; l < len(*trainers); l++ {
		isReserved := false
		trainer := &((*trainers)[l])
		for _, newDate := range trainer.BusyDates {
			if newDate.Year == date.Year && newDate.Month == date.Month && date.Day == newDate.Day && date.SubName == newDate.SubName {
				isReserved = true
				break
			}
		}
		if !isReserved {
			*freeTrainers = append(*freeTrainers, *trainer)
			count++
		}
		if uint(count) == quantity {
			return freeTrainers
		}
	}
	return nil
}

// IsDateIn function to check whether the date is in the list of dates
func IsDateIn(date etc.Date, dates []etc.Date) bool {
	for _, dato := range dates {
		if date.Month == dato.Month && dato.Year == date.Year && date.Day == dato.Day {
			return true
		}
	}
	return false
}

// MarshalThis function
func MarshalThis(inter interface{}) []byte {
	val, era := json.Marshal(inter)
	if era != nil {
		return nil
	}
	return val
}
