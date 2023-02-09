package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

type CandidatesStruct struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

func ExelParser(fileName string) []CandidatesStruct {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(ex)
	candidates := []CandidatesStruct{}
	var pathDirectory string
	pathDirectory = path + "/" + fileName
	f, err := excelize.OpenFile(pathDirectory)
	if err != nil {
		fmt.Println(err)
		return candidates
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	sheetName := f.GetSheetName(0)
	rows, _ := f.GetRows(sheetName)
	for idx, row := range rows {
		if idx == 0 {
			continue
		}
		candidate := CandidatesStruct{
			Name:     row[0],
			Email:    row[1],
			Password: row[2],
			Phone:    row[3],
		}
		candidates = append(candidates, candidate)
	}
	fmt.Println(candidates)
	return candidates
}
