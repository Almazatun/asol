package info

import (
	"context"
	"fmt"

	"github.com/Almazatun/asol/constant"
	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/manifoldco/promptui"
)

const (
	pubKeyQuestion = "Please enter your public key"
)

func GetAccountInfo(args []string) error {
	client := rpc.New(rpc.MainNetBeta_RPC)

	pubKeyPrompt := promptui.Prompt{
		Label: pubKeyQuestion + constant.QUESTION_PROMPT_EXIT_PART,
	}

	result, err := pubKeyPrompt.Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("prompt failed %v\n", err))
	}

	accountPubKey, err := solana.PublicKeyFromBase58(result)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to parse public key %v\n", result))
	}

	resp, err := client.GetAccountInfo(
		context.TODO(),
		accountPubKey,
	)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to get account info %v\n", err))
	}
	spew.Dump(resp)

	// var mint token.Mint
	// err = client.GetAccountDataInto(
	// 	context.TODO(),
	// 	accountPubKey,
	// 	&mint,
	// )
	// if err != nil {
	// 	return fmt.Errorf(fmt.Sprintf("failed to get account info %v\n", err))
	// }

	// spew.Dump(mint)
	return nil
}
