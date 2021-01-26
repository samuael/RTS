package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/Projects/RidingTrainingSystem/pkg/Helper"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"
	"github.com/Projects/RidingTrainingSystem/pkg/form"
)

// This file is created only to handle and Perform ROLE="SECRETARY" things

// SecretaryRegisterStudent method for secretary only niggoye
func (adminh *AdminHandler) SecretaryRegisterStudent(response http.ResponseWriter, request *http.Request) {

	response.Header().Add("Content-Type", "application/json")
	newSession := adminh.SessionHandler.GetSession(request)
	admin := adminh.AdminService.GetAdminByID(newSession.ID, newSession.Username)
	if admin == nil {
		fmt.Println(admin)
	}
	branch := adminh.BranchHandler.BranchService.GetBranchByID(newSession.BranchID)
	if request.Method == http.MethodGet {
		NewInput := &entity.Informing{
			Admin:    *admin,
			Branch:   *branch,
			CSRF:     adminh.SessionHandler.RandomToken(),
			HasError: false,
			Success:  false,
		}
		jsonReturn, _ := json.Marshal(NewInput)
		response.Write(jsonReturn)
		return
	} else if request.Method == http.MethodPost {
		// fmt.Println("InThe Nothing ")
		// Check whether the Request has valid CSRF or not
		request.ParseMultipartForm(123434545344)
		if request.FormValue("_csrf") != "" {
			if !adminh.SessionHandler.ValidateForm(request.FormValue("_csrf")) {
				adminh.Templatehandler.PageNotFound(response, request)
				return
			}
		}
		// ---------------------------------------------------
		signUpMethod := form.Input{
			Values:  request.PostForm,
			VErrors: form.ValidationErrors{},
		}
		signUpMethod.Required(
			"firstname", "lastname", "grandname", "accstatus", "sex", "birth-day", "birth-month", "birth-year", "maritial-status", "phone",
			"family-count", "guarantor-fullname", "guarantor-phone", "previous-licence-type", "previous-licence-number", "categoryid",
			"roundid", "lang", "student-kebele", "student-woreda", "student-zone", "student-region", "guarantor-kebele",
			"guarantor-woreda", "guarantor-zone", "guarantor-region")
		// The mendatory files are  Correctly Posted therefor now i am going to check the profilepic thing
		imagedir := "img/Students/"
		imageurl := ""
		file, head, fileError := request.FormFile("profilepic")
		if fileError != nil {
			log.Println("Error while fetching the file ")
			imageurl = ""
		}
		imageurl = Helper.GenerateRandomString(10, Helper.CHARACTERS)
		if head != nil {
			if head.Filename != "" {
				// fmt.Println("../../web/templates/" + imagedir + imageurl + head.Filename)
				// appripriateFileName := Helper.JPEGFileName(imageurl + head.Filename)
				name := Helper.JPEGFileName(head.Filename)
				filenew, fileError := os.Create("../../web/templates/" + imagedir + imageurl + name)
				imageurl = imagedir + imageurl + name
				if fileError != nil {
					log.Println("Error while saving the  file ")
					imageurl = ""
					signUpMethod.VErrors.Add("Image", "Error While Saving The Image")
				}
				defer func() {
					file.Close()
					filenew.Close()
				}()
				_, err := io.Copy(filenew, file)
				if err != nil {
					imageurl = ""
				}
			}
		} else {
			signUpMethod.VErrors.Add("Profile Picture", "Profile Picture Must Have to Be Submitted ")
		}
		if !signUpMethod.Valid() {
			signUpMethod.CSRF = adminh.SessionHandler.RandomToken()
			NewInput := &entity.Informing{
				CSRF:     adminh.SessionHandler.RandomToken(),
				HasError: true,
				Input:    signUpMethod,
				Admin:    *admin,
				Branch:   *branch,
				Host:     request.Host,
				Message:  "Please Fille the Required Infos Correctly ",
			}
			jsonReturn, _ := json.Marshal(NewInput)
			response.Write(jsonReturn)
			return
		}
		// Populating the Studdent Struct for the Sake of Saving to the Database
		day, a := strconv.Atoi(request.FormValue("birth-day"))
		month, a := strconv.Atoi(request.FormValue("birth-month"))
		year, a := strconv.Atoi(request.FormValue("birth-year"))
		familyQuantity, fc := strconv.Atoi(request.FormValue("family-count"))
		Categoryid, ci := strconv.Atoi(request.FormValue("categoryid"))
		Roundid, ri := strconv.Atoi(request.FormValue("roundid"))
		if ri != nil {
			signUpMethod.VErrors.Add("Round Id", "Round ID Has to be An Integer Type Not a text")
		}
		if ci != nil {
			signUpMethod.VErrors.Add("CategoryID", "FieldCategoryId Has to Be Integer")
		}
		if fc != nil {
			signUpMethod.VErrors.Add("FamilyCount", " Family Count Has To Be Number not A Text")
		}
		if a != nil {
			signUpMethod.VErrors.Add("General", "Incorrect Date Value Birth Date ")
		}
		if !signUpMethod.Valid() {
			signUpMethod.CSRF = adminh.SessionHandler.RandomToken()
			NewInput := &entity.Informing{
				CSRF:     adminh.SessionHandler.RandomToken(),
				HasError: true,
				Success:  false,
				Input:    signUpMethod,
				Admin:    *admin,
				Branch:   *branch,
				Host:     request.Host,
				Message:  "Please Fille the Required Infos Correctly ",
			}
			jsonReturn, _ := json.Marshal(NewInput)
			response.Write(jsonReturn)
			// adminh.Templatehandler.Templates.ExecuteTemplate(response, "sam_registration.html", NewInput)
			return
		}
		birthDate := &etc.Date{
			Day:   (day),
			Month: (month),
			Year:  (year),
		}
		// "student-kebele", "student-woreda", "student-zone"
		studentAddress := entity.Address{
			Kebele: request.FormValue("student-kebele"),
			Zone:   request.FormValue("student-zone"),
			Woreda: request.FormValue("student-woreda"),
			Region: request.FormValue("student-region"),
		}
		guarantorAddress := entity.Address{
			Kebele: request.FormValue("guarantor-kebele"),
			Zone:   request.FormValue("guarantor-zone"),
			Woreda: request.FormValue("guarantor-woreda"),
			Region: request.FormValue("guarantor-region"),
		}
		round := adminh.RoundService.GetRoundByID(uint(Roundid))
		if round == nil {
			NewInput := &entity.Informing{
				CSRF:     adminh.SessionHandler.RandomToken(),
				HasError: true,
				Success:  false,
				Input:    signUpMethod,
				Admin:    *admin,
				Branch:   *branch,
				Host:     request.Host,
				Message:  " No Such Round  ",
				// Specific: student,
			}
			jsonreturn, _ := json.Marshal(NewInput)
			response.Write(jsonreturn)
			return
		}
		student := &entity.Student{
			AcademicStatus:        request.FormValue("accstatus"),
			Firstname:             request.FormValue("firstname"),
			Lastname:              request.FormValue("lastname"),
			GrandFatherName:       request.FormValue("grandname"),
			Sex:                   request.FormValue("sex"),
			Nickname:              request.FormValue("nickname"),
			Lang:                  request.FormValue("lang"),
			MaritialStatus:        request.FormValue("maritial-status"),
			PhoneNumber:           request.FormValue("phone"),
			BirthDate:             birthDate,
			PartnerPhoneNumber:    request.FormValue("partner-phone"),
			PartnerFullname:       request.FormValue("partner-fullname"),
			GuarantorAddress:      &guarantorAddress,
			GuarantorFullName:     request.FormValue("guarantor-fullname"),
			GuarantorPhoneNumber:  request.FormValue("guarantor-fullname"),
			PreviousLicenceType:   request.FormValue("previous-licence-type"),
			PreviousLicenceNumber: request.FormValue("previous-licence-number"),
			CategoryID:            uint(Categoryid),
			FamilyCount:           uint(familyQuantity),
			RoundRefer:            uint(Roundid),
			Address:               &studentAddress,
			Imageurl:              imageurl,
			Round:                 round,
			BranchID:              newSession.BranchID,
			Username: round.Category.Title +
				"/" +
				strconv.Itoa(int(round.Roundnumber)) + "/" +
				strconv.Itoa(int(round.Studentscount+1)) + "/" +
				strconv.Itoa(int(etc.NewDate(0).Year)),
			// Category :
			Password: Helper.GenerateRandomString(4, Helper.NUMBERS),
		}
		student.Round = round
		round.Students = append(round.Students, *student)
		round.Studentscount = uint(len(round.Students))
		round = adminh.RoundService.SaveRound(round)
		if round == nil {
			fmt.Println("Error While Saving Round ...........")
		}
		student = adminh.StudentService.SaveStudent(student)
		if student == nil || round == nil {
			NewInput := &entity.Informing{
				CSRF:     adminh.SessionHandler.RandomToken(),
				HasError: true,
				Success:  false,
				Input:    signUpMethod,
				Admin:    *admin,
				Branch:   *branch,
				Host:     request.Host,
				Message:  " Error While Saving the Data ",
				Specific: student,
			}
			jsonreturn, _ := json.Marshal(NewInput)
			response.Write(jsonreturn)
		}
		admin.Branch = entity.Branch{}
		NewInput := &entity.Informing{
			CSRF:     adminh.SessionHandler.RandomToken(),
			HasError: false,
			Success:  true,
			Input:    signUpMethod,
			Admin:    *admin,
			Branch:   *branch,
			Host:     request.Host,
			Message:  "Succesfully Registered the student  ",
			Specific: student,
		}
		jsonreturn, _ := json.Marshal(NewInput)
		response.Write(jsonreturn)
	}
}

