package handlers

import (
	"github.com/alex-bogomolov/timebot_go/models"
	"fmt"
	"github.com/alex-bogomolov/timebot_go/sender"
)

const showProjectsRegexp = "^ *projects *$"

func handleShowProjects(uid string) {
	projects, err := models.GetAllProjects()

	if err != nil {
		handleError(uid, err)
	}

	stringArray := models.NewStringArray()

	largestLength := longestProjectNameLength(projects)

	for _, project := range projects {
		stringArray.Add(fmt.Sprintf("%s Alias: %s", rightPad(project.Name, largestLength), project.Alias.String))
	}

	sender.SendMessage(uid, fmt.Sprintf("```%s```", stringArray.Join("\n")))
}


func rightPad(s string, length int) string {
	out := s

	for len(out) < length {
		out += " "
	}

	return out
}

func longestProjectNameLength(projects []*models.Project) int {
	max := 0

	for _, project := range projects {
		if len(project.Name) > max {
			max = len(project.Name)
		}
	}

	return max
}
