package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Projects/RidingTrainingSystem/cmd/DTIS/MediaRouter"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Section"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Section/SectionRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Section/SectionService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/BreakDates"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Resource"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Resource/ResourceRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Resource/ResourceService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Question"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Question/QuestionRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Question/QuestionService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Information"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Information/InformationRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Information/InformationService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Schedule"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Schedule/ScheduleRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Schedule/ScheduleService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Payment"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Course"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Course/CourseRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Course/CourseService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/BreakDates/BreakDateRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/BreakDates/BreakDateService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Room"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Room/RoomRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Room/RoomService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round/RoundRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round/RoundService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Payment/PaymentRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Payment/PaymentService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Vehicle"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Vehicle/VehicleRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Vehicle/VehicleService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Category"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Student"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer/TrainerRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Trainer/TrainerService"

	apihandler "github.com/Projects/RidingTrainingSystem/api/ApiHandler"
	DB "github.com/Projects/RidingTrainingSystem/internal/pkg/Db"
	stdrepo "github.com/Projects/RidingTrainingSystem/internal/pkg/Student/repository"

	adminre "github.com/Projects/RidingTrainingSystem/internal/pkg/Admin/AdminRepository"
	adminse "github.com/Projects/RidingTrainingSystem/internal/pkg/Admin/AdminService"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch/BranchRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch/BranchService"
	
	stdserv "github.com/Projects/RidingTrainingSystem/internal/pkg/Student/service"
	teachRepo "github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher/TeacherRepository"
	teachSer "github.com/Projects/RidingTrainingSystem/internal/pkg/Teacher/TeacherService"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Category/CategoryRepo"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Category/CategoryService"

	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/web/handler"
	"github.com/jinzhu/gorm"
)

var once = sync.Once{}
var th *handler.TemplateHandler
var systemTemplate *template.Template
var db *gorm.DB
var erro error
var sessionHandler *session.Cookiehandler
var studentHandler *handler.StudentHandler

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			return
		}
		next.ServeHTTP(w, r)
	})
}

var resourceRepo Resource.ResourceRepo
var resourceService Resource.ResourceService

var courseRepo Course.CourseRepo
var courseService Course.CourseService

var roundRepo Round.RoundRepo
var roudnService Round.RoundService

var vehicleRepo Vehicle.VehicleRepo
var vehicleService Vehicle.VehicleService

var branchService Branch.BranchService
var branchrepo Branch.BranchRepo

var userhandler *handler.UserHandler
var adminhandler *handler.AdminHandler

var studentservice Student.StudentService
var studentrepo Student.StudentRepository

var teacherRepo teachRepo.TeacherRepo
var teacherService teachSer.TeacherService

var trinServ Trainer.TrainerService
var trainRepo Trainer.TrainerRepo

var categoryRepo Category.CategoryRepo
var categoryService Category.CategoryService

var apiHandler *apihandler.APIHandler

var roomRepo Room.RoomRepo
var roomService Room.RoomService

var paymentser Payment.PaymentService
var paymentrepo Payment.PaymentRepo

// var roomRpo Room.RoomRepo
// var roomService Room.RoomService

var sectionRepo Section.SectionRepo
var sectionService Section.SectionService

var scheduleRepo Schedule.ScheduleRepo
var scheduleService Schedule.ScheduleService

var infoRepo Information.InformationRepo
var infoService Information.InformationService

var questionRepo Question.QuestionRepo
var questionService Question.QuestionService

var breakRepo BreakDates.BreakDateRepo
var breakService BreakDates.BreakDateService

