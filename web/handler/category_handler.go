package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Projects/RidingTrainingSystem/pkg/Helper"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/Category"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/Round"
	session "github.com/Projects/RidingTrainingSystem/internal/pkg/Session"
	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/Projects/RidingTrainingSystem/pkg/form"
)

// CategoryHandler struct
type CategoryHandler struct {
	CategoryService Category.CategoryService
	SessionHandler  *session.Cookiehandler
	RoundService    Round.RoundService
}

// NewCategoryHandler function
func NewCategoryHandler(cateservice Category.CategoryService, session *session.Cookiehandler, round Round.RoundService) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: cateservice,
		RoundService:    round,
		SessionHandler:  session,
	}
}

// CreateCategory method for saving a category having an Image returning a json DAta Type
func (catehandler *CategoryHandler) CreateCategory(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	signUpForm := form.Input{
		VErrors: form.ValidationErrors{},
	}
	sessionOfPage := catehandler.SessionHandler.GetSession(request)
	_ = request.ParseMultipartForm(1324364758886)
	var imagePath string
	imagePath = ""
	// if erro != nil {
	// 	signUpForm.VErrors.Add("Parse Error", "Error While Parsing the Data Please Re Submit the Form ")
	// }
	signUpForm.Values = request.PostForm
	signUpForm.Required("title", "_csrf")
	returnValue := struct {
		Success          bool
		Message          string
		ValidationErrors form.ValidationErrors
		Category         entity.Category
	}{
		Success: false,
	}
	if !signUpForm.Valid() {
		returnValue.Message = "Creating the Category Was Not Succesful "
		returnValue.ValidationErrors = signUpForm.VErrors
		jsonReturn, _ := json.Marshal(
			returnValue,
		)
		response.Write(jsonReturn)
		// DATASTRUCTURE FOR CONTROLL PAGE
		//  Route to the Controll Page
	}
	CSRF := request.FormValue("_csrf")
	if !catehandler.SessionHandler.ValidateForm(CSRF) {
		fmt.Println("Erorro Over ere r")
		http.Redirect(response, request, "/page/Not/Found/", 300)
		return
	}
	//  Reading th einage if there
	image, header, errors := request.FormFile("image")
	if errors == nil && header != nil {
		title := Helper.GenerateRandomString(10, Helper.CHARACTERS) + header.Filename
		imageHost := "../../web/templates/CategoryImages"
		main := "img/CategoryImages"
		newImage, errro := os.Create(imageHost + title)
		if errro != nil {
			signUpForm.VErrors.Add("imageError", "Internal Server Error")
		}
		defer newImage.Close()
		_, Erroa := io.Copy(newImage, image)
		if Erroa != nil {
			signUpForm.VErrors.Add("imageError", "Internal Server Error")
			returnValue.Message = " Error While Saving the Message Please try Again  "
			returnValue.ValidationErrors = signUpForm.VErrors
			jsonReturn, _ := json.Marshal(
				returnValue,
			)
			response.Write(jsonReturn)
			//  Route to other page
		}
		imagePath = main + title
	}
	if image != nil {
		defer image.Close()
	}

	category := &entity.Category{
		Branchid:        sessionOfPage.BranchID,
		ImageURL:        imagePath,
		LastRoundNumber: 0,
		Title:           request.FormValue("title"),
	}
	fmt.Println(category.Branchid, "Branch Id Of A course ")
	category = catehandler.CategoryService.CreateCategory(category)
	if category != nil {
		returnValue.Success = true
		returnValue.Message = " Succesfully Created New Category "
		returnValue.ValidationErrors = signUpForm.VErrors
		returnValue.Category = *category
		jsonReturn, _ := json.Marshal(
			returnValue,
		)
		response.Write(jsonReturn)
	}
	returnValue.Success = false
	returnValue.Message = " Internal Server Error Please Try Again "
	returnValue.ValidationErrors = signUpForm.VErrors
	jsonReturn, _ := json.Marshal(
		returnValue,
	)
	response.Write(jsonReturn)
}

