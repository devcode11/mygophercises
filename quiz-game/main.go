package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	csvFileName := flag.String("csv", "problems.csv", "a csv `file` in format of 'question,answer'")
	quizTime := flag.Int("timer", 5, "timer for quiz in `seconds`")
	flag.Parse()

	csvFile, err := os.Open(*csvFileName)
	exitIfErr(err)
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	exitIfErr(err)

	problems := parseRecords(records)
	var score int8
	ansChan := make(chan string)

	/*timer := time.AfterFunc(time.Duration(*quizTime)*time.Second, func() {
		fmt.Println("Timer expired")
		fmt.Printf("\nYour score is %d out of %d\n", score, len(problems))
		os.Exit(0)
	})*/

	fmt.Printf("quiz time: %d seconds. Timer started\n", *quizTime)
	timer := time.NewTimer(time.Duration(*quizTime) * time.Second)

problemloop:
	for i, problem := range problems {
		fmt.Printf("Question #%d: %s: ", i+1, problem.q)
		go func() {
			var scannedAns string
			fmt.Scanln(&scannedAns)
			ansChan <- scannedAns
		}()
		select {
			case answer := <-ansChan: {
				if (answer == problem.a) {
					score++
				}
			}			
			case <-timer.C: {
				break problemloop
			}
		}
	}
	fmt.Printf("\nYour score is %d out of %d\n", score, len(problems))
	timer.Stop()
}

type problem struct {
	q string
	a string
}

func parseRecords(records [][]string) []problem {
	problems := make([]problem, len(records))
	for i, record := range records {
		problems[i] = problem{q: record[0], a: record[1]}
	}
	return problems
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
