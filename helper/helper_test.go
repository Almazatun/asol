package helper

import (
	"testing"
)

const (
	solVal = "1"
	lamVal = uint64(1000000000)
)

func TestConvertSolToLamports(t *testing.T) {
	lamports, err := ConvertSolToLamports(solVal)
	if err != nil {
		t.Errorf("Error converting SOL to lamports: %v", err)
	}

	if lamports != lamVal {
		t.Fatalf("Expected %d lamports but got %d", lamVal, lamports)
	}
}

func TestLamportsToSol(t *testing.T) {
	sol := ConvertLamportsToSol(lamVal)

	if sol != "1.000000000" {
		t.Errorf("Expected 1.000000000 SOL but got %s", sol)
	}
}
