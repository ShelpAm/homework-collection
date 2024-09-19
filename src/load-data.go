package main

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

func LoadStudents(accounts *map[Student]struct{}) error {
	f, err := excelize.OpenFile("students.xlsx")
	if err != nil {
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}

	for _, row := range rows {
		name := row[0]
		schoolId := row[1]

		(*accounts)[Student{name, schoolId}] = struct{}{}
	}

	return nil
}

func LoadAssignments(assignments *map[string]Assignment) error {
	(*assignments)["第二周"] = Assignment{"第二周", time.Now(), time.Now().Add(time.Hour * 24 * 7)}
	(*assignments)["第三周"] = Assignment{"第三周", time.Now(), time.Now().Add(time.Hour * 24 * 7)}

	return nil
}
