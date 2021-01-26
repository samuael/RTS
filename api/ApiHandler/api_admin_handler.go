package ApiHandler

import (
	"github.com/Projects/Inovide/Admin"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
)

// APIAdminHandler struct
type APIAdminHandler struct {
	AdminService   Admin.AdminService
	SessionService *session.Cookiehandler
}

// TODO:  Remember Me if any Type Has Logged in it Has To Save its Role in the Session
