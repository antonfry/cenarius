package userinput

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Input(w string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter %s: ", w)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Error(err)
		return ""
	}
	text = strings.Replace(text, "\r\n", "", -1)
	fmt.Printf("Got: %s ", text)
	return strings.Trim(text, "\n")
}

func InputID() int {
	InputID := Input("Id of secret")
	id, err := strconv.Atoi(InputID)
	if err != nil {
		log.Errorf("Wrong id: %s", err.Error())
		return -1
	}
	return id
}
