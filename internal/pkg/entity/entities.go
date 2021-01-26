package entity

import (
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/Projects/RidingTrainingSystem/pkg/form"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	// "github.com/samuael/RidingTrainingSystem/entity"
)

var (
	// PathToStudentsFromTemplates constant
	PathToStudentsFromTemplates = "img/Students/"
	// TemplateDirectoryFromMain  varaibe
	TemplateDirectoryFromMain = "../../web/templates/"
	// PathToPdfs constant
	PathToPdfs = "../../web/templates/pdfs/"
	// MediaDirectoryFromTemplates const varaible
	MediaDirectoryFromTemplates = "Source/Resources/media/"
	// VIDEOS file format
	VIDEOS = []string{"mp4", "m4a", "m4v", "f4v", "f4a", "m4b",
		"m4r", "f4b", "mov", "3gp", "3gp2", "3g2", "3gpp", "3gpp2",
		"wmv", "wma", "asf*", "mkv", "webm", "flv", "avi*", "dvr"}
	// AUDIOS file format
	AUDIOS = []string{
		"mp3", "wav", "aac", "wma", "flac", "amr", "ogg", "mpeg", "midi",
	}
	// PICTURES file format
	PICTURES = []string{"jpeg", "png", "jpg", "gif", "btmp"}
	// PDFS file format
	PDFS = []string{"pdf", "ebook"}
	// DynamicRoundPopulatingHeadersEnglish variables
	DynamicRoundPopulatingHeadersEnglish = []string{
		"Full name", "Sex", "Age", "Acadamic Status",
		"Language", "City", "Kebele", "Phone"}
)

/*

$mimeTypes = array(
	'pdf' => 'application/pdf',
	'txt' => 'text/plain',
	'html' => 'text/html',
	'exe' => 'application/octet-stream',
	'zip' => 'application/zip',
	'doc' => 'application/msword',
	'xls' => 'application/vnd.ms-excel',
	'ppt' => 'application/vnd.ms-powerpoint',
	'gif' => 'image/gif',
	'png' => 'image/png',
	'jpeg' => 'image/jpg',
	'jpg' => 'image/jpg',
	'php' => 'text/plain'
);

*/

const (
	// PathToQuestionImages  constant Path to Images Directory OF Questions
	PathToQuestionImages = "../../web/templates/img/QuestionImages/"
	// PathToQuestionImagesFromTemplates constant
	PathToQuestionImagesFromTemplates = "img/QuestionImages/"
	// RoundTemplatingCSV file name of the Default file Populating file
	RoundTemplatingCSV = "populating_template.csv"
	// PathToStaticTemplateFiles  constant path to PopulatingResource
	PathToStaticTemplateFiles = "../../web/templates/Source/PopulatingResource/"
	// PathToBranchImagesFromTemplates  constant
	PathToBranchImagesFromTemplates = "img/BranchImages/"
	// All constant Signature
	All = 0
	// Active constant
	Active = 1
	// Passive constant
	Passive = 2
	// ShiftOnceADay constant
	ShiftOnceADay = 1
	// ShiftWheneFreeADay const
	ShiftWheneFreeADay = 2
	// PathToAudioIcon  conant path to system Audio Image
	PathToAudioIcon = "img/SystemImages/music_icon.png"
	// PathToPdfsIcon path to pdf Icon
	PathToPdfsIcon = "img/SystemImages/pdf_image.png"
	// PathToFilesIcon path to files image icon
	PathToFilesIcon = "img/SystemImages/file_image.jpg"
	// PathToFirstFrames constant
	PathToFirstFrames = "../../web/templates/Source/Resources/FirstFrames/"
	// PathToTemplates const
	PathToTemplates = "../../web/templates/"
	// PATHToZipFiles  constant
	PATHToZipFiles = "../../web/templates/Source/ZipFiles/"
	// PATHToVehiclesImageFromMain constant
	PATHToVehiclesImageFromMain = "../../web/templates/Source/Vehicles/"
	// TemporaryFilesDirectoryFromMain constant
	TemporaryFilesDirectoryFromMain = "../../web/templates/Source/Resources/TemporaryFiles"
	// PathToResources  constant
	PathToResources = "../../web/templates/Source/Resources/"
	// FileSchema  const
	FileSchema = "file:///"
	// MediaFromMain constant
	MediaFromMain = "../../web/templates/Source/Resources/media"
	// RESOURCE_PATH constant
	RESOURCE_PATH = "Source/Resources/"
	// VIDEO  resource
	VIDEO = 1
	// AUDIO resource
	AUDIO = 2
	// PDF resource
	PDF = 3
	// IMAGES resource
	IMAGES = 4
	// FILES resource
	FILES = 5
	// PROTOCOL string
	PROTOCOL = "http://"
	// HOST string
	HOST = "localhost:9900/"
	// HALF value
	HALF = "half"
	// FULL value
	FULL = "full"
	// OWNER role only One Person Should Have
	OWNER = "OWNER"
	// SUPERADMIN role of admin type
	SUPERADMIN = "SUPERADMIN"
	// SECRETART role of admin type
	SECRETART = "SECRETARY"
	// TEACHER role of type Teacher
	TEACHER = "TEACHER"
	// STUDENT role of type teacher
	STUDENT = "STUDENT"
	// FIELDMAN role of type fieldman
	FIELDMAN = "FIELDMAN"
	// MaxLectureDuration constant
	MaxLectureDuration = 2
	// MaxBreakTimeMinute constant  number of minute a break time the studdents could have
	MaxBreakTimeMinute = 15
)

