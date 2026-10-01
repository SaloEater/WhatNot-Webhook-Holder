package api

import (
	"encoding/json"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
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

	var thumbData []byte
	thumbFilename := ""
	if tf, th, err := r.FormFile("thumbnail"); err == nil {
		defer tf.Close()
		thumbData, err = io.ReadAll(tf)
		if err != nil {
			return nil, err
		}
		thumbFilename = th.Filename
	}

	seriesIDStr := r.FormValue("series_id")
	seriesID, err := strconv.ParseInt(seriesIDStr, 10, 64)
	if err != nil {
		fmt.Println("An error occurred during parsing series_id of photo upload: " + err.Error())
		return nil, err
	}

	name := r.FormValue("name")
	team := r.FormValue("team")
	price, _ := strconv.ParseInt(r.FormValue("price"), 10, 64)
	rotation, _ := strconv.ParseInt(r.FormValue("rotation"), 10, 64)

	var artBounds, labelBounds *entity.Bounds
	if v := r.FormValue("art_bounds"); v != "" {
		artBounds = &entity.Bounds{}
		if err := json.Unmarshal([]byte(v), artBounds); err != nil {
			return nil, fmt.Errorf("invalid art_bounds: %w", err)
		}
	}
	if v := r.FormValue("label_bounds"); v != "" {
		labelBounds = &entity.Bounds{}
		if err := json.Unmarshal([]byte(v), labelBounds); err != nil {
			return nil, fmt.Errorf("invalid label_bounds: %w", err)
		}
	}
	var labelText *entity.LabelText
	if v := r.FormValue("label_text"); v != "" {
		var err error
		if labelText, err = entity.ParseLabelText(v); err != nil {
			return nil, fmt.Errorf("invalid label_text: %w", err)
		}
	}

	return a.Service.PhotoUpload(seriesID, data, name, team, price, rotation, header.Filename, thumbData, thumbFilename, artBounds, labelBounds, labelText)
}
