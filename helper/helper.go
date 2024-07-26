package helper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gagliardetto/solana-go"
)

const lamportsSol = 1000_000_000

func ConvertSolToLamports(sol string) (uint64, error) {
	// Parse the SOL amount from string to float64
	solValue, err := strconv.ParseFloat(sol, 64)
	if err != nil {
		return 0, err
	}

	// Convert SOL to lamports (1 SOL = 1,000,000,000 lamports)
	lamports := uint64(solValue * lamportsSol)
	return lamports, nil
}

// ConvertLamportsToSol converts lamports to SOL.
func ConvertLamportsToSol(lamports uint64) string {
	// Convert lamports to SOL (1,000,000,000 lamports = 1 SOL)
	solValue := float64(lamports) / lamportsSol
	return fmt.Sprintf("%.9f", solValue)
}

func ValidatePrivateKey(input string) error {
	pk, err := solana.PrivateKeyFromBase58(input)

	// Only way to check is correct private key
	// But this method call panic instead of return error if invalid private key
	// pk.PublicKey()
	if err != nil || len(pk.String()) == 0 {
		return errors.New("Invalid private key")
	}

	return nil
}
