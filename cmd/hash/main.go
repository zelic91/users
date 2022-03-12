package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"zelic91/users/shared"
)

func main() {
	fmt.Println("[ GENERATE HASH FOR SCORE ]")
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter score: ")
	scoreString, _ := reader.ReadString('\n')
	scoreString = strings.TrimSuffix(scoreString, "\n")

	score, err := strconv.ParseInt(scoreString, 10, 64)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print("Enter time: ")
	timeString, _ := reader.ReadString('\n')
	timeString = strings.TrimSuffix(timeString, "\n")

	hash := shared.GetHashFromScoreRequest(score, timeString)

	fmt.Printf("Hash: %s\n", hash)
}