// User to be used by all the User of the system and the to Be Included in the Session Header
type User struct {
	Username string
	Password string
	Imageurl string
}

// Users type Interface
type Users interface {
	GetUsername() string
	GetID() uint
	GetImageURL() string
	GetLang() string
	GetBranchID() uint
	GetRole() string
	SetImageURL(val string)
}

// CourseToDuration this structure is created to Tell whether the Student is taking the Appropriate amount of Course or not
// with out missing course
// the relation ship is manu to many with the Section Table and inside the Structure
// the Course struct will have a belongs to relation ship
// meaning the structure  StudentCourseToDuration Struct will have a belongs to relationship with Course
type CourseToDuration struct {
	gorm.Model
	Course       uint `gorm:"foreignkey:SectionRefer"`
	Elapsed      uint
	SectionRefer uint // Section Refer represents
}

// Payment representing payment transaction held by Secretary
type Payment struct {
	gorm.Model
	Student Student `gorm:"foreignkey:StudentRefer"`
	Admin   Admin   `gorm:"foreignkey:AdminRefer"`
	// Branch       Branch  `gorm:"foreignkey:BranchRefer"`
	Round        *Round `gorm:"foreignkey:RoundRefer"`
	Amount       float64
	Date         *etc.Date `gorm:"foreignkey:DateRefer"`
	BranchRefer  uint
	StudentRefer uint
	RoundRefer   uint
	AdminRefer   uint
	DateRefer    uint
}

// Address  representing address
type Address struct {
	gorm.Model
	Region  string
	Zone    string
	Woreda  string
	Kebele  string
	City    string
	OwnerID uint
}

// StudentCsv  struct
type StudentCsv struct {
	Fullname       string `csv:"fullname"`
	Sex            string `csv:"sex"`
	City           string `csv:"city"`
	Kebele         string `csv:"kebele"`
	Phone          string `csv:"phone"`
	Lang           string `csv:"lang"`
	AcademicStatus string `csv:"academic_status"`
	Age            int    `csv:"age"`
}

