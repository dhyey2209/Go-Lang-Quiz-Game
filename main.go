package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func problemPuller(fileName string) ([]problem, error) {
	// Open the file

	if fObj, err := os.Open(fileName); err == nil {
		csvR := csv.NewReader(fObj)

		if cLines, err := csvR.ReadAll(); err == nil {
			return parseProblem(cLines), nil
		} else {
			return nil, fmt.Errorf("error in reading data in csv"+"format from %s file; %s", fileName, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file; %s", fileName, err.Error())
	}
}

func main() {
	// Input the name of the file
	fName := flag.String("f", "quiz.csv", "path of csv file")

	// Set the duration of the timer
	timer := flag.Int("t", 30, "timer of the quiz")
	flag.Parse()

	// calling problem puller function
	problems, err := problemPuller(*fName)

	// Handle the error
	if err != nil {
		exit(fmt.Sprintf("something went wrong:%s", err.Error()))
	}
	// Create a variabe to count our current answers
	correctAnswer := 0

	// Using the duration of the timer,initialize the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)

	// Loop through the problems, print the questions, we'll accept the answers
problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s = \n ", i+1, p.q)

		go func() {
			fmt.Scanf("%s\n", &answer)
			ansC <- answer
		}()
		select {
		case <-tObj.C:
			fmt.Println()
			break problemLoop
		case iAns := <-ansC:
			if iAns == p.a {
				correctAnswer++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}
	}
	// Calculate and print out the result
	fmt.Printf("Your result is %d out of %d\n", correctAnswer, len(problems))
	fmt.Printf("Press enter to exit")
	<-ansC
}

func parseProblem(lines [][]string) []problem {
	// Go over the lines and parse them
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r

}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
