package balance

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
)

const (
	privateKeyQuestion = "Please enter your private key, or Cntrl+C to exit"
)

func GetBalance(args []string) {
	privateKeyPrompt := promptui.Prompt{
		Label:    privateKeyQuestion,
		Validate: validatePrivateKey,
	}

	result, err := privateKeyPrompt.Run()

	if err != nil {
		log.Printf("Prompt failed %v\n", err)
		return
	}

	pk, _ := solana.PrivateKeyFromBase58(result)

	endpoint := selectNetworkPrompt()
	client := rpc.New(endpoint)

	out, err := client.GetBalance(
		context.TODO(),
		pk.PublicKey(),
		rpc.CommitmentFinalized,
	)

	if err != nil {
		panic(err)
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	t.SetCaption("Account SOL balance")

	t.AppendHeader(table.Row{"PublicKey", "SOL"})
	lamportsOnAccount := new(big.Float).SetUint64(uint64(out.Value))
	t.AppendRow(table.Row{pk.PublicKey(), new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))})

	fmt.Println(t.Render())
}

func validatePrivateKey(input string) error {
	_, err := solana.PrivateKeyFromBase58(input)

	if err != nil {
		return errors.New("Invalid private key")
	}

	return nil
}

func selectNetworkPrompt() string {
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
