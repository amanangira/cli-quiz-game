package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// TODO
// 1.Refactor timeout implementation to leverage timer.C
// 2. Implement select between timeout and answer channel

func main() {
	// Configure CLI flags
	filePath := flag.String("csv", "./questions.csv", "CSV file path")
	timeout := flag.Int("time", 30, "time in seconds") //// time in seconds
	flag.Parse()
	// Open file
	file, fileErr := os.Open(*filePath)
	if fileErr != nil {
		panic(fileErr)
	}
	// Read file
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	questionList, readErr := csvReader.ReadAll()
	if readErr != nil {
		panic(readErr)
	}
	var score int
	var newLine string
	doneCh := make(chan bool)
	answerCh := make(chan bool)
	fmt.Printf("A timer of %d seconds will begin with the first question, are you ready? (y/n)\n", *timeout)
	if _, scanErr := fmt.Scanln(&newLine); scanErr != nil {
		panic(scanErr)
	}

	if strings.ToLower(newLine) != "y" {
		return
	}

	fmt.Println()
	go writeTimeout(doneCh, *timeout)
	go printQuestions(answerCh, questionList, &score)

	select {
	case <-doneCh:
		fmt.Printf("\nTimed out!")
	case <-answerCh:
	}

	fmt.Printf("\nYou scored %d out of %d\n", score, len(questionList))
}

func writeTimeout(c chan<- bool, t int) {
	timer := time.NewTimer(time.Second * time.Duration(t))
	<-timer.C

	c <- true
}

func printQuestions(c chan<- bool, questionList [][]string, score *int) {
	var tempScore int
	var response string
	for index := range questionList {
		partsLength := len(questionList[index])
		questionParts := questionList[index][:partsLength-1]
		question := strings.Join(questionParts, "")
		answer := questionList[index][partsLength-1]

		fmt.Printf("%s ", question)
		_, scanErr := fmt.Scan(&response)
		if scanErr != nil {
			panic(scanErr)
		}
		if answer == response {
			if score == nil {
				*score = tempScore
			} else {
				tempScore = *score + 1
				*score = tempScore
			}
		}
	}

	c <- true
}
