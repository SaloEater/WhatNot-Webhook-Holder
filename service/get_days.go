package service

import (
	"encoding/json"
	"errors"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"os"
	"path/filepath"
)

const dataDir = "data"
const daysFile = "days.json"

var pwd string

func GetDays() ([]byte, error) {
	var daysJSON []byte
	var err error

	err = ensureFileExists(dataDir, daysFile)
	if err != nil {
		return nil, errors.New("Error while ensuring days file existence: " + err.Error())
	}

	daysJSON, err = os.ReadFile(getFilepath(dataDir, daysFile))
	if err != nil {
		return nil, errors.New("Error while reading days file: " + err.Error())
	}

	return daysJSON, nil
}

func getFilepath(dir string, file string) string {
	return filepath.Join(getPwd(), dir, file)
}

func ensureFileExists(dir string, file string) error {
	path := getFilepath(dir, file)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		emptyDays := entity.Days{Days: []entity.Day{}}
		emptyDaysData, err := json.Marshal(emptyDays)
		if err != nil {
			return err
		}
		return os.WriteFile(path, emptyDaysData, 0644)
	}

	return nil
}

func getPwd() string {
	if pwd == "" {
		pwd, _ = os.Getwd()
	}

	return pwd
}
