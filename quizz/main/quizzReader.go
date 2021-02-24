package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Quizz struct {
	mu             sync.Mutex
	totalQuestions int
	rightAnswers   int
}

func (q *Quizz) printFinalResult() {
	fmt.Printf("%d answer correct of %d\n", q.rightAnswers, q.totalQuestions)
}

func (q *Quizz) addRightAnswer() {
	q.mu.Lock()
	q.rightAnswers = q.rightAnswers + 1
	q.mu.Unlock()
}

func programKickOff() {
	var wg sync.WaitGroup
	fmt.Print("Type Y if you want to start the quizz: ")
	kickoffReader := bufio.NewReader(os.Stdin)
	text, _ := kickoffReader.ReadString('\n')
	if strings.TrimSpace(text) == "Y" {
		quizz := Quizz{}
		wg.Add(2)
		go timer(&quizz)
		go execQuizz(&quizz)
		wg.Wait()
	} else {
		os.Exit(0)
	}
}

func execQuizz(quizz *Quizz) {
	file, err := os.Open("problems.csv")
	if err != nil {
		log.Fatalf("Error opening file %s", "problems.csv")
	}
	csvDataReader := csv.NewReader(file)
	csvRows, _ := csvDataReader.ReadAll()
	quizz.totalQuestions = len(csvRows)
	inputReader := bufio.NewReader(os.Stdin)
	for _, row := range csvRows {
		fmt.Printf("%s: ", row[0])
		text, inputErr := inputReader.ReadString('\n')
		fmt.Println(text)
		if inputErr != nil {
			log.Fatalln("can not read user answer")
		}
		if row[1] == strings.TrimSpace(text) {
			quizz.addRightAnswer()
		}
	}
	quizz.printFinalResult()
}

func timer(quizz *Quizz) {
	now := time.Now()
	tick := <-time.NewTicker(10 * time.Second).C
	if now.Before(tick) {
		fmt.Println("Time out, exiting")
		quizz.printFinalResult()
		os.Exit(0)
	}
}

func main() {
	programKickOff()
}
