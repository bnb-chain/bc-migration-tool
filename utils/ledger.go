package utils

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/usbwallet"
	"strings"
)

var (
	ledgerBasePath = accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0}
)

func OpenLedgerAccount(index uint32) (accounts.Wallet, accounts.Account, error) {
	ledgerHub, err := usbwallet.NewLedgerHub()
	if err != nil {
		return nil, accounts.Account{}, err
	}

	if len(ledgerHub.Wallets()) == 0 {
		return nil, accounts.Account{}, errors.New("empty ledger wallet")
	}

	wallet := ledgerHub.Wallets()[0]
	if err := wallet.Open(""); err != nil {
		return nil, accounts.Account{}, err
	}

	if wallet == nil {
		return nil, accounts.Account{}, errors.New("ledger account not found")
	}

	status, err := wallet.Status()
	if err != nil {
		return wallet, accounts.Account{}, err
	}
	if strings.Contains(status, "offline") {
		return wallet, accounts.Account{}, errors.New("please open Ethereum app on ledger")
	}
	fmt.Println("Ledger status", status)

	ledgerPath := make(accounts.DerivationPath, len(ledgerBasePath))
	copy(ledgerPath, ledgerBasePath)
	ledgerPath[2] = ledgerPath[2] + index
	ledgerAccount, err := wallet.Derive(ledgerPath, true)
	if err != nil {
		return wallet, accounts.Account{}, err
	}
	return wallet, ledgerAccount, nil
}
