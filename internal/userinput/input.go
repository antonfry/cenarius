package userinput

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Input(w string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter %s: ", w)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	text = strings.Replace(text, "\r\n", "", -1)
	fmt.Printf("Got: %s ", text)
	return strings.Trim(text, "\n")
}

func InputId() int {
	inputId := Input("Id of secret")
	id, err := strconv.Atoi(inputId)
	if err != nil {
		log.Fatalf("Wrong id: %s", err.Error())
	}
	return id
}
