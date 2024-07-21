package helper

import (
	"log"

	"github.com/manifoldco/promptui"
)

func YesOrNoPromptByQuestion(question string) string {
	prompt := promptui.Select{
		Label: question,
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return result
}