// Student representing the Student or Trainee
type Student struct {
	gorm.Model
	ImageNumber           string
	Username              string `gormSigningKey:"type:varchar(255);not null; unique"` // categoryname/id/round
	Firstname             string
	Lastname              string
	GrandFatherName       string
	MotherName            string
	Nickname              string
	Sex                   string
	BirthDate             *etc.Date `gorm:"foreignkey:BirthDateRefer"`
	BirthDateRefer        uint
	Address               *Address `gorm:"foreignkey:AddressRefer"`
	BirthAddress          *Address `gorm:"foreignkey:BirthAddressRefer"`
	BirthAddressRefer     uint
	AddressRefer          uint
	MaritialStatus        string `gorm:"default:'Not Married'"`
	PhoneNumber           string
	FamilyCount           uint
	PartnerFullname       string
	PartnerPhoneNumber    string
	GuarantorFullName     string
	GuarantorPhoneNumber  string
	GuarantorAddress      *Address `gorm:"foreignkey:GuarantorAddressRefer"`
	GuarantorAddressRefer uint
	PreviousLicenceType   string
	PreviousLicenceNumber string
	CategoryID            uint
	Category              Category `gorm:"foreignkey:CategoryID"`
	AcademicStatus        string
	Password              string
	Round                 *Round `gorm:"foreignkey:RoundRefer"`
	Lang                  string
	RedisteredBy          string
	Imageurl              string
	// Active      	          bool    `gorm:"default:true"`
	Section               *Section `gorm:"foreignkey:SectionRefer"`
	AskedQuestionsCount   uint     `json:"askedcount,omitempty"`
	AnsweredQuestionCount uint     `json:"answerdcount,omitempty"`
	Active                bool     `gorm:"default:true"` // This tells whether the student is active or not
	SectionRefer          uint     /// meaning this struct belongs to this section niggoye
	RoundRefer            uint     // meaning this Student Struct belongs to this Round Instance
	PaidAmount            float64  `gorm:"default:0"`
	BranchID              uint
}

// GetUsername for implementing the Users interface
func (user *Student) GetUsername() string {
	return user.Username
}

// GetID returning the Id of the Student
func (user *Student) GetID() uint {
	return user.ID
}

// GetImageURL returns string
func (user *Student) GetImageURL() string {
	return user.Imageurl
}

// GetLang returns String
func (user *Student) GetLang() string {
	return user.Lang
}

// GetBranchID returns uint
func (user *Student) GetBranchID() uint {
	return user.BranchID
}

// GetRole returning the role of the Student
func (user *Student) GetRole() string {
	return "Student"
}

// SetImageURL returning strign
func (user *Student) SetImageURL(val string) {
	user.Imageurl = val
}

// Session representing the Sesstion to Be sent with the request body
// no saving of a session in the database so i Will use this session in place of
type Session struct {
	jwt.StandardClaims
	Username string
	ID       uint
	Imageurl string
	Lang     string
	Role     string
	BranchID uint // Imageurl and Username are Included in this iggo
}

// Teacher representing the Lecture for the Training Center
type Teacher struct {
	gorm.Model
	Username     string // categoryname/id/round
	Email        string
	Password     string
	Firstname    string
	Lastname     string
	GrandName    string
	BranchNumber uint
	Createdby    string
	Imageurl     string
	Lang         string
	Phonenumber  string
	Active       bool       `gorm:"default:true"` // represents whether the teacher is teaching or not
	BusyDates    []etc.Date `gorm:"many2many:lectures_bussy_date;"`
}

// GetUsername for implementing the Users interface
func (user *Teacher) GetUsername() string {
	return user.Username
}

// GetID returning the Id of the Teacher
func (user *Teacher) GetID() uint {
	return user.ID
}

// GetImageURL returns string
func (user *Teacher) GetImageURL() string {
	return user.Imageurl
}

// GetLang returns String
func (user *Teacher) GetLang() string {
	return user.Lang
}

// GetBranchID returns uint
func (user *Teacher) GetBranchID() uint {
	return user.BranchNumber
}

// GetRole returning the role of the Teacher
func (user *Teacher) GetRole() string {
	return "Teacher"
}

// SetImageURL returning strign
func (user *Teacher) SetImageURL(val string) {
	user.Imageurl = val
}

