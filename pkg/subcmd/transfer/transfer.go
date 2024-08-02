package transfer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	"github.com/spf13/cobra"
)

type account struct {
	PublicKey string `json:"publicKey"`
	Amount    string `json:"amount"`
}

const (
	privateKeyQuestion      = "Please enter your private key"
	transferAccountQuestion = "Please enter account address to transfer"
	transferAmountQuestion  = "Please enter transfer amount"
)

var privateKey solana.PrivateKey

func TransferBalance(cmd *cobra.Command, args []string) error {
	endpoint := prompt.SelectNetworkPrompt()
	rpcClient := rpc.New(endpoint)
	// Transfer balance to list of accounts
	pathJsonFile, _ := cmd.Flags().GetString("path")

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

	// Get balance from private key
	out, err := rpcClient.GetBalance(
		context.TODO(),
		privateKey.PublicKey(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to get balance %v\n", err))
	}

	if strings.TrimSpace(pathJsonFile) != "" {
		// JSON file struct
		// [
		// {
		// "publicKey": "some_address",
		// "amount": "0.5"
		// },
		// {
		// "publicKey": "some_address",
		// "amount": "0.5"
		// }
		// ]
		// example --path="/homeDIR/subDIR/file.json"
		info, err := os.Stat(pathJsonFile)

		if os.IsNotExist(err) {
			return fmt.Errorf(fmt.Sprintf("failed to get file %v\n", err))
		}

		if info.IsDir() {
			return fmt.Errorf(fmt.Sprintf("failed to get file %v\n", pathJsonFile))
		}

		if err = checkFileFormat(pathJsonFile); err != nil {
			return err
		}

		file, err := os.Open(pathJsonFile)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("failed to open file: %v", err))
		}
		defer file.Close()

		// Read the file content
		byteValue, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("error reading file: %v", err))
		}

		var listAccounts []account

		err = json.Unmarshal(byteValue, &listAccounts)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("error decoding JSON %v:", err))
		}

		// validation public key list
		validateAccountPubKeyList(listAccounts)
		// validation from transfer amount
		validateFromTransferAmount(out, listAccounts)

		instructions, err := getTransferInstructions(privateKey, listAccounts)
		if err != nil {
			return err
		}

		if err = singAndSendTransaction(rpcClient, wsClient, privateKey, instructions); err != nil {
			return err
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleColoredBlackOnGreenWhite)
		t.SetCaption("Transfer SOL")

		t.AppendHeader(table.Row{"From", "To", "Amount"})

		for _, acc := range listAccounts {
			t.AppendRow(table.Row{privateKey.PublicKey(), acc.PublicKey, acc.Amount})
		}

		fmt.Println(t.Render())
		return nil
	}
	solAmountPrompt := promptui.Prompt{
		Label: transferAmountQuestion + constant.QUESTION_PROMPT_EXIT_PART,
	}

	amount, err := solAmountPrompt.Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("prompt failed %v\n", err))
	}

	// Check current account balance for transfer
	validateFromTransferAmount(out, []account{{PublicKey: "", Amount: amount}})

	pubKeyAccountPrompt := promptui.Prompt{
		Label: transferAccountQuestion + constant.QUESTION_PROMPT_EXIT_PART,
	}

	publicKeyAccount, err := pubKeyAccountPrompt.Run()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("prompt failed %v\n", err))
	}

	instructions, err := getTransferInstructions(privateKey, []account{{PublicKey: publicKeyAccount, Amount: amount}})
	if err != nil {
		return err
	}

	if err = singAndSendTransaction(rpcClient, wsClient, privateKey, instructions); err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	t.SetCaption("Transfer SOL")

	t.AppendHeader(table.Row{"From", "To", "Amount"})
	t.AppendRow(table.Row{privateKey.PublicKey().String(), publicKeyAccount, amount})

	fmt.Println(t.Render())
	return nil
}

func getWSRpc(clientRPCNet string) string {
	if clientRPCNet == rpc.MainNetBeta_RPC {
		return rpc.MainNetBeta_WS
	}

	return rpc.DevNet_WS
}

func checkFileFormat(file string) error {
	if filepath.Ext(file) != ".json" {
		return errors.New("invalid file format")
	}
	return nil
}

func validateAccountPubKeyList(list []account) {
	for _, acc := range list {
		_, err := solana.PublicKeyFromBase58(acc.PublicKey)

		if err != nil {
			log.Fatalf("failed to parse public key: %v, err: %v\n", acc.PublicKey, err)
		}
	}
}

// lamport value
func validateFromTransferAmount(balanceFrom *rpc.GetBalanceResult, accountListTo []account) {
	var sumToTransfer uint64
	// TODO fee check

	for _, acc := range accountListTo {
		lamportTransferAmount, err := helper.ConvertSolToLamports(acc.Amount)

		if err != nil {
			log.Fatalf("failed to convert sol to lamport %v\n", err)
		}

		sumToTransfer += lamportTransferAmount
	}

	if balanceFrom.Value < sumToTransfer || balanceFrom.Value <= 0 {
		log.Fatalf("your account balance is not enough for transfer SOL")
	}
}

func getTransferInstructions(fromPK solana.PrivateKey, accountListTo []account) ([]solana.Instruction, error) {
	res := []solana.Instruction{}

	if len(accountListTo) == 0 {
		return res, nil
	}

	for _, acc := range accountListTo {
		lamportTransferAmount, err := helper.ConvertSolToLamports(acc.Amount)
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("failed to convert sol to lamport %v\n", err))
		}

		accountTo, err := solana.PublicKeyFromBase58(acc.PublicKey)
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("failed to parse public key %v\n", acc.PublicKey))
		}

		res = append(res, system.NewTransferInstruction(
			lamportTransferAmount,
			fromPK.PublicKey(),
			accountTo,
		).Build())
	}

	return res, nil
}

func singAndSendTransaction(
	rpcClient *rpc.Client,
	wsClient *ws.Client,
	privateKey solana.PrivateKey,
	instructions []solana.Instruction,
) error {
	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to get recent blockhash %v\n", err))
	}

	tx, err := solana.NewTransaction(
		instructions,
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
		return fmt.Errorf(fmt.Sprintf("unable to sign transaction %v\n", err))
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
		return err
	}

	spew.Dump(sig)
	return nil
}
