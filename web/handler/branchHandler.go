package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Projects/RidingTrainingSystem/pkg/Helper"
	etc "github.com/Projects/RidingTrainingSystem/pkg/ethiopianCalendar"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Branch"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/Projects/RidingTrainingSystem/pkg/form"
)

// BranchHandler struct
type BranchHandler struct {
	TemplateHandler *TemplateHandler
	SessionHandler  *session.Cookiehandler
	BranchService   Branch.BranchService
}

// NewBranchHandler function
func NewBranchHandler(
	th *TemplateHandler,
	sess *session.Cookiehandler,
	bs Branch.BranchService,
) *BranchHandler {
	return &BranchHandler{
		TemplateHandler: th,
		SessionHandler:  sess,
		BranchService:   bs,
	}
}

// *******************************Methods*********************************************
//  UpdateBranch    ChangeLogo   EditPhones  DeleteBranch  ChangeEmail  CreateBranch
// ***********************************************************************************

// CreateBranch method
// Method POST
// Authorization OWNER
// Input is Multipart Post Request
func (branchhandler *BranchHandler) CreateBranch(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	rsession := branchhandler.SessionHandler.GetSession(request)
	if rsession == nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	ds := struct {
		Success bool
		Message string
		Branch  entity.Branch
		VErrors form.ValidationErrors
	}{
		Success: false,
	}
	errro := request.ParseMultipartForm(9999999999999999)
	if errro != nil {
		ds.Message = "Parsing Error , Invalid Request Body "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	signupForm := form.Input{VErrors: form.ValidationErrors{}, Values: request.PostForm}
	signupForm.Required("name", "country", "region",
		"zone", "woreda",
		"licence_day", "licence_month", "licence_year", "lang",
		// "license_date" ,
		"email", "moto", "branch_name_local", "acronym",
		"branch_name_english", "licence_number",
		"city")
	if !signupForm.Valid() {
		ds.VErrors = signupForm.VErrors
		ds.Message = "Required Field Missing "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	phones := []string{}
	var values url.Values
	values = request.PostForm
	theArray := values["phone"]
	for _, vaue := range theArray {
		// if len(vaue) > 20 {
		// 	signupForm.VErrors.Add("Phones ", "Invalid Phone Number ")
		// }
		phones = append(phones, vaue)
	}
	branchImage, header, err := request.FormFile("image")
	var branchLogo *os.File
	var branchLogoError error
	var imageUrl string
	if err == nil && header != nil {
		filename := header.Filename
		fileextension := Helper.GetExtension(filename)
		randomName := fmt.Sprintf("%sbranch_logo_%s.%s", entity.PathToBranchImagesFromTemplates, Helper.GenerateRandomString(4, Helper.NUMBERS), fileextension)
		branchLogoRelDirectory := fmt.Sprintf("%s%s", entity.PathToTemplates, randomName)
		branchLogo, branchLogoError = os.Create(branchLogoRelDirectory)
		if branchLogoError != nil {
			imageUrl = ""

		} else {
			imageUrl = randomName
			defer branchLogo.Close()
		}
		defer branchImage.Close()
	}
	var licenseDay, licenseMonth, licenseYear int
	licenseDay, era := strconv.Atoi(request.FormValue("licence_day"))
	licenseMonth, era = strconv.Atoi(request.FormValue("licence_month"))
	licenseYear, era = strconv.Atoi(request.FormValue("licence_year"))
	if era != nil {
		signupForm.VErrors.Add("LicenseDate", "License Date Values Are Not Entered Correctly ")
	}
	if !signupForm.Valid() {
		ds.VErrors = signupForm.VErrors
		ds.Message = "Required Field Missing "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	branch := &entity.Branch{
		Name:    request.FormValue("name"),
		Country: request.FormValue("country"),
		Address: entity.Address{
			Region: request.FormValue("region"),
			Zone:   request.FormValue("zone"),
			Woreda: request.FormValue("woreda"),
			Kebele: request.FormValue("kebele"),
			City:   request.FormValue("city"),
		},
		LicenceGivenDate: etc.Date{
			Day:   (licenseDay),
			Month: (licenseMonth),
			Year:  (licenseYear),
		},
		Email:                 request.FormValue("email"),
		BranchAcronym:         request.FormValue("acronym"),
		BranchFullnameAmharic: request.FormValue("branch_name_local"),
		BranchFullnameEnglish: request.FormValue("branch_name_english"),
		LicenceNumber:         request.FormValue("licence_number"),
		City:                  request.FormValue("city"),
		Lang:                  request.FormValue("lang"),
		Logourl:               imageUrl,
		Phones:                phones,
		Moto:                  request.FormValue("moto"),
		Createdby:             rsession.Username,
	}
	branch = branchhandler.BranchService.CreateBranch(branch)
	if branch == nil {
		signupForm.VErrors.Add("DTIS : ERROR ", "Internal Server Error")
	}
	if !signupForm.Valid() {
		ds.VErrors = signupForm.VErrors
		ds.Message = "Required Field Missing "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	if branchImage != nil || imageUrl != "" {
		io.Copy(branchLogo, branchImage)
	}
	ds.Success = true
	ds.Message = "Succesful"
	ds.Branch = *branch
	jr, _ := json.Marshal(ds)
	response.Write(jr)
}

// ChangeEmail method to change the email of the Branch
// Methos POST
// Authorizatio for Owner
func (branchhandler *BranchHandler) ChangeEmail(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	ds := struct {
		Success  bool
		Message  string
		NewEmail string
		BranchID int
	}{
		Success: false,
	}
	email := request.FormValue("email")
	branchID, era := strconv.Atoi(request.FormValue("branch_id"))
	if era != nil || branchID <= 0 {
		ds.Message = "Invalid Input Request "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	ds.BranchID = branchID
	ds.NewEmail = email
	success := branchhandler.BranchService.ChangeEmail(branchID, email)
	if success {
		ds.Success = success
		ds.Message = "Branch Email Succesfuly Changed"
	} else {
		ds.Message = "Request Was Not Succesful "
	}
	jr, _ := json.Marshal(ds)
	response.Write(jr)
}

//DeleteBranch   only owner should Do
//  method post
// variable branch_id
// since this is  a highly Risky Action it should have _csrf Check
func (branchhandler *BranchHandler) DeleteBranch(response http.ResponseWriter, request *http.Request) { // Not Finished
	response.Header().Set("Content-Type", "application/json")
	request.ParseForm()
	rsession := branchhandler.SessionHandler.GetSession(request)
	if rsession == nil || rsession.Role != entity.OWNER {
		return
	}
	ds := struct {
		Success  bool
		Message  string
		BranchID int
	}{
		Success: false,
	}
	branchID, era := strconv.Atoi(request.FormValue("branch_id"))
	if era != nil || branchID <= 0 {
		ds.Message = "Invalid Input Number "
		ds.BranchID = branchID
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	ds.BranchID = branchID
	success := branchhandler.BranchService.DeleteBranch(branchID)
	if success {
		ds.Success = success
		ds.Message = "Branch Succesfully Deleted "
	} else {
		ds.Message = fmt.Sprintf(" Deleting Branch  With ID : %d Was Not Succesful\n", branchID)
	}
	jr, _ := json.Marshal(ds)
	response.Write(jr)
}

// EditPhones method to add a phone number to the Branch Phone Numbers List
// Method Post
// Authorization  OWNER
// Variable phone list of phons
func (branchhandler *BranchHandler) EditPhones(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	// rsession := branchhandler.SessionHandler.GetSession(request)
	// if rsession != nil {
	// 	response.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	ds := struct {
		Success  bool
		Message  string
		BranchID int
		Phones   []string
	}{
		Success: false,
	}
	branchID, era := strconv.Atoi(request.FormValue("branch_id"))
	if era != nil || branchID <= 0 {
		ds.Message = "   Invalid Branch ID Value "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	var values url.Values
	values = request.PostForm
	theArray := values["phone"]
	if len(theArray) == 0 {
		ds.Message = "Invalid Length Of array "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}

	ds.BranchID = branchID
	success := branchhandler.BranchService.ChangePhones(branchID, theArray)
	if success {
		ds.Message = "Succesful"
		ds.Success = true
		ds.Phones = theArray
	} else {
		ds.Message = "Change Phone Numbers Was No Succesful"
	}
	jr, _ := json.Marshal(ds)
	response.Write(jr)
}

// ChangeLogo method to change the logo of the bRanch
// Authorization  OWNER
// variable image
// response   the Path to Profiel Picture
func (branchhandler *BranchHandler) ChangeLogo(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	era := request.ParseMultipartForm(990999999999)
	ds := struct {
		Success  bool
		Message  string
		Logourl  string
		BranchID uint
	}{
		Success: false,
	}
	branchID, era := strconv.Atoi(request.FormValue("branch_id"))
	if era != nil || branchID <= 0 {
		ds.Message = "Invalid Branch ID "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	branch := branchhandler.BranchService.GetBranchByID(uint(branchID))
	if branch == nil {
		ds.Message = "No Branch Found this ID "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	ds.BranchID = uint(branchID)
	var imageURl string
	image, header, headerError := request.FormFile("image")
	if headerError == nil && header != nil {
		extension := Helper.GetExtension(header.Filename)
		if branch.Logourl == "" {
			randomname := fmt.Sprintf("%sbranch_logo_%s.%s", entity.PathToBranchImagesFromTemplates, Helper.GenerateRandomString(4, Helper.CHARACTERS), extension)
			fmt.Println(randomname)
			branch.Logourl = randomname
		}
		imageURl = branch.Logourl
		imageURl = entity.PathToTemplates + imageURl
		newImage, era := os.Create(imageURl)
		_, era = io.Copy(newImage, image)
		if era != nil {
			ds.Message = "Internal Server Error "
			jr, _ := json.Marshal(ds)
			response.Write(jr)
			return
		}
		ds.Success = true
		ds.Message = "Succesfully Upgraded the Image"
		ds.Logourl = branch.Logourl
		// saving the Branch after Changing the pahth of the logo
		branch = branchhandler.BranchService.CreateBranch(branch)
		if branch == nil {
			ds.Message = "Internal Server ERROR while Saving the Updated Data \n Please Try Again "
			ds.Success = false
		}
	} else {
		ds.Message = "Logo Upload Not Succesful"
		ds.Logourl = branch.Logourl
	}
	jr, _ := json.Marshal(ds)
	response.Write(jr)
}

// UpdateBranch for Updating some Datas Especially Those which are not frequently Changing
// Method GEt
// AUTHORIZATION OWNER
// Variables branch_id
// there Will Be name  , country  , moto  , acronym
//, city  , licence_number , local_branch_name  , english_branch_name
func (branchhandler *BranchHandler) UpdateBranch(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	rsession := branchhandler.SessionHandler.GetSession(request)
	if rsession == nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	ds := struct {
		BranchID uint
		Success  bool
		Message  string
		Branch   entity.Branch
	}{
		Success: false,
	}
	branchID, era := strconv.Atoi(request.FormValue("branch_id"))
	if era != nil || branchID <= 0 {
		ds.Message = fmt.Sprintf("Invalid Branch Id : %d ", branchID)
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	ds.BranchID = uint(branchID)
	branch := branchhandler.BranchService.GetBranchByID(uint(branchID))
	if branch == nil {
		ds.Message = "Invalid Request"
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	changes := 0
	branch.Lang = IsChanged(branch.Lang, request.FormValue("lang"), &changes)
	branch.Country = IsChanged(branch.Country, request.FormValue("country"), &changes)
	branch.Moto = IsChanged(branch.Moto, request.FormValue("moto"), &changes)
	branch.BranchAcronym = IsChanged(branch.BranchAcronym, request.FormValue("acronym"), &changes)
	branch.City = IsChanged(branch.City, request.FormValue("city"), &changes)
	branch.LicenceNumber = IsChanged(branch.LicenceNumber, request.FormValue("licence_number"), &changes)
	branch.BranchFullnameAmharic = IsChanged(branch.BranchFullnameAmharic, request.FormValue("local_branch_name"), &changes)
	branch.BranchFullnameEnglish = IsChanged(branch.BranchFullnameEnglish, request.FormValue("english_branch_name"), &changes)
	branch.Address.Region = IsChanged(branch.Address.Region, request.FormValue("region"), &changes)
	branch.Address.Woreda = IsChanged(branch.Address.Woreda, request.FormValue("woreda"), &changes)
	branch.Address.Kebele = IsChanged(branch.Address.Kebele, request.FormValue("kebele"), &changes)
	branch.Address.Zone = IsChanged(branch.Address.Zone, request.FormValue("zone"), &changes)
	branch.Address.City = IsChanged(branch.Address.City, request.FormValue("city"), &changes)
	if changes > 0 {
		branch = branchhandler.BranchService.CreateBranch(branch)
		if branch == nil {
			ds.Message = "Internal Server Error"
		} else {
			ds.Success = true
			ds.Message = "Succesfuly Updated "
			ds.Branch = *branch
		}
	} else {
		ds.Message = "No Change Is Made "
	}
	jr, _ := json.Marshal(ds)
	response.Write(jr)
	return
}

// IsChanged function
func IsChanged(previousValue, newValue string, changes *int) string {
	if newValue != previousValue && strings.Trim(newValue, " ") != "" {
		(*changes)++
		return newValue
	}
	return previousValue
}

// Branches returning the Branches in the Database
func (branchhandler *BranchHandler) Branches(resposne http.ResponseWriter, request *http.Request) {
	resposne.Header().Set("Content-Type", "application/json")
	ds := struct {
		Success  bool
		Message  string
		Branches []entity.Branch
	}{
		Success: false,
	}

	branches := branchhandler.BranchService.GetBranchs()
	if branches == nil {
		ds.Message = " Not Brancs Found Please Create One "
		ds.Branches = []entity.Branch{}
		resposne.Write(Helper.MarshalThis(ds))
		return
	}
	ds.Message = fmt.Sprintf(" Succesfuly Found %d Brnches", len(*branches))
	ds.Branches = *branches
	resposne.Write(Helper.MarshalThis(ds))
}

// UpdateLicenseGivenDate method for uploading changing the Date In Which The Company Got The License
// and if the company has many branches using their branch_id
// Method POST
// AUTHORIZATION OWNER and variables     day   , month  , year  , branch_id
func (branchhandler *BranchHandler) UpdateLicenseGivenDate(response http.ResponseWriter, request *http.Request) {
	session := branchhandler.SessionHandler.GetSession(request)
	if session == nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	ds := struct {
		Success  bool
		Message  string
		BranchID uint
		Date     etc.Date
	}{
		Success: false,
	}
	day, era := strconv.Atoi(request.FormValue("day"))
	month, era := strconv.Atoi(request.FormValue("month"))
	year, era := strconv.Atoi(request.FormValue("year"))
	branchID, era := strconv.Atoi(request.FormValue("branch_id"))
	if era != nil || day <= 0 || day > 30 || month <= 0 || month > 13 || year <= 1900 || year > int(etc.ETYear())+1 || branchID <= 0 {
		ds.Message = "Invalid Data Input >Try AGAIN!"
		response.Write(Helper.MarshalThis(ds))
		return
	}
	branch := branchhandler.BranchService.GetBranchByID(uint(branchID))
	if branch == nil {
		ds.Message = "  No Branch Found By this ID  "
		response.Write(Helper.MarshalThis(ds))
		return
	}
	lastDate := branch.LicenceGivenDateRefer
	branchhandler.BranchService.DeleteDate(int64(lastDate))
	ds.BranchID = uint(branchID)
	date := etc.Date{
		Day:   day,
		Month: month,
		Year:  year,
	}
	branch.LicenceGivenDate = date
	branch = branchhandler.BranchService.CreateBranch(branch)
	if branch == nil {
		ds.Message = "  No Branch Found By this ID  "
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(Helper.MarshalThis(ds))
		return
	}
	ds.Success = true
	ds.Message = " Branch Licence Given Data Succesfully Updated"
	ds.Date = branch.LicenceGivenDate
	response.Write(Helper.MarshalThis(ds))
}
