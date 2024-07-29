package balance

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"

	"github.com/Almazatun/asol/constant"
	"github.com/Almazatun/asol/helper"
	"github.com/Almazatun/asol/pkg/prompt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
)

const (
	privateKeyQuestion = "Please enter your private key"
)

var keygenPath string
var privateKey solana.PrivateKey

func GetBalance(args []string) error {
	endpoint := prompt.SelectNetworkPrompt()
	client := rpc.New(endpoint)

	selectOptToGetBalance := fromPKOrKGFPrompt()

	if selectOptToGetBalance == "KGF" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("failed to get user's home directory: %v", err))
		}

		subDir := prompt.SubDirPrompt(homeDir)
		if err := helper.CheckSubDirExists(homeDir, subDir); err != nil {
			return err
		}

		keygenPath = filepath.Join(homeDir, subDir)
		pk, err := solana.PrivateKeyFromSolanaKeygenFile(keygenPath)

		if err != nil {
			return fmt.Errorf(fmt.Sprintf("keygen file read failed %v\n", err))
		}

		privateKey = pk
	} else {
		privateKeyPrompt := promptui.Prompt{
			Label: privateKeyQuestion + constant.QUESTION_PROMPT_EXIT_PART,
		}

		result, err := privateKeyPrompt.Run()
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("prompt failed %v\n", err))
		}

		pk, err := solana.PrivateKeyFromBase58(result)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("parse private key failed %v\n", err))
		}

		// validate private key
		pk.Sign(pk)

		privateKey = pk
	}

	out, err := client.GetBalance(
		context.TODO(),
		privateKey.PublicKey(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("get balance failed %v\n", err))
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	t.SetCaption("Account SOL balance")

	t.AppendHeader(table.Row{"PublicKey", "SOL"})
	lamportsOnAccount := new(big.Float).SetUint64(uint64(out.Value))
	t.AppendRow(table.Row{privateKey.PublicKey(), new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))})

	fmt.Println(t.Render())
	return nil
}

// Get balance from by using private key or keygen file
func fromPKOrKGFPrompt() string {
	const (
		pk  = "private key"
		kgf = "keygen file"
	)

	prompt := promptui.Select{
		Label: "Please select option to get balance by",
		Items: []string{kgf, pk},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	if result == kgf {
		return "KGF"
	}

	return "PK"
}
