package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

/*Problem represents a pair of question answer from the quiz database*/
type Problem struct {
	question, answer string
}

func loadProblems(filename string) (problems []Problem, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	csvr := csv.NewReader(f)
	rows, err := csvr.ReadAll()
	if err != nil {
		return nil, err
	}
	problems = make([]Problem, len(rows))
	for index, row := range rows {
		problems[index] = Problem{strings.TrimSpace(row[0]), strings.TrimSpace(row[1])}
	}
	return problems, nil
}

func quiz(problem *Problem, answerCh chan bool) {
	fmt.Printf("%s?\n", problem.question)
	fmt.Printf("Answer:")
	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Printf("Fatal Error: Something went wrong! %s", err)
	}
	if strings.TrimSpace(answer) == problem.answer {
		answerCh <- true
	} else {
		answerCh <- false
	}
}

func quizLoop(problems []Problem, seconds int) int {
	score := 0
	timer := time.NewTimer(time.Duration(seconds) * time.Second)
	answerCh := make(chan bool)
	for index, problem := range problems {
		fmt.Printf("Question %d:\t", index+1)
		go quiz(&problem, answerCh)

		select {
		case <-timer.C:
			fmt.Printf("\n\n Times up! \n\n")
			return score

		case answer := <-answerCh:
			if answer {
				score++
			}
		}
	}
	return score
}

func printScore(score int, total int) {
	fmt.Println("----------------------------------------")
	fmt.Println("Quiz Done! Here are the results")
	fmt.Printf("You answered %d questions correct out of %d\n", score, total)
	fmt.Println("----------------------------------------")
}

func main() {
	var fileName string
	var timerSeconds int
	flag.StringVar(&fileName, "csv", "problems.csv", "The quiz dataset in csv file")
	flag.IntVar(&timerSeconds, "timerSeconds", 20, "The timer value for the quiz")
	flag.Parse()

	problems, err := loadProblems(fileName)
	if err != nil {
		fmt.Printf("Fatal error while loading csv file : %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("All ready for the quiz, you will have %d seconds to answer all questions!\n", timerSeconds)
	fmt.Println("Press 'Enter' to start the quiz!")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	score := quizLoop(problems, timerSeconds)
	printScore(score, len(problems))
}
