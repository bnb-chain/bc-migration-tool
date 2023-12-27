package utils

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/usbwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

func TestLedger(t *testing.T) {
	opAccount := "0xdEaD"
	ledgerHub, err := usbwallet.NewLedgerHub()
	if err != nil {
		t.Error(err)
	}

	var wallet accounts.Wallet
	var acc accounts.Account
	for _, w := range ledgerHub.Wallets() {
		if w.Contains(accounts.Account{Address: common.HexToAddress(opAccount)}) {
			wallet = w
			acc = accounts.Account{Address: common.HexToAddress(opAccount)}
			break
		}
	}
	if wallet == nil {
		t.Error("ledger account not found")
	}
	if err := wallet.Open(""); err != nil {
		t.Error(err)
	}
	defer wallet.Close()

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       &StakeHubAddress,
		Value:    nil,
		Gas:      DefaultGasLimit,
		GasPrice: big.NewInt(500000000),
	})

	chainId := big.NewInt(714)
	signedTx, err := wallet.SignTx(acc, tx, chainId)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(signedTx)
}