// FieldAssistant Representing the filed Trainer
type FieldAssistant struct {
	gorm.Model
	Username     string // categoryname/id/round
	Firstname    string
	Lastname     string
	GrandName    string
	Email        string
	Imageurl     string
	Phonenumber  string
	Lang         string
	BranchNumber uint
	Createdby    string
	Password     string
	CategoryID   uint
	VehicleID    uint
	Reserved     bool       `gorm:"default:false"`
	Vehicle      Vehicle    `gorm:"foreignkey:VehicleID"`
	Categoty     Category   `gorm:"foreignkey:CategoryID"`
	BusyDates    []etc.Date `gorm:"many2many:field_dates"`
	Active       bool       `gorm:"default:true"` // Represents whether the admin is working in the Company or not
}

// FieldDate string
type FieldDate struct {
	FieldAssistantID uint
	DateID           uint
}

// GetUsername for implementing the Users interface
func (user *FieldAssistant) GetUsername() string {
	return user.Username
}

// GetID returning the Id of the Teacher
func (user *FieldAssistant) GetID() uint {
	return user.ID
}

// GetImageURL returns string
func (user *FieldAssistant) GetImageURL() string {
	return user.Imageurl
}

// GetLang returns String
func (user *FieldAssistant) GetLang() string {
	return user.Lang
}

// GetBranchID () uint returns uint
func (user *FieldAssistant) GetBranchID() uint {
	return user.BranchNumber
}

// GetRole returning the role of the Teacher
func (user *FieldAssistant) GetRole() string {
	return "FieldAssistant"
}

// SetImageURL returning strign
func (user *FieldAssistant) SetImageURL(val string) {
	user.Imageurl = val
}

// Vehicle struct representing the
type Vehicle struct {
	gorm.Model
	Imageurl    string
	BoardNumber string `gorm:"unique;"`
	TrainerID   uint
	// Trainer     FieldAssistant `gorm:"foreignkey:TrainerID"`
	BranchNo   uint
	CategoryID uint
	Branch     Branch   `gorm:"foreignkey:BranchRefer"`
	Category   Category `gorm:"foreginkey:CategoryRefer"`
	Reserved   bool     `gorm:"default:false"`
}

// Admin representing the Controllers of each Branch
type Admin struct {
	gorm.Model
	Uploader    // Username Imgurl
	Firstname   string
	Lastname    string
	GrandName   string
	Password    string
	Email       string
	Lang        string
	Phone       string
	Createdby   string
	Branch      Branch `gorm:"foreignkey:BranchRefer"`
	BranchRefer uint   // the Branch and the Admin will have a Belongs relationship after now on
	Role        string
	Active      bool `gorm:"default:true"` // represents whether the Admin is working in the system or not
}

// GetUsername for implementing the Users interface
func (user *Admin) GetUsername() string {
	return user.Username
}

// GetID returning the Id of the Admin
func (user *Admin) GetID() uint {
	return user.ID
}

// GetImageURL returns string
func (user *Admin) GetImageURL() string {
	return user.Imageurl
}

// GetLang returns String
func (user *Admin) GetLang() string {
	return user.Lang
}

// GetBranchID  returns uint
func (user *Admin) GetBranchID() uint {
	return user.BranchRefer
}

// GetRole returning the role of the Admin
func (user *Admin) GetRole() string {
	return user.Role
}

// SetImageURL returning strign
func (user *Admin) SetImageURL(val string) {
	user.Imageurl = val
}

// Course representing the Lecture given Courses
type Course struct {
	gorm.Model
	Title      string `json:"title"`
	Duration   uint   `json:"duration"`
	Categoryid uint   `json:"categoryid"`
	OwnerID    uint
	BranchID   uint `json:"branchid"`
}