func startUp() {
	once.Do(
		func() {
			db, erro = DB.InitializPostgres()
			if erro != nil {
				return
			}
			systemTemplate = (template.Must(template.New("dtis").Funcs(handler.FuncMap).ParseGlob("../../web/templates/*.html")))
		},
	)
}
func main() {
	startUp()

	sectionRepo = SectionRepo.NewSectionRepo(db)
	sectionService = SectionService.NewSectionService(sectionRepo)

	breakRepo = BreakDateRepo.NewBreakDateRepo(db)
	breakService = BreakDateService.NewBreakDateService(breakRepo)

	resourceRepo = ResourceRepo.NewResourceRepo(db)
	resourceService = ResourceService.NewResourceService(resourceRepo)

	questionRepo = QuestionRepo.NewQuestionRepo(db)
	questionService = QuestionService.NewQuestionService(questionRepo)

	infoRepo = InformationRepo.NewInfoRepo(db)
	infoService = InformationService.NewInformationService(infoRepo)

	scheduleRepo = ScheduleRepo.NewScheduleRepo(db)
	scheduleService = ScheduleService.NewScheduleService(scheduleRepo)

	paymentrepo = PaymentRepo.NewGormPaymentRepo(db)
	paymentser = PaymentService.NewGormPaymentService(paymentrepo)

	courseRepo = CourseRepo.NewGormCourseRepo(db)
	courseService = CourseService.NewGormCourseService(courseRepo)

	roomRepo = RoomRepo.NewRoomRepo(db)
	roomService = RoomService.NewRoomService(roomRepo)

	roundRepo = RoundRepo.NewRoundRepo(db)
	roudnService = RoundService.NewRoundService(roundRepo)

	vehicleRepo = VehicleRepo.NewVehicleRepo(db)
	vehicleService = VehicleService.NewVehicleService(vehicleRepo)

	categoryRepo = CategoryRepo.NewCategoryRepo(db)
	categoryService = CategoryService.NewCategoryService(categoryRepo)

	trainRepo = TrainerRepo.NewTrainerRepo(db)
	trinServ = TrainerService.NewTrainerService(trainRepo)

	teacherRepo := teachRepo.NewTeacherRepo(db)
	teacherService := teachSer.NewTeacherService(teacherRepo)

	branchrepo := BranchRepo.NewBranchRepo(db)
	branchService := BranchService.NewBranchService(branchrepo)

	sessionHandler = session.NewCookieHandler()

	studentrepo := stdrepo.NewStudentRepository(db)
	studentservice := stdserv.NewStudentService(studentrepo)

	adminrepo := adminre.NewAdminRepo(db)
	adminservice := adminse.NewAdminService(adminrepo)

	th = handler.NewTemplateHandler(systemTemplate, branchService, sessionHandler, adminservice)
	studentHandler := handler.NewStudentHandler(
		studentservice,
		th,
		sessionHandler,
		branchService,
		adminservice,
		roudnService,
	)
	courseHandler := handler.NewCourseHandler(sessionHandler, courseService, roudnService)
	categoryhandler := handler.NewCategoryHandler(categoryService, sessionHandler, roudnService)
	branchhandler := handler.NewBranchHandler(th, sessionHandler, branchService)
	questionHandler := handler.NewQuestionHandler(
		questionService,
		sessionHandler,
		studentservice,
		branchService , 
		th,
	)
	paymenthandler := handler.NewPaymentHandler(
		paymentser,
		sessionHandler,
		adminservice,
		studentservice,
		roudnService,
		branchService,
		th,
	)
	roundhandler := handler.NewRoundHandler(
		roudnService,
		sessionHandler,
		categoryService,
		adminservice,
		trinServ,
		teacherService,
		studentservice,
		courseService,
	)
	adminHandler := handler.NewAdminHandler(
		adminservice,
		th,
		sessionHandler,
		branchhandler,
		&studentservice,
		teacherService,
		trinServ,
		categoryService,
		vehicleService,
		roudnService,
		roomService,
	)
	apiHandler = apihandler.NewAPIHandler(
		vehicleService,
		sessionHandler,
		adminservice,
		studentservice,
		teacherService,
		trinServ,
		roomService,
	)
	userhandler := handler.NewUserHandler(
		adminHandler,
		studentHandler,
		th,
		sessionHandler,
		branchService,
		studentservice,
		teacherService,
		trinServ)
	roomhandler := handler.NewRoomHandler(
		roomService,
		sessionHandler,
	)
	schedulehandler := handler.NewScheduleHandler(
		scheduleService,
		sessionHandler,
		roudnService,
		roomService,
		sectionService,
		// LectureService,
		trinServ,
		teacherService,
		breakService,
	)
	infohandler := handler.NewInfoHandler(
		infoService,
		sessionHandler,
		branchService,
	)

	breakdateshandler := handler.NewBreakDateHandler(
		breakService,
		sessionHandler,
	)
	resourcehandler := handler.NewResourceHandler(
		resourceService,
		sessionHandler,
		branchService,
		adminservice,
		studentservice,
		teacherService,
		trinServ,
		th,
	)

	htmltopdfhandler := handler.NewHTMLToPdfHandler(
		sectionService,
		th,
		sessionHandler,
		branchService,
		roudnService,
	)
	vehiclehandler := handler.NewVehicleHandler(
		vehicleService,
		sessionHandler,
		trinServ,
	)

	trainerhandler := handler.NewTrainerHandler(
		trinServ,
		sessionHandler,
		vehicleService,
	)
	mux := http.NewServeMux() //.StrictSlash(true)
	fs := http.FileServer(http.Dir("../../web/templates/css"))
	mux.Handle("/css/", http.StripPrefix("/css/", neuter(fs)))
	fsjs := http.FileServer(http.Dir("../../web/templates/js"))
	mux.Handle("/js/", http.StripPrefix("/js/", neuter(fsjs)))
	fsimg := http.FileServer(http.Dir("../../web/templates/img"))
	mux.Handle("/img/", http.StripPrefix("/img/", neuter(fsimg)))
	fspublic := http.FileServer(http.Dir("../../web/templates/public"))
	mux.Handle("/public/", http.StripPrefix("/public/", neuter(fspublic)))
	fsfonts := http.FileServer(http.Dir("../../web/templates/fonts"))
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", neuter(fsfonts)))
	fssource := http.FileServer(http.Dir("../../web/templates/Source"))
	mux.Handle("/Source/", http.StripPrefix("/Source/", neuter(fssource)))
	mux.HandleFunc("/", th.IndexPage)
	mux.HandleFunc("/chart/", th.ChartPage)
	mux.HandleFunc("/login", userhandler.Login)                                                                        // Handling Login Request
	mux.HandleFunc("/logout", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(userhandler.Logout)))) // Handling Login Request
	mux.HandleFunc("/user/password/new", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(apiHandler.ChangePassword))))
	mux.HandleFunc("/admin/controll/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(adminHandler.Controll))))
	mux.HandleFunc("/admin/registration/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(adminHandler.SuperadminRegistration))))
	mux.HandleFunc("/admin/new", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(adminHandler.RegisterAdmin))))
	mux.HandleFunc("/admin/teacher/new", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(adminHandler.RegisterTeacher))))
	mux.HandleFunc("/admin/room/new", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(apiHandler.RegisterRoom))))
	mux.HandleFunc("/admin/trainer/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(adminHandler.RegisterFieldMan))))
	mux.HandleFunc("/admin/student/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(adminHandler.SecretaryRegisterStudent))))
	mux.HandleFunc("/api/vehicles/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(apiHandler.GetVehicles))))
	/*	*
		Branch Related Routes
		CreateBranch
		CategoriesOfABranch
		ChangeEmail
		EditPhones
		ChangeLogo
		UpdateBranch
		**/
	mux.HandleFunc("/branch/categories/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(categoryhandler.CategoriesOfABranch)))) // Doesn't Ask Authorization
	mux.HandleFunc("/branch/rooms/", roomhandler.GetRoomsOfBranch)                                                                                  // Doesn't Ask Authorization
	mux.HandleFunc("/branch/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(branchhandler.CreateBranch))))
	mux.HandleFunc("/branch/email/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(branchhandler.ChangeEmail))))
	mux.HandleFunc("/branch/phone/update/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(branchhandler.EditPhones))))
	mux.HandleFunc("/branch/logo/update/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(branchhandler.ChangeLogo))))
	mux.HandleFunc("/branch/update/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(branchhandler.UpdateBranch))))

	// Category Related Routes GetRoundsOfCategory   CreateVehicle DeleteVehicle  FreeVehiclesOfCategory  AssignVehicleForFieldMan  DetachVehicleFromFieldMan
	// ActiveStudentsOfCategory
	mux.HandleFunc("/admin/category/activation/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(categoryhandler.ActivateCategory))))
	mux.HandleFunc("/admin/category/deactivation/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(categoryhandler.DeactivateCategory))))
	mux.HandleFunc("/admin/category/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(categoryhandler.CreateCategory))))
	mux.HandleFunc("/admin/category/update/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(categoryhandler.EditCategory))))
	mux.HandleFunc("/admin/category/delete/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(categoryhandler.DeleteCategoryByID))))
	mux.HandleFunc("/admin/category/rounds/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.GetRoundsOfCategory))))
	mux.HandleFunc("/admin/category/vehicle/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(vehiclehandler.CreateVehicle))))
	mux.HandleFunc("/admin/category/vehicle/delete/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(vehiclehandler.DeleteVehicle))))
	mux.HandleFunc("/admin/category/vehicles/free/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(vehiclehandler.FreeVehiclesOfCategory))))
	mux.HandleFunc("/category/students/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(studentHandler.StudentsOfCategory))))

	// Trainers Related Routes ListOfTrainersOfCategory
	mux.HandleFunc("/category/trainers/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(trainerhandler.ListOfTrainersOfCategory))))
	mux.HandleFunc("/category/trainers/free/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(trainerhandler.ListOfFreeTrainers))))
	// Trainer Related Routes GetTrainerByID
	mux.HandleFunc("/admin/fieldman/vehicle/assign", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(trainerhandler.AssignVehicleForFieldMan))))
	mux.HandleFunc("/admin/fieldman/vehicle/detach", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(trainerhandler.DetachVehicleFromFieldMan))))
	mux.HandleFunc("/admin/fieldman/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(trainerhandler.GetTrainerByID))))

	//Course Related Routings  courseHandler
	mux.HandleFunc("/admin/course/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(courseHandler.CreateCourse))))
	mux.HandleFunc("/api/admin/category/courses/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(courseHandler.GetCoursesByCategoryID))))
	mux.HandleFunc("/api/admin/branch/courses/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(courseHandler.GetCoursesOfABranch))))
	mux.HandleFunc("/api/admin/course/edit/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(courseHandler.EditCourse))))
	mux.HandleFunc("/api/admin/round/course/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(courseHandler.CreateCourseForRound))))

	// GetRoundByID  AddTeacherToRound  RemoveStudentFromRound   PopulateRoundUsingXlsx
	mux.HandleFunc("/admin/round/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.CreateRound))))
	mux.HandleFunc("/admin/round/populate/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.PopulateRoundUsingXlsx))))
	mux.HandleFunc("/round/populator/", roundhandler.GetPassiveRoundPopulatingForm)

	// StudentsIDCardGenerate
	// mux.HandleFunc("/admin/round/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.CreateRound))))
	mux.HandleFunc("/admin/branch/rounds/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.GetRounds))))
	mux.HandleFunc("/api/admin/round/update/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.UpdateRound))))
	mux.HandleFunc("/api/admin/round/trainer/add/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.AddTrainerToRound))))
	mux.HandleFunc("/api/admin/round/teacher/add/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.AddTeacherToRound))))
	mux.HandleFunc("/api/admin/round/course/add/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.AddCourseToRound))))
	mux.HandleFunc("/api/admin/round/student/remove/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.RemoveStudentFromRound))))
	mux.HandleFunc("/api/admin/round/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.GetRoundByID))))
	mux.HandleFunc("/admin/round/course/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.NewCourseToRound))))
	mux.HandleFunc("/admin/round/students/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.GetStudentsOfRound))))
	mux.HandleFunc("/round/students/id/", studentHandler.StudentsIDCardGenerate) //userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(roundhandler.GetStudentsOfRound))))

	// Payments Route   CreatePayment  PaymentsOfRound  PaymentsOfStudent  PaymentsOfSecretary
	mux.HandleFunc("/api/admin/payment/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(paymenthandler.CreatePayment))))
	mux.HandleFunc("/api/admin/round/payments/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(paymenthandler.PaymentsOfRound))))
	mux.HandleFunc("/api/admin/student/payments/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(paymenthandler.PaymentsOfStudent))))
	mux.HandleFunc("/api/admin/secretary/payments/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(paymenthandler.PaymentsOfSecretary))))
	mux.HandleFunc("/admin/payment/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(paymenthandler.ReciptForPayment))))

	// Schedule Related Routes GenerateScheduleForRound GenerateSchedulePDF
	mux.HandleFunc("/api/round/schedule/new", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(schedulehandler.GenerateScheduleForRound))))
	mux.HandleFunc("/round/schedule/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(htmltopdfhandler.GenerateSchedulePDF))))
	mux.HandleFunc("/schedule/remove/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(schedulehandler.DeleteSchedule))))

	// Students related Routes     GetRoundStudentsImageZip
	mux.HandleFunc("/round/students/image/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(studentHandler.GetRoundStudentsImageZip))))

	// Information Related Routes  GetAllInfos  GetInfoByID
	mux.HandleFunc("/api/admin/info/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(infohandler.CreateInformation))))
	mux.HandleFunc("/api/admin/info/update/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(infohandler.UpdateInfo))))
	mux.HandleFunc("/api/admin/info/delete/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(infohandler.DeleteInformation))))
	mux.HandleFunc("/api/admin/info/activation/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(infohandler.ActivateInformation))))
	mux.HandleFunc("/api/admin/info/deactivation/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(infohandler.DeactivateInformation))))
	mux.HandleFunc("/api/info/branch/active/", infohandler.GetActiveInformations)
	mux.HandleFunc("/api/info/active/all/", infohandler.GetAllActiveInformations)
	mux.HandleFunc("/api/info/all/", infohandler.GetAllInfos)
	mux.HandleFunc("/api/info/each/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(infohandler.GetInfoByID))))

	// Question and Answer Related Routes  GetGradeResultForStudent
	mux.HandleFunc("/api/question/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(questionHandler.CreateQuestion))))
	mux.HandleFunc("/api/question/remove/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(questionHandler.DeleteQuestion))))
	mux.HandleFunc("/api/question/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(questionHandler.GetQuestions))))
	mux.HandleFunc("/api/question/answer/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(questionHandler.GetAnswer))))
	mux.HandleFunc("/api/question/result/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(questionHandler.ShowQuestionResult))))
	mux.HandleFunc("/api/question/total_result/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(questionHandler.GetGradeResultForStudent))))

	// Resource Related Routes  UploadResource  DeleteResource  GetResources  SearchResource DownloadResource
	mux.HandleFunc("/learning/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.UploadResource))))
	mux.HandleFunc("/learning/resource/remove/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.DeleteResource))))
	mux.HandleFunc("/learning/resources/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.GetResources))))
	mux.HandleFunc("/learning/resource/videos/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.GetResourceVideos))))
	mux.HandleFunc("/learning/resource/audios/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.GetResourceAudios))))
	mux.HandleFunc("/learning/resource/pictures/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.GetResourcePicture))))
	mux.HandleFunc("/learning/resource/files/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.GetResourceFiles))))
	mux.HandleFunc("/learning/resource/pdfs/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.GetResourcePDFS))))
	mux.HandleFunc("/learning/resource/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.GetResourceByID))))
	mux.HandleFunc("/learning/search/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.SearchResource))))
	mux.HandleFunc("/download/", resourcehandler.DownloadResource)

	// Reserved Dates handler  CreateBreakDatesHandler   DeleteBreakDatesHandler
	//   GetBreakDates     DateInformationHandler

	mux.HandleFunc("/date/new/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(breakdateshandler.CreateBreakDatesHandler))))
	mux.HandleFunc("/date/remove/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(breakdateshandler.DeleteBreakDatesHandler))))
	mux.HandleFunc("/dates/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(breakdateshandler.GetBreakDates))))
	mux.HandleFunc("/date/info/", breakdateshandler.DateInformationHandler)

	mux.HandleFunc("/image/new/", adminHandler.ChangeProfilePicture)
	//  Student Related Routes   TheoryTestingForm  GenerateTreatyPage  RoundRegisteredStudentsInfoPDF
	mux.HandleFunc("/admin/student/biography/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(studentHandler.StudentBiographyPDF))))
	mux.HandleFunc("/admin/student/theory_testing_form/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(studentHandler.TheoryTestingForm))))
	mux.HandleFunc("/admin/student/treaty/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(studentHandler.GenerateTreatyPage))))
	mux.HandleFunc("/admin/student/info/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(htmltopdfhandler.RoundRegisteredStudentsInfoPDF))))

	// GetAdminsOfSystem
	mux.HandleFunc("/admin/delete", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(adminHandler.DeleteAdmin))))  /////////
	mux.HandleFunc("/admins/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(adminHandler.GetAdminsOfSystem)))) /////////
	// mux.HandleFunc("/admin/deactivate/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(adminhandler.DeactivateAdmin))))
	// mux.HandleFunc("/admin/activate/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(adminhandler.ActivateAdmin))))
	// Handling Media Straming
	// MediaRouter.SetTemplate(systemTemplate)
	mux.Handle("/media/", MediaRouter.MediaHandler())
	// Template Pages Route
	mux.HandleFunc("/learning/", userhandler.Authenticated(userhandler.Authorized(http.HandlerFunc(resourcehandler.TemplateLearningPage))))

	// Tests
	mux.HandleFunc("/paymento/", paymenthandler.DailyPaymentReport)

	log.Fatal(http.ListenAndServe(":9900", mux))
}

// DirectoryListener  representing
func DirectoryListener(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
