package handlers

import (
	"regexp"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
	"fmt"
)

const addProjectRegexp = "^ *add project (\\w.*?) *$"

const projectNameMinimumLength = 4

func handleAddProject(uid, text string) {
	r := regexp.MustCompile(addProjectRegexp)
	projectName := r.FindStringSubmatch(text)[1]

	_, err := models.FindProjectByNameOrAlias(projectName)

	if err == nil {
		sender.SendMessage(uid, fmt.Sprintf("Project with name \"%s\" already exists.", projectName))
		return
	} else if _, ok := err.(models.NotFoundError); !ok {
		handleError(uid, err)
		return
	}

	if len(projectName) < projectNameMinimumLength {
		sender.SendMessage(uid, fmt.Sprintf("The minimum lenth for a project's name is %d.", projectNameMinimumLength))
		return
	}

	p := models.Project{}

	p.Name = projectName

	err = p.Create()

	if err != nil {
		handleError(uid, err)
		return
	}

	sender.SendMessage(uid, fmt.Sprintf("The project with name \"%s\" was successfully created.", p.Name))
}
