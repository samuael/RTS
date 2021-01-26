package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/Projects/RidingTrainingSystem/pkg/Helper"
	"github.com/Projects/RidingTrainingSystem/pkg/translation"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
)

// FuncMap func Map For The Templat
var FuncMap = template.FuncMap{
	"Minus":              Minus,
	"GetDate":            GetDayString,
	"GetAge":             GetAge,
	"IsAmharic":          IsAmharic,
	"GetIcon":            GetIcon,
	"GetDateString":      GetDateString,
	"DirectOpen":         DirectOpen,
	"DirectDownloadLink": DirectDownloadLink,
	"DateName":           DateName,
	"GetExtension":       Helper.GetExtension,
	"GetJsonResources":   GetJSONResources,
	"ShortTitle":         ShortTitle,
	"Tran":               translation.Translate,
	"Attach":             Attach,
}

// Attach function
func Attach(val uint, vul int) int {
	vals := fmt.Sprintf("%d%d", val, vul)
	num, ers := strconv.Atoi(vals)
	if ers != nil {
		othar, _ := strconv.Atoi(Helper.GenerateRandomString(3, Helper.NUMBERS))
		return othar
	}
	return num
}

// ShortTitle func
func ShortTitle(title string) string {
	if len(title) <= 47 {
		return title
	}
	return fmt.Sprintf(title[:45], "...")
}

// GetJSONResources functio for Converting to json
func GetJSONResources(inter interface{}) []byte {
	val, _ := json.Marshal(inter)
	return val
}

// Minus function For Subtracting float64 params
func Minus(a, b float64) float64 {
	return (a - b)
}

// GetDayString struct
func GetDayString(day uint) string {
	today := etc.NewDate(int(day))
	return fmt.Sprintf("%d/%d/%d ", today.Day, today.Month, today.Year)
}

// GetAge function
func GetAge(date *etc.Date) string {
	fmt.Println("Calles ", date)
	now := etc.NewDate(0)
	deltayear := now.Year - date.Year
	deltamonth := now.Month - date.Month
	deltaday := now.Day - date.Day
	return fmt.Sprintf("%d", ((deltayear*365 + deltamonth*30 + deltaday) / 365))
}

// IsAmharic function
func IsAmharic(value string) bool {
	variables := strings.Split(value, "")
	bytes := []byte(value)
	if 2*len(variables) == len(bytes) {
		return true
	}
	return false
}

// GetIcon function for getting the Icon for theTo Be Shown
func GetIcon(resource *entity.Resource) string {
	switch resource.Type {
	case entity.VIDEO:
		return resource.SnapShootImage
	case entity.AUDIO:
		return entity.PathToAudioIcon
	case entity.PDF:
		return entity.PathToPdfsIcon
	case entity.FILES:
		return entity.PathToFilesIcon
	case entity.IMAGES:
		return resource.Path
	default:
		return "img/taxi.jpg"
	}
}

// GetDateString functionh
func GetDateString(date etc.Date) string {
	return fmt.Sprintf("%s %d/%d/%d %d:%d", date.Name, date.Day, date.Month, date.Year, date.Hour, date.Minute)
}

// DirectOpen function
func DirectOpen(resource entity.Resource) bool {
	typ := resource.Type
	if typ == entity.AUDIO || typ == entity.VIDEO {
		return false
	}
	return true
}

// DirectDownloadLink function
func DirectDownloadLink(host string, resource entity.Resource) string {
	base := "/learning/resource/download/?rs-id="
	return fmt.Sprintf("%s%d", base, resource.ID)
}

// DateName function
func DateName(date etc.Date) string {
	date = *date.Modify()
	date = *date.FulFill()
	return date.Name
}

// OtherSceneOpenable function
// func OtherSceneOpenable(  resource entity.Resource ) bool {
// 	typ := resource.Type
// 	if typ == entity.PDF || typ == entity.IMAGES ||
// }
