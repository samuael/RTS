package MediaRouter

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Projects/RidingTrainingSystem/internal/pkg/entity"
	"github.com/gorilla/mux"
)

// MediaHandler function to handle Video and Audio Streaming
func MediaHandler() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/media/{mID:[0-999999999]+}/stream/", StreamHandler).Methods("GET")
	router.HandleFunc("/media/{mID:[0-999999999]+}/stream/{segName:index[0-9999]+.ts}", StreamHandler).Methods("GET")
	return router
}

// StreamHandler for Handling streaming related Routes
func StreamHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("I get called .... ")
	vars := mux.Vars(request)
	fmt.Println("The Variables are  ... ", vars)
	mID, era := strconv.Atoi(vars["mID"])
	if era != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	SegName, ok := vars["segName"]
	fmt.Println(vars)
	path := GetMediaBase(uint64(mID))
	fmt.Println(ok)
	if !ok {
		ServeHLSM3u8(response, request, path)
	} else {
		ServeHlsTS(response, request, path, SegName)
	}
}

// GetMediaBase function
func GetMediaBase(mID uint64) string {
	return fmt.Sprintf("%smedia/%d/", entity.PathToResources, mID)
}

// ServeHLSM3u8  function for handling m3
func ServeHLSM3u8(response http.ResponseWriter, request *http.Request, path string) {
	fileName := "index.m3u8"
	theFile := fmt.Sprintf("%s%s", path, fileName)
	response.Header().Set("Content-Type", "application/x-mpegURL")
	http.ServeFile(response, request, theFile)
}

// ServeHlsTS  function for handling m3
func ServeHlsTS(response http.ResponseWriter, request *http.Request, path string, segName string) {
	theFile := fmt.Sprintf("%s%s", path, segName)
	fmt.Println("The File Directory is : ",theFile)
	response.Header().Set("Content-Type", "video/MP2T")
	http.ServeFile(response, request, theFile)
}

// ServeContent
/*
	ServeContent replies to the request
	using the content in the provided
	ReadSeeker. The main benefit of ServeContent
	 over io.Copy is that it handles Range requests properly,
	  sets the MIME type, and handles
	  If-Match, If-Unmodified-Since,
	  If-None-Match, If-Modified-Since,
	   and If-Range requests.

*/
