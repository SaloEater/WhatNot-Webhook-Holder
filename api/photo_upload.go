package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (a *API) PhotoUpload(w http.ResponseWriter, r *http.Request) (any, error) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		fmt.Println("An error occurred during parsing multipart form of photo upload: " + err.Error())
		return nil, err
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("An error occurred during reading file of photo upload: " + err.Error())
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("An error occurred during reading file bytes of photo upload: " + err.Error())
		return nil, err
	}

	seriesIDStr := r.FormValue("series_id")
	seriesID, err := strconv.ParseInt(seriesIDStr, 10, 64)
	if err != nil {
		fmt.Println("An error occurred during parsing series_id of photo upload: " + err.Error())
		return nil, err
	}

	name := r.FormValue("name")
	team := r.FormValue("team")

	return a.Service.PhotoUpload(seriesID, data, name, team, header.Filename)
}