// DeleteAdmin  function
// Method POST
// Variables admin id and if the Admin is SUPERADMIN the authority will be only for OWNER
// if the ADmin is Secretary the authority is going given to superadmin where the branch of the the secretary resides
func (adminh *AdminHandler) DeleteAdmin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	rsession := adminh.SessionHandler.GetSession(request)
	if rsession == nil || !(rsession.Role == entity.OWNER || rsession.Role == entity.SUPERADMIN) {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	ds := struct {
		Success bool
		Message string
		AdminID uint
	}{
		Success: false,
	}
	adminID, era := strconv.Atoi(request.FormValue("admin_id"))
	if era != nil || adminID <= 0 {
		ds.Message = "Invalid Request Body "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	admin := adminh.AdminService.GetAdminByID(uint(adminID), "")
	ds.AdminID = uint(adminID)
	if admin == nil {
		ds.Message = "Record Not Found "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	if admin.BranchRefer != rsession.BranchID || (strings.ToUpper(admin.Role) == entity.SUPERADMIN && rsession.Role != entity.OWNER) {
		ds.Message = "You Are Not Authorized "
		jsonReturn, _ := json.Marshal(ds)
		response.Write(jsonReturn)
		return
	}
	success := adminh.AdminService.DeleteAdmin(uint(adminID))
	if !success {
		ds.Message = "Error While Deleting the Admin"
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	profileImage := admin.Imageurl
	if profileImage != "" {
		profileImage = entity.PathToTemplates + profileImage
		os.Remove(profileImage)
	}
	ds.Message = "Succesfuly Deleted "
	ds.Success = true
	ds.AdminID = uint(adminID)
	jr, _ := json.Marshal(ds)
	response.Write(jr)
}

// GetAdminsOfSystem  for fetching the superadmins of the system
// Variables branchid -1 ,0 ,1  status entity.All entity.Active entity.Passive
// Method Get
// Authorization for superAdmin Only
func (adminh *AdminHandler) GetAdminsOfSystem(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	ds := struct {
		Success      bool
		Message      string
		Admins       []entity.Admin
		BranchFilter int
		StatusFilter int
	}{
		Success: false,
	}
	branchID, era := strconv.Atoi(request.FormValue("branch_id"))
	statusNo, era := strconv.Atoi(request.FormValue("status"))
	status := int(statusNo)
	fmt.Printf(" Branch ID : %d \n Status : %d\n", branchID, status)
	if era != nil || !(status == entity.All || status == entity.Active || status == entity.Passive) || branchID < -1 {
		ds.Message = "Invalid Request Body "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	admins := adminh.AdminService.GetAdminsOfSystem(branchID, status)
	ds.BranchFilter = int(branchID)
	ds.StatusFilter = int(status)
	if admins == nil {
		ds.Message = fmt.Sprintf("No Admin Found By this Filter Branch %d Status %d ", branchID, status)
	} else {
		ds.Message = "succesful"
		ds.Success = true
		ds.Admins = *admins
	}
	jr, _ := json.Marshal(ds)
	response.Write(jr)
}

// AdminsActivation to activate or deactivate admin
// methdo get
// variable admin id
// response json ds
// Authorization ALL
func (adminh *AdminHandler) AdminsActivation(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	rsession := adminh.SessionHandler.GetSession(request)
	if rsession == nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	ds := struct {
		Success bool
		Message string
		AdminID uint
		Active  bool
	}{
		Success: false,
	}
	adminID, era := strconv.Atoi(request.FormValue("admin_id"))
	// action, era := strconv.ParseBool(request.FormValue("activate"))
	if era != nil || adminID < 0 {
		ds.Message = "Invalid Input "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	admin := adminh.AdminService.GetAdminByID(uint(adminID), "")
	if admin == nil {
		ds.Message = "Record Not Found "
		ds.AdminID = uint(adminID)
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	if admin.Active {
		admin.Active = false
		ds.Message = fmt.Sprintf("Admin : %s ID : %d Deactivated Succesfully ", admin.Username, admin.ID)
		ds.Success = true
		ds.Active = false
		success := adminh.AdminService.DeactivateAdmin(uint(adminID))
		if !success {
			ds.Success = false
			ds.Message = "Error While Deactivating Admin"
		}
	} else {
		admin.Active = true
		ds.Message = fmt.Sprintf("Admin : %s id : %d Activated Succesfully ", admin.Username, admin.ID)
		ds.Active = true
		ds.Success = true
		success := adminh.AdminService.ActivateAdmin(uint(adminID))
		if !success {
			ds.Success = false
			ds.Message = "Error While Activating Admin "
		}
	}
	jr, _ := json.Marshal(ds)
	response.Write(jr)
}
