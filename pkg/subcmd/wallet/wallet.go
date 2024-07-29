package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Almazatun/asol/helper"
	"github.com/Almazatun/asol/pkg/prompt"
	"github.com/gagliardetto/solana-go"
	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	errMsg               = "Please specify correct argument by wallet command"
	keyWord              = "new"
	defaultFileName      = "data.json"
	createJsonFilePrompt = "Create JSON file to save created wallet ?"
)

type walletData struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

const maxAccountToCreate = 100

var subDir string
var homeDir string
var enterFileName string

func CreateWallet(cmd *cobra.Command, args []string) error {
	listCountToCreateStr, _ := cmd.Flags().GetString("list")

	// Create a list of account
	if listCountToCreateStr != "" {
		listCountToCreate, err := strconv.Atoi(listCountToCreateStr)

		if err != nil || listCountToCreate <= 0 {
			return errors.New("create wallet count should be a Int without any symbols and more than 0")
		}

		if err = checkMaxAccountVal(listCountToCreate); err != nil {
			return err
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleColoredBlackOnGreenWhite)
		t.AppendHeader(table.Row{"#", "PublicKey", "PrivateKey"})
		t.SetCaption("Wallets")

		createJsonPromptRes := prompt.YesOrNoPromptByQuestion(createJsonFilePrompt)

		// Check subdirectory
		if createJsonPromptRes == "Yes" {
			// Get the user's home directory
			hDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf(fmt.Sprintf("failed to get user's home directory: %v", err))
			}

			homeDir = hDir

			sDir := prompt.SubDirPrompt(hDir)
			if err := helper.CheckSubDirExists(homeDir, subDir); err != nil {
				return err
			}

			subDir = sDir
		}

		list := []walletData{}
		for i := 0; i < listCountToCreate; i++ {
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

		// Create JSON file
		if createJsonPromptRes == "Yes" {
			enterFileName = fileNamePrompt()
			// Create the full path
			fullPath := filepath.Join(homeDir, subDir, enterFileName)

			if err := createJsonAccountFile(fullPath, list); err != nil {
				return err
			}
		}

		fmt.Println(t.Render())
		return nil
	}

	// Create single account
	acc := solana.NewWallet()
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	t.SetCaption("Wallet")

	t.AppendHeader(table.Row{"PublicKey", "PrivateKey"})
	t.AppendRow(table.Row{acc.PublicKey(), acc.PrivateKey})
	createJsonPromptRes := prompt.YesOrNoPromptByQuestion(createJsonFilePrompt)

	// Check subdirectory
	if createJsonPromptRes == "Yes" {
		// Get the user's home directory
		hDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("failed to get user's home directory: %v", err))
		}

		homeDir = hDir

		sDir := prompt.SubDirPrompt(hDir)
		if err := helper.CheckSubDirExists(homeDir, subDir); err != nil {
			return err
		}

		subDir = sDir
	}

	if createJsonPromptRes == "Yes" {
		walletData := walletData{
			PrivateKey: acc.PrivateKey.String(),
			PublicKey:  acc.PublicKey().String(),
		}

		enterFileName = fileNamePrompt()
		// Create the full path
		fullPath := filepath.Join(homeDir, subDir, enterFileName)

		if err := createJsonAccountFile(fullPath, walletData); err != nil {
			return err
		}

	}

	fmt.Println(t.Render())
	return nil
}

func createJsonAccountFile(fullPath string, data interface{}) error {
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to open file: %v", err))
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("failed to encode JSON: %v", err))
	}

	return nil
}

func checkMaxAccountVal(i int) error {
	if math.Floor(float64(i)) > float64(maxAccountToCreate) {
		return fmt.Errorf(fmt.Sprintf("maximum number to create wallets should be less than or qual to %v", maxAccountToCreate))
	}

	return nil
}

func fileNamePrompt() string {
	prompt := promptui.Prompt{
		Label: "Please enter file name",
	}

	result, err := prompt.Run()

	if err != nil {
		log.Fatalf("prompt failed %v\n", err)
	}

	if strings.Contains(result, ".json") {
		log.Fatalf("prompt failed: please enter only file name without {.format}\n")
	}

	if result == "" {
		return defaultFileName
	}

	return result + ".json"
}