// Category representing the types (kinds ) of training given By the Training Center
type Category struct {
	gorm.Model
	Title string `json:"title"`
	// TitleAmharic string `json:"titleAmharic"`
	ImageURL string `json:"image,omitempty"`
	// Description string
	Rounds                   []Round `gorm:"many2many:category_round"`
	Branchid                 uint    `json:"branchid"`
	LastRoundNumber          uint
	TrainedNumberOfStudents  uint
	NumberOfStudentsLearning uint
	Active                   bool `gorm:"default:true"`
}

// Branch represening Speific Branch of the Training Center
type Branch struct {
	gorm.Model
	Name                  string
	Country               string
	Address               Address `gorm:"foreignkey:AddressRefer"`
	AddressRefer          uint
	LicenceGivenDate      etc.Date `gorm:"foreignkey:LicenceGivenDateRefer"`
	LicenceGivenDateRefer uint
	Email                 string
	Moto                  string
	Lang                  string
	BranchAcronym         string // example SHA
	BranchFullnameAmharic string
	BranchFullnameEnglish string
	LicenceNumber         string
	City                  string
	Createdby             string
	Logourl               string         // Imageurl for specific Branch this can be created using the
	Phones                pq.StringArray `gorm:"type:varchar(20)[];"  json:"phones"`
}

// Resource representing teaching resourses
// such as Videos  , Audios  , and Books
type Resource struct {
	gorm.Model
	Path       string
	UploadDate etc.Date `gorm:"foreignkey:UDID"`
	UDID       uint
	// NumberOfView   uint
	Title          string
	Description    string
	Uploadedby     string
	UploaderImage  string
	SnapShootImage string
	HLSDirectory   string
	FirstFrame     string
	Type           uint
}

// Uploader  representing the Uploader
type Uploader struct {
	// Just For Inheritance
	Username string
	Imageurl string `gorm:"type:varchar(255);not null"`
}

// Room representing the Lecture Room for Lecture
// Having intake capacity
type Room struct {
	gorm.Model
	Number        uint `json:"number"`
	Branchid      uint
	CreatedBy     string
	Capacity      uint       `json:"capacity"`
	ReservedDates []etc.Date `gorm:"many2many:room_to_date"`
}

// RoomToDate struct
type RoomToDate struct {
	gorm.Model
	RoomID uint
	DateID uint
}

// Section  representing a class or Section for Specific Catagory in specific Branch
type Section struct {
	gorm.Model
	Sectionname   string
	Room          *Room `gorm:"foreignkey:RoomRefer"`
	RoomRefer     uint
	Categoryid    uint
	Round         Round `gorm:"foreignkey:RoundRefer"`
	RoundRefer    uint
	Trainings     []FieldSession `gorm:"many2many:section_trainings"`
	TrainingDates []etc.Date     `gorm:"many2many:section_training_dates"`
	Lectures      []Lecture      `gorm:"many2many:section_lectures"`
	Students      []Student      `gorm:"many2many:section_students"`
	// SectionCourseToDurations []*CourseToDuration `gorm:"many2many:section_course_to_duration"`
	ClassDates []etc.Date `gorm:"many2many:section_class_dates"`
	OwnerID    uint
}

// SectionTraining method
type SectionTraining struct {
	FieldSessionID uint
	SectionID      uint
}

// SectionTrainingDate struct
type SectionTrainingDate struct {
	DateID    uint
	SectionID uint
}

// SectionLecture struct
type SectionLecture struct {
	LectureID uint
	SectionID uint
}

// SectionStudent struct
type SectionStudent struct {
	StudentID uint
	SectionID uint
}

// SectionClassDate struct
type SectionClassDate struct {
	DateID    uint
	SectionID uint
}

// Lecture representing Theory Teaching Schedulecou
type Lecture struct {
	gorm.Model
	Branchid       uint
	Teacher        Teacher `gorm:"foreignkey:TeacherRefer"`
	TeacherRefer   uint
	CourseRefer    uint
	SectionRefer   uint
	Duration       uint
	Course         Course   `gorm:"foreignkey:CourseRefer"`
	StartDate      etc.Date `gorm:"foreignkey:StartDateRefer"`
	StartDateRefer uint
	EndDate        etc.Date `gorm:"foreignkey:EndDateRefer"`
	EndDateRefer   uint
	Passed         bool
	Roundid        uint
	Round          Round `gorm:"foreignkey:Roundid"`
}

