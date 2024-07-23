package wallet

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/Almazatun/asol/pkg/prompt"
	"github.com/gagliardetto/solana-go"
	"github.com/jedib0t/go-pretty/table"
)

const (
	errMsg               = "Please specify correct argument by wallet command"
	keyWord              = "new"
	fileName             = "data.json"
	createJsonFilePrompt = "Create JSON file to save created wallet ?"
)

type walletData struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

func CreateWallet(args []string) {
	if len(args) == 1 {
		if args[0] != keyWord {
			log.Fatal(errMsg)
		}

		acc := solana.NewWallet()
		t := table.NewWriter()
		t.SetStyle(table.StyleColoredBlackOnGreenWhite)
		t.SetCaption("Wallet")

		t.AppendHeader(table.Row{"PublicKey", "PrivateKey"})
		t.AppendRow(table.Row{acc.PublicKey(), acc.PrivateKey})
		createJsonPromptRes := prompt.YesOrNoPromptByQuestion(createJsonFilePrompt)

		if createJsonPromptRes == "Yes" {
			walletData := walletData{
				PrivateKey: acc.PrivateKey.String(),
				PublicKey:  acc.PublicKey().String(),
			}

			file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				log.Fatalf("Failed to open file: %v", err)
			}

			defer file.Close()

			encoder := json.NewEncoder(file)

			err = encoder.Encode(walletData)
			if err != nil {
				log.Fatalf("Failed to encode JSON: %v", err)
			}
		}

		fmt.Println(t.Render())

		return
	}

	maxCountWalletsToCreate := 100
	if len(args) == 2 {
		if args[0] != keyWord {
			log.Fatal(errMsg)
		}

		num, err := strconv.Atoi(args[1])
		if err != nil {
			log.Println("Create wallet count should be a number without any symbols")

			return
		}

		if math.Floor(float64(num)) > float64(maxCountWalletsToCreate) {
			log.Fatalf("Maximum number to create wallets should be less than or qual to %v", maxCountWalletsToCreate)

			return
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleColoredBlackOnGreenWhite)
		t.AppendHeader(table.Row{"#", "PublicKey", "PrivateKey"})
		t.SetCaption("Wallets")

		createJsonPromptRes := prompt.YesOrNoPromptByQuestion(createJsonFilePrompt)

		list := []walletData{}
		for i := 0; i < num; i++ {
			acc := solana.NewWallet()
			t.AppendRow(table.Row{i + 1, acc.PublicKey(), acc.PrivateKey})

			if createJsonPromptRes == "Yes" {
				walletData := walletData{
					PrivateKey: acc.PrivateKey.String(),
					PublicKey:  acc.PublicKey().String(),
				}

				list = append(list, walletData)
			}
		}

		if createJsonPromptRes == "Yes" {
			file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				log.Fatalf("Failed to open file: %v", err)
			}

			defer file.Close()

			encoder := json.NewEncoder(file)

			err = encoder.Encode(list)
			if err != nil {
				log.Fatalf("Failed to encode JSON: %v", err)
			}
		}

		fmt.Println(t.Render())
		return
	}

	log.Println(errMsg)
}
