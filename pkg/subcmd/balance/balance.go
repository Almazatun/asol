package balance

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/Almazatun/asol/helper"
	"github.com/Almazatun/asol/pkg/prompt"
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
		Validate: helper.ValidatePrivateKey,
	}

	result, err := privateKeyPrompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	// https://github.com/gagliardetto/solana-go/issues/232
	pk, _ := solana.PrivateKeyFromBase58(result)

	endpoint := prompt.SelectNetworkPrompt()
	client := rpc.New(endpoint)

	out, err := client.GetBalance(
		context.TODO(),
		pk.PublicKey(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		log.Fatalf("Get balance failed %v\n", err)
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	t.SetCaption("Account SOL balance")

	t.AppendHeader(table.Row{"PublicKey", "SOL"})
	lamportsOnAccount := new(big.Float).SetUint64(uint64(out.Value))
	t.AppendRow(table.Row{pk.PublicKey(), new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))})

	fmt.Println(t.Render())
}
