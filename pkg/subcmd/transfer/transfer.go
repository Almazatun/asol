package transfer

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Almazatun/asol/constant"
	"github.com/Almazatun/asol/helper"
	"github.com/Almazatun/asol/pkg/prompt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/gagliardetto/solana-go/text"
	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
)

const (
	privateKeyQuestion      = "Please enter your private key"
	transferAccountQuestion = "Please enter account address to transfer"
	transferAmountQuestion  = "Please enter transfer amount"
)

var privateKey solana.PrivateKey

func TransferBalance(args []string) error {
	endpoint := prompt.SelectNetworkPrompt()
	rpcClient := rpc.New(endpoint)

	// Create a new WS client (used for confirming transactions)
	wsClient, err := ws.Connect(context.Background(), getWSRpc(endpoint))
	if err != nil {
		return err
	}

	selectOptToGetBalance := prompt.FromPKOrKGFPrompt()

	// TODO refactor
	if selectOptToGetBalance == "KGF" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("failed to get user's home directory: %v", err))
		}

		subDir := prompt.SubDirPrompt(homeDir)
		if err := helper.CheckSubDirExists(homeDir, subDir); err != nil {
			return err
		}

		keygenPath := filepath.Join(homeDir, subDir)
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

	out, err := rpcClient.GetBalance(
		context.TODO(),
		privateKey.PublicKey(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to get balance %v\n", err))
	}

	solAmountPrompt := promptui.Prompt{
		Label: transferAmountQuestion + constant.QUESTION_PROMPT_EXIT_PART,
	}

	amount, err := solAmountPrompt.Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("prompt failed %v\n", err))
	}

	lamportTransferAmount, err := helper.ConvertSolToLamports(amount)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to convert sol to lamport %v\n", err))
	}

	// Check current account balance for transfer
	if out.Value < lamportTransferAmount || lamportTransferAmount <= 0 {
		return fmt.Errorf("your account balance is not enough for transfer SOL")
	}

	pubKeyAccountPrompt := promptui.Prompt{
		Label: transferAccountQuestion + constant.QUESTION_PROMPT_EXIT_PART,
	}

	publicKeyAccount, err := pubKeyAccountPrompt.Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("prompt failed %v\n", err))
	}

	accountTo, err := solana.PublicKeyFromBase58(publicKeyAccount)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to parse public key %v\n", err))
	}

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to get recent blockhash %v\n", err))
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				lamportTransferAmount,
				privateKey.PublicKey(),
				accountTo,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(privateKey.PublicKey()),
	)

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to create transaction %v\n", err))
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if privateKey.PublicKey().Equals(key) {
				return &privateKey
			}
			return nil
		},
	)
	if err != nil {
		log.Fatalf("Unable to sign transaction %v\n", err)
	}

	log.Printf("🚀 Sending transaction...\n")
	spew.Dump(tx)
	// Pretty print the transaction:
	tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))

	sig, err := confirm.SendAndConfirmTransaction(
		context.TODO(),
		rpcClient,
		wsClient,
		tx,
	)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(sig)

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	t.SetCaption("Transfer SOL")

	t.AppendHeader(table.Row{"FromAddress", "ToAddress", "Amount"})
	t.AppendRow(table.Row{privateKey.PublicKey().String(), accountTo.String(), amount})

	fmt.Println(t.Render())
	return nil
}

func getWSRpc(clientRPCNet string) string {
	if clientRPCNet == rpc.MainNetBeta_RPC {
		return rpc.MainNetBeta_WS
	}

	return rpc.DevNet_WS
}