// Question representing the Questions for the Theoretical Exercise for the Students
type Question struct {
	gorm.Model
	Body           string         `gorm:"type:TEXT;not null" json:"body"`
	Answers        pq.StringArray `gorm:"type:varchar(100)[];not null"  json:"answers"`
	Answerindex    uint           `gorm:"type:;not null" json:"answer_index"`
	ImageDirectory string
}

// AskedQuetion representing the Questions the User got Asked
type AskedQuetion struct {
	gorm.Model
	Studentid   uint
	Questionsid pq.Int64Array `gorm:"type:Int[] "`
}

// FieldSession representing the Field shedule
type FieldSession struct {
	gorm.Model
	Trainer       FieldAssistant `gorm:"foreignkey:FieldmanRefer"`
	StartDate     etc.Date       `gorm:"foreignkey:StartDateID"`
	StartDateID   uint
	EndDate       etc.Date `gorm:"foreignkey:EndDateID"`
	EndDateID     uint
	Passed        bool
	Round         Round `gorm:"foreignkey:RoundRefer"`
	RoundRefer    uint
	FieldmanRefer uint
	Students      []Student `gorm:"many2many:field_student"`
}

// Round representing the Round with ti's status
type Round struct {
	gorm.Model
	Year               uint      `json:"year"`
	Roundnumber        uint      `json:"roundnumber"`
	Courses            []Course  `gorm:"many2many:round_course"`
	Sections           []Section `gorm:"many2many:round_section"`
	Branchnumber       uint      `json:"branchnumber"`
	Category           Category  `gorm:"foreignkey:CategoryRefer"`
	CategoryRefer      uint
	Active             bool  `json:"active,omitempty"`
	CreatedBY          Admin `gorm:"foreignkey:AdminRefer"`
	AdminRefer         uint
	Cost               float64 `json:"cost"`
	TrainingSchiftCode uint
	TrainingDuration   uint
	TotalPaid          float64
	Trainers           []FieldAssistant `gorm:"many2many:round_trainers"`
	Teachers           []Teacher        `gorm:"many2many:round_teachers"`
	Lectures           []Lecture        `gorm:"many2many:round_lectures"`
	OnRegistration     bool             `json:"onregistration" gorm:"default:true"`
	Learning           bool             `json:"learning" gorm:"default:false"`
	Students           []Student        `gorm:"many2many:round_students"`
	Studentscount      uint             `json:"Studentscount"`
	MaxStudents        uint             `json:"max_students"`
	SchedulePath       string
}

// ActiveRound representing the Rounds Running the Company is Teaching
type ActiveRound struct {
	ID       uint
	Branchid uint
	Roundid  uint
}

// SchedulePageStructure  representing the datastructure for Schcedule Page
type SchedulePageStructure struct {
	Round         Round
	Sections      []Section
	Lectures      []Lecture
	FieldSessions []FieldSession
}

// Information  representing Admin aploaded informations
type Information struct {
	gorm.Model
	BranchID    uint
	Username    string
	Title       string
	Description string
	BranchName  string
	Active      bool `gorm:"default:true;"`
}

// SystemService representing SystemWide Message
type SystemService struct {
	Success bool
	Message string
}

// Informing representing Template embedding informations
type Informing struct {
	CSRF     string
	Message  string
	HasError bool
	Success  bool
	Branch   Branch
	Admin    Admin
	Host     string
	Branches []Branch
	Input    form.Input
	Specific interface{}
}

// MainPageInfo struct
type MainPageInfo struct {
	Rooms      []*Room
	Branchs    *[]Branch
	Branch     Branch
	Categories []*Category
}
