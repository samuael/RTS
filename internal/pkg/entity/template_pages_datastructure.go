package entity

// LearningPageDS struct
type LearningPageDS struct {
	Success            bool
	Message            string
	Branch             Branch
	User               interface{}
	Resources          []Resource
	RouteMap           map[string]string
	VRouteMap          map[string]string
	ActiveResource     Resource
	HOST               string
	VerticalNavIndex   uint
	HorizontalNavIndex uint
	Lang               string
	ResCount 	   int 
}
