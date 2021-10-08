package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"strconv"
	"strings"
)

// Calculate score based on GPA
// Rule 54
func getScore(averageGPA float64) int {
	fmt.Println("averageGPA:", averageGPA)
	if averageGPA >= 3.25 {
		return 4
	} else if averageGPA >= 2.0 {
		return int(((8 * averageGPA) - 6) / 5)
	} else {
		return 0
	}
}

func getScoreFromGPA(gpa []string) int {

	var sum float64 = 0

	for _, each := range gpa {
		each := strings.ReplaceAll(each, ",", ".")
		val, err := strconv.ParseFloat(each, 64)
		if err != nil {
			fmt.Println(err)
		} else {
			sum += val
		}
	}
	return getScore(sum / float64(len(gpa)))
}

func readGPA() {
	xlsx, err := excelize.OpenFile("./ipk_lulusan.xlsx")
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	sheet1Name := "Tabel 8.a"

	var gpa []string
	for i := 0; i < 3; i++ {
		gpa = append(gpa, xlsx.GetCellValue(sheet1Name, fmt.Sprintf("D%d", i+4)))
	}

	var score int = getScoreFromGPA(gpa)

	fmt.Printf("score: %d \n", score)
}

func changeGPA(value []string) {
	xlsx, err := excelize.OpenFile("./ipk_lulusan.xlsx")
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	sheet1Name := "Tabel 8.a"

	for i, each := range value {
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+4), each)
	}

	err = xlsx.SaveAs("./ipk_lulusan.xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

func runExampleGPA() {
	changeGPA([]string{"1,5", "1,5", "1,5"})
	readGPA()
	changeGPA([]string{"2,1", "2,1", "2,1"})
	readGPA()
	changeGPA([]string{"3,5", "2,3", "3,5"})
	readGPA()
	changeGPA([]string{"4,0", "3,65", "3,6"})
	readGPA()
}
