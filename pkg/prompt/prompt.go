package prompt

import (
	"fmt"
	"log"
	"strings"

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
		Label: fmt.Sprintf("please enter subdirectory path where your keygen file, $HOME=%v/", homeDir),
	}

	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("prompt failed %v\n", err)
	}

	return result
}

// Get balance from by using private key or keygen file
func FromPKOrKGFPrompt() string {
	prompt := promptui.Select{
		Label: "Please select option to continue command" + constant.QUESTION_PROMPT_EXIT_PART,
		Items: []string{"get private key from keygen file", "enter private key"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if strings.Contains(result, "keygen") {
		return "KGF"
	}

	return "PK"
}