// DeleteCategoryByID method for Deleting the Category
func (catehandler *CategoryHandler) DeleteCategoryByID(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	sessionData := catehandler.SessionHandler.GetSession(request)
	newDecoder := json.NewDecoder(request.Body)
	// DataStructure For DeleteCategory Page
	datastructure := struct {
		CategoryID uint `json:"categoryid"`
		Success    bool
		Message    string
		DeletedBy  string
	}{
		Success:   false,
		Message:   "Can't Delete The Category ",
		DeletedBy: sessionData.Username,
	}
	category := &entity.Category{}
	decodeError := newDecoder.Decode(&datastructure)
	if decodeError != nil {
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	category.ID = datastructure.CategoryID
	fmt.Println(datastructure.CategoryID, "In The Handler ")
	result := catehandler.CategoryService.DeleteCategoryByID(category.ID)
	if result {
		datastructure.Success = true
		datastructure.CategoryID = category.ID
		datastructure.Message = "Succesfully Deleted The Category "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	datastructure.CategoryID = category.ID
	if category.ID <= 0 {
		datastructure.Message = " Invalid Category Number  "
	} else {
		datastructure.Message = "Can't Delete the Category "
	}
	jsonReturn, _ := json.Marshal(datastructure)
	response.Write(jsonReturn)
	return
}

// EditCategory method to edit the Category returning a Json and taking a json
//  Authority for entity.SUPERADMINS only
//  iNput Method could Be Using the Form ajax json nigga
func (catehandler *CategoryHandler) EditCategory(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	datastructure := struct {
		Success  bool
		Message  string
		Category entity.Category
	}{
		Success: false,
	}
	request.ParseMultipartForm(7878787878787878)
	signUpForm := form.Input{
		VErrors: form.ValidationErrors{},
		Values:  request.PostForm,
	}
	signUpForm.Required("id", "title")
	signUpForm.MinLength("title", 2)
	if !signUpForm.Valid() {
		datastructure.Success = false
		datastructure.Message = "Please Fill the DataStructures Correctly "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	ID, erroa := strconv.Atoi(request.FormValue("id"))
	if erroa != nil {
		datastructure.Success = false
		datastructure.Message = "Invalid ID Number !"
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	category := &entity.Category{
		Title: request.FormValue("title"),
	}
	category.ID = uint(ID)
	gettingCategoryByID := catehandler.CategoryService.GetCategoryByID(uint(ID))
	if gettingCategoryByID == nil {
		datastructure.Message = "No Category By the specified ID !"
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	image, header, newError := request.FormFile("newImage")
	var newImage string
	newImage = ""
	if newError == nil && header != nil {
		defer image.Close()
		imageHost := "../../web/templates/img/CategoryImages"
		main := "img/CategoryImages"
		title := Helper.GenerateRandomString(10, Helper.CHARACTERS) + header.Filename
		newImageFile, era := os.Create(imageHost + title)
		if era != nil {
			datastructure.Message = "Internal Server Error "
			jsonReturn, _ := json.Marshal(datastructure)
			response.Write(jsonReturn)
			return
		}
		defer newImageFile.Close()
		io.Copy(newImageFile, image)
		newImage = main + title
		// Closing the Previous File(image )
		erra := os.Remove(imageHost + gettingCategoryByID.ImageURL)
		if erra != nil {
			datastructure.Message = "Internal Server Error Please Try Again "
			jsonReturn, _ := json.Marshal(datastructure)
			response.Write(jsonReturn)
			return
		}
	}
	category.ImageURL = newImage
	gettingCategoryByID.ImageURL = newImage
	gettingCategoryByID.Title = category.Title
	// Saving the New Category struct
	category = catehandler.CategoryService.SaveCategory(gettingCategoryByID)
	if category != nil {
		datastructure.Success = true
		datastructure.Message = "Category is Succesfuly Updated !"
		datastructure.Category = *category
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	datastructure.Message = "Internal erver Error !"
	datastructure.Category = *category
	jsonReturn, _ := json.Marshal(datastructure)
	response.Write(jsonReturn)
	return
}

// ActivateCategory method
func (catehandler *CategoryHandler) ActivateCategory(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requestSession := catehandler.SessionHandler.GetSession(request)
	datastructure := struct {
		ID        uint `json:"categoryid"`
		Succesful bool
		Message   string
	}{
		Succesful: false,
	}
	newDecoder := json.NewDecoder(request.Body)
	newDecodeError := newDecoder.Decode(&datastructure)
	if newDecodeError != nil {
		fmt.Println("Decode Error In the Decoder Class ")
		datastructure.Message = "Decoder Error "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	if datastructure.ID == 0 {
		datastructure.Message = " Invalid Category Number  "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	category := &entity.Category{
		Branchid: requestSession.BranchID,
	}
	category.ID = datastructure.ID
	// Getting the Full Category form The Service
	category = catehandler.CategoryService.GetCategoryByID(category.ID)
	if category == nil {
		datastructure.Message = " Internal Server Erorr "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	if category.Branchid != requestSession.BranchID {
		datastructure.Message = " Sorry You Can Not Authorized To Activate the Category "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	if category.Active {
		datastructure.Message = " The Category was Active "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	category.Active = true
	category = catehandler.CategoryService.SaveCategory(category)
	if category == nil {
		datastructure.Message = "Internal Server Error "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	datastructure.Succesful = true
	datastructure.Message = "Succesfully Activated "
	jsonReturn, _ := json.Marshal(datastructure)
	response.Write(jsonReturn)
	return
}

// DeactivateCategory method
func (catehandler *CategoryHandler) DeactivateCategory(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	requestSession := catehandler.SessionHandler.GetSession(request)
	datastructure := struct {
		ID        uint `json:"categoryid"`
		Succesful bool
		Message   string
	}{
		Succesful: false,
	}
	newDecoder := json.NewDecoder(request.Body)
	newDecodeError := newDecoder.Decode(&datastructure)
	if newDecodeError != nil {
		fmt.Println("Decode Error In the Decoder Class ")
		datastructure.Message = "Decoder Error "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	if datastructure.ID == 0 {
		datastructure.Message = " Invalid Category Number  "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}

	category := &entity.Category{
		Branchid: requestSession.BranchID,
	}
	category.ID = datastructure.ID

	// Getting the Full Category form The Service
	category = catehandler.CategoryService.GetCategoryByID(category.ID)
	if category == nil {
		datastructure.Message = " Internal Server Erorr "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	if category.Branchid != requestSession.BranchID {
		datastructure.Message = " Sorry You Can Not Authorized To Deactivate the Category "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	if !category.Active {
		datastructure.Message = " The Category was Deactivated "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	category.Active = false

	category = catehandler.CategoryService.SaveCategory(category)
	if category == nil {
		datastructure.Message = "Internal Server Error "
		jsonReturn, _ := json.Marshal(datastructure)
		response.Write(jsonReturn)
		return
	}
	datastructure.Succesful = true
	datastructure.Message = "Succesfully Deactivated "
	jsonReturn, _ := json.Marshal(datastructure)
	response.Write(jsonReturn)
	return
}

// CategoriesOfABranch  method to find the Categories of a branch
// Method GET
// Authorization SUPERADMIN and OWNER of the System
func (catehandler *CategoryHandler) CategoriesOfABranch(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")
	rsession := catehandler.SessionHandler.GetSession(request)
	if rsession == nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	ds := struct {
		Success    bool
		Message    string
		BranchID   uint
		Categories []entity.Category
	}{
		Success: false,
	}
	branchID, era := strconv.Atoi(request.FormValue("branch_id"))
	if era != nil || branchID <= 0 {
		ds.Message = "Invalid Input Nigga "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	if strings.ToUpper(rsession.Role) != entity.OWNER && branchID != int(rsession.BranchID) {
		ds.Message = " You Are Not Authorized For this Request ! You Are Branch Limited "
		jr, _ := json.Marshal(ds)
		response.Write(jr)
		return
	}
	ds.BranchID = uint(branchID)
	categories := catehandler.CategoryService.GetCategories(uint(branchID))
	if len(categories) == 0 {
		ds.Message = "No Category Found"
	} else {
		ds.Success = true
		ds.Message = fmt.Sprintf(" Succesfully Found %d Categories ", len(categories))
	}
	ds.Categories = categories
	jr, _ := json.Marshal(ds)
	response.Write(jr)
}
