package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {

	testFile := flag.String("fileName", "problems.csv", "a CSV file in the format of 'question,answer' (default \"problems.csv\")")
	testTime := flag.Int("timer", 30, "the time limit for the quiz in seconds (default 30)")
	flag.Parse()

	rows := readFileToStruct(*testFile)
	printQuestionsAndTakeAnswers(rows, *testTime)
}

func readFileToStruct(testFileName string) []problem {

	file, err := os.Open(testFileName)
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	fileReader := csv.NewReader(file)
	records, err := fileReader.ReadAll()
	if err != nil {
		log.Fatal("Error reading csv file into records")
	}

	problems := make([]problem, len(records))

	for index, row := range records {
		problems[index] = problem{
			question: row[0],
			answer:   row[1],
		}
	}
	return problems
}

func printQuestionsAndTakeAnswers(questions []problem, timeLimit int) {
	var correctAnswers = 0
	var totalQuestions = len(questions)

	fmt.Println("Questions appear one by one, please answer correspondingly")
	fmt.Printf("\n")
	fmt.Print("Press ENTER key to start")
	fmt.Scanln()

	testTimer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	for index, row := range questions {

		fmt.Printf("Question #%d: %v = ", index+1, row.question)
		answerCh := make(chan string)
		go func() {
			var userAnswer string
			fmt.Scanf("%s", &userAnswer)
			answerCh <- userAnswer
		}()

		select {
		case <-testTimer.C:
			fmt.Println("time is up!")
			return
		case answer := <-answerCh:
			if answer == row.answer {
				correctAnswers++
			}
		}
	}
	fmt.Printf("\nYou have scored %d out of %d\n", correctAnswers, totalQuestions)
}
