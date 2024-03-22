package service

import (
	"fmt"
	"math/rand"
	"os"
)

const BreakFilePostfix = "break"

type GetBreakRequest struct {
	Year  int
	Month int
	Day   int
	Name  string
}

func GetBreak(r *GetBreakRequest) ([]byte, error) {
	breakFilepath := getFilepath(dataDir, createBreakFilename(r.Year, r.Month, r.Day, r.Name))
	breakData, err := os.ReadFile(breakFilepath)
	if err != nil {
		return nil, err
	}

	return breakData, nil
}

func createBreakFilename(year int, month int, day int, name string) string {
	return fmt.Sprintf("%d_%d_%d_%s", year, month, day, getBreakPostfix(name))
}

func createDeletedBreakFilename(year int, month int, day int, name string) string {
	return fmt.Sprintf("%d_%d_%d_%s.%s_deleted", year, month, day, getBreakPostfix(name), RandStringBytes(8))
}

func getBreakPostfix(name string) string {
	return fmt.Sprintf("%s_%s.json", BreakFilePostfix, name)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
