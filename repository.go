package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
)

func getCommits(
	repository *git.Repository,
	fileName string,
	config BockConfig,
) []Revision {
	o, err := exec.Command(
		"git",
		"-C",
		config.articleRoot,
		"log",
		`--pretty=format:'{"id": "%H", "shortId": "%h", "subject": "%f", "body": "%b", "date": "%aD"}'`,
		fileName,
	).Output()

	res := []Revision{}

	if err != nil {
		fmt.Println("Error getting commits:", err)
		return res
	} else {
		s := string(o)
		s = strings.ReplaceAll(s, "'\n'", ",")
		s = strings.ReplaceAll(s, "'", "")
		s = strings.ReplaceAll(s, "\n", "")
		s = "[" + s + "]"

		json.Unmarshal([]byte(s), &res)
	}

	return res
}
