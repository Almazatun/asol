package prompt

import (
	"fmt"
	"log"

	"github.com/Almazatun/asol/constant"
	"github.com/gagliardetto/solana-go/rpc"
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

func SelectNetworkPrompt() string {
	prompt := promptui.Select{
		Label: "Please select network" + constant.QUESTION_PROMPT_EXIT_PART,
		Items: []string{"Mainnet", "Devnet"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if result == "Mainnet" {
		return rpc.MainNetBeta_RPC
	}

	return rpc.DevNet_RPC
}

func SubDirPrompt(homeDir string) string {
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("please enter subdirectory path ..., $HOME=%v/...", homeDir),
	}

	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("prompt failed %v\n", err)
	}

	return result
}
