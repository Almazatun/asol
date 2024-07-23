package prompt

import (
	"log"

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
		Label: "Please select network",
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
