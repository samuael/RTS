package permission

import (
	"net/http"
	"strings"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
)

type permission struct {
	roles   []string
	methods []string
}
type authority map[string]permission

var authorities = authority{
	"/admin/controll/": permission{
		roles:   []string{"SUPERADMIN", "SECRETARY"},
		methods: []string{"GET", "POST"},
	},
	"/admin/registration/": permission{
		roles:   []string{"SUPERADMIN"},
		methods: []string{"GET"},
	},

	"/admin/teacher/new": permission{
		roles:   []string{"SUPERADMIN"},
		methods: []string{"POST"},
	},
	"/admin/new": permission{
		roles:   []string{"SUPERADMIN"},
		methods: []string{"POST"},
	},
	"/admin/fieldman/new": permission{
		roles:   []string{"SUPERADMIN"},
		methods: []string{"POST"},
	},
	"/api/vehicles/": permission{
		roles:   []string{"SUPERADMIN", "SECRETARY"},
		methods: []string{"POST", "GET"},
	},
	"/admin/trainer/new/": permission{
		roles:   []string{"SUPERADMIN"},
		methods: []string{"POST"},
	},
	"/admin/student/new/": permission{
		roles:   []string{entity.SECRETART},
		methods: []string{"POST"},
	},
	"/user/password/new": permission{
		roles:   []string{entity.SUPERADMIN, entity.TEACHER, entity.STUDENT, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/admin/room/new": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/admin/category/activation/": permission{
		roles:   []string{entity.SUPERADMIN},
		methods: []string{"POST", "GET"},
	},
	"/admin/category/deactivation/": permission{
		roles:   []string{entity.SUPERADMIN},
		methods: []string{"POST"},
	},
	"/admin/category/new/": permission{
		roles:   []string{entity.SUPERADMIN},
		methods: []string{"POST"},
	},
	"/admin/category/update/": permission{
		roles:   []string{entity.SUPERADMIN},
		methods: []string{"POST"},
	},
	"/admin/category/delete/": permission{
		roles:   []string{entity.SUPERADMIN},
		methods: []string{"POST"},
	},
	"/logout/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART, entity.STUDENT, entity.FIELDMAN, entity.TEACHER},
		methods: []string{http.MethodGet},
	},
	"/admin/course/new/": permission{
		roles:   []string{entity.SUPERADMIN},
		methods: []string{"POST"},
	}, "/api/admin/category/courses/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{http.MethodGet},
	}, "/api/admin/course/edit/": permission{
		roles:   []string{entity.SUPERADMIN},
		methods: []string{"POST"},
	},
	"/api/admin/round/course/new/": permission{
		roles:   []string{entity.SUPERADMIN},
		methods: []string{"POST"},
	},
	"/api/admin/branch/courses/": permission{
		roles: []string{
			entity.SUPERADMIN, entity.SECRETART,
			entity.TEACHER, entity.FIELDMAN},
		methods: []string{"GET"},
	},
	"/admin/round/populate/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},

	/*
	*
	*   Round Related Routes
	*
	 */
	"/admin/branch/rounds/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART, entity.TEACHER, entity.FIELDMAN},
		methods: []string{"GET"},
	},
	"/admin/round/new/": permission{
		roles:   []string{entity.SUPERADMIN},
		methods: []string{http.MethodPost},
	},
	"/api/admin/round/update/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/admin/round/trainer/add/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/admin/round/teacher/add/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/admin/round/trainer/remove/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/admin/round/": permission{
		roles: []string{
			entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"GET"},
	},
	"/api/admin/round/student/remove/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/admin/round/students/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART, entity.STUDENT, entity.FIELDMAN, entity.TEACHER},
		methods: []string{"GET"},
	},
	"/category/students/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART, entity.STUDENT, entity.FIELDMAN, entity.TEACHER},
		methods: []string{"GET"},
	},
	/*
	 *
	 *
	 *
	 *  Category Related Routes
	 *
	 *
	 */
	"/api/admin/round/course/add/": permission{
		roles:   []string{entity.SUPERADMIN},
		methods: []string{"POST"},
	},
	"/admin/round/course/new/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	}, "/admin/category/rounds/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART, entity.STUDENT, entity.FIELDMAN, entity.TEACHER},
		methods: []string{"GET"},
	},
	"/admin/category/vehicle/new/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/admin/category/vehicle/delete/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/admin/category/vehicles/free/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"GET"},
	},

	// Trainers Related Routes
	"/category/trainers/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART, entity.STUDENT, entity.FIELDMAN, entity.TEACHER},
		methods: []string{"GET"},
	},
	"/category/trainers/free/": permission{
		roles:   []string{entity.SUPERADMIN, entity.FIELDMAN, entity.SECRETART},
		methods: []string{http.MethodGet},
	}, "/trainers/vehicle/add/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{http.MethodPost},
	},
	"/admin/fieldman/vehicle/assign": permission{ /// assigning a vehicle for a trainer
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/admin/fieldman/vehicle/detach": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"GET"},
	},
	"/admin/fieldman/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART, entity.FIELDMAN, entity.TEACHER, entity.STUDENT},
		methods: []string{"GET"},
	},

	/*
		*
		*
		Student Related Routes *
		*
		*
		*
		*
	*/
	"/round/students/image/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"GET"},
	},

	"/admin/student/biography/": permission{
		roles:   []string{entity.SECRETART, entity.SUPERADMIN},
		methods: []string{"GET"},
	}, "/admin/student/theory_testing_form/": permission{
		roles:   []string{entity.SECRETART, entity.SUPERADMIN},
		methods: []string{"GET"},
	},

	"/admin/student/treaty/": permission{
		roles:   []string{entity.SECRETART, entity.SUPERADMIN},
		methods: []string{"GET"},
	},
	"/admin/student/info/": permission{
		roles:   []string{entity.SECRETART, entity.SUPERADMIN},
		methods: []string{"GET"},
	},
	/*
	*  End
	 */

	/*
		Payment Related Routes And Methods
	*/
	"/api/admin/payment/new/": permission{
		roles:   []string{entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/admin/round/payments/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"GET"},
	},
	"/api/admin/student/payments/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{http.MethodPost},
	},
	"/api/admin/secretary/payments/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"GET"},
	},
	"/admin/payment/": permission{
		roles:   []string{entity.SECRETART},
		methods: []string{"GET"},
	},
	/*	Payment Related Routes And Methods
	 */
	// Schedule related routes
	/*
	*
	*
	 */
	"/api/round/schedule/new": permission{ // Activation and Deactivation
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"GET"},
	},
	"/round/schedule/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART, entity.STUDENT, entity.TEACHER, entity.FIELDMAN},
		methods: []string{"GET"},
	},
	"/schedule/remove/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"GET"},
	},
	// Schedule related routes End
	/*
	*
	*
	 */
	// Info Related Paths
	"/api/admin/info/new/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/admin/info/update/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/admin/info/delete/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/admin/info/activation/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/admin/info/deactivation/": permission{
		roles:   []string{entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/info/each/": permission{
		roles:   []string{entity.TEACHER, entity.STUDENT, entity.SECRETART, entity.FIELDMAN, entity.SUPERADMIN},
		methods: []string{"GET"},
	},
	// Question and Answer related routes
	"/api/question/new/": permission{
		roles: []string{
			entity.TEACHER, entity.FIELDMAN, entity.SUPERADMIN, entity.SECRETART,
		},
		methods: []string{"POST"},
	}, "/api/question/remove/": permission{
		roles:   []string{entity.TEACHER, entity.FIELDMAN, entity.SUPERADMIN, entity.SECRETART},
		methods: []string{"POST"},
	},
	"/api/question/": permission{
		roles:   []string{entity.STUDENT, entity.TEACHER, entity.SECRETART, entity.SUPERADMIN, entity.FIELDMAN},
		methods: []string{"GET"},
	},
	"/api/question/answer/": permission{
		roles:   []string{entity.STUDENT},
		methods: []string{"GET"},
	},
	"/api/question/result/": permission{
		roles: []string{
			entity.STUDENT,
		},
		methods: []string{"GET"},
	},
	"/api/question/total_result/": permission{
		roles: []string{
			entity.STUDENT,
		},
		methods: []string{"GET"},
	},
	"/learning/new/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART},
		methods: []string{http.MethodPost},
	},
	"/learning/resources/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART},
		methods: []string{http.MethodGet},
	},
	"/learning/resource/remove/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART},
		methods: []string{http.MethodGet},
	},
	"/learning/resource/videos/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{http.MethodGet},
	},
	"/learning/resource/audios/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{http.MethodGet},
	},
	"/learning/resource/pictures/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{http.MethodGet},
	},
	"/learning/resource/files/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{http.MethodGet},
	},
	"/learning/resource/pdfs/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{http.MethodGet},
	},
	"/learning/resource/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{http.MethodGet},
	}, "/learning/search/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{http.MethodGet},
	},
	"/learning/resource/download/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{http.MethodGet},
	},
	// BreakDates related Handlers Path
	"/date/new/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{"POST"},
	},
	"/date/remove/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{"POST"},
	},
	"/dates/": permission{
		roles:   []string{entity.FIELDMAN, entity.TEACHER, entity.SUPERADMIN, entity.SECRETART, entity.STUDENT},
		methods: []string{"GET"},
	},
	// TemplatePage Related Routes
	"/learning/": permission{
		roles:   []string{entity.TEACHER, entity.STUDENT, entity.SECRETART, entity.FIELDMAN, entity.SUPERADMIN},
		methods: []string{"GET"},
	},

	"/image/new/": permission{
		roles:   []string{entity.TEACHER, entity.STUDENT, entity.SECRETART, entity.FIELDMAN, entity.SUPERADMIN},
		methods: []string{"POST"},
	},
	/*
	*
	* Super Admin Related Functionalities
	*
	 */
	"/admin/delete": permission{
		roles:   []string{entity.OWNER, entity.SUPERADMIN},
		methods: []string{"POST"},
	},
	"/admins/": permission{
		roles:   []string{entity.OWNER},
		methods: []string{http.MethodGet},
	},
	/*
	*
	*  BranchRelated Routes
	*
	 */
	"/branch/new/": permission{
		roles:   []string{entity.OWNER},
		methods: []string{"POST"},
	},
	"/branch/categories/": permission{
		roles:   []string{entity.OWNER, entity.SUPERADMIN},
		methods: []string{"POST"},
	},
	"/branch/email/new/": permission{
		roles:   []string{entity.OWNER},
		methods: []string{"POST"},
	},
	"/branch/phone/update/": permission{
		roles:   []string{entity.OWNER},
		methods: []string{"POST"},
	},
	"/branch/update/": permission{
		roles:   []string{entity.OWNER},
		methods: []string{"POST"},
	},
}

// HasPermission checks if a given role has permission to access a given route for a given method
func HasPermission(path string, role string, method string) bool {
	perm := authorities[path]
	checkedRole := checkRole(role, perm.roles)
	checkedMethod := checkMethod(method, perm.methods)
	if !checkedRole || !checkedMethod {
		return false
	}
	return true
}

func checkRole(role string, roles []string) bool {
	for _, r := range roles {
		if strings.ToUpper(r) == strings.ToUpper(role) {
			return true
		}
	}
	return false
}
func checkMethod(method string, methods []string) bool {
	for _, m := range methods {
		if strings.ToUpper(m) == strings.ToUpper(method) {
			return true
		}
	}
	return false
}
