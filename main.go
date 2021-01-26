package main

import (
	"fmt"

	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

func main() {
	date := etc.NewDate(0)
	// nextTC := 0
	date = date.Modify()
	for k := 0; k < 9; k++ {
		// date = date.NextTime()
		fmt.Printf("%s >> %d/%d/%d  %s  %d:%d \n Shift : %d \n", date.Name, date.Day, date.Month, date.Year, date.SubName, date.Hour, date.Minute, date.Shift)
		// if nextTC < 2 {
		// 	nextTC++
		// 	date = date.NextTime()
		// } else {
		date = date.NextTimeHalf()
		// date = date.BackShiftTime()
		// date = date.BackShiftTime()
		// 	nextTC = 0
		// }
	}
	*date = date.BackShiftTime()
	fmt.Printf("Back Shifted :\n%s >> %d/%d/%d  %s  %d:%d \n Shift : %d \n", date.Name, date.Day, date.Month, date.Year, date.SubName, date.Hour, date.Minute, date.Shift)

}

// fmt.Printf("%s >> %d/%d/%d  %s  %d:%d \n Shift : %d \n", date.Name, date.Day, date.Month, date.Year, date.SubName, date.Hour, date.Minute, date.Shift)
//
