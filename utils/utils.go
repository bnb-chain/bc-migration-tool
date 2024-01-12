package utils

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SignTxByPrivateKey(client *ethclient.Client, privateKeyHex, opAccount string, data []byte, value *big.Int) (*types.Transaction, error) {
	if privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}
	privKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	pubKey := privKey.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(*pubKeyECDSA)
	if addr.String() != opAccount {
		return nil, fmt.Errorf("private key does not match the operator account")
	}

	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &StakeHubAddress,
		Value:    value,
		Gas:      DefaultGasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	signer := types.NewEIP155Signer(chainId)
	signedTx, err := types.SignTx(tx, signer, privKey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func SignTxByKeystore(client *ethclient.Client, keystorePath, password, opAccount string, data []byte, value *big.Int) (*types.Transaction, error) {
	ks := keystore.NewKeyStore(keystorePath, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.Find(accounts.Account{Address: common.HexToAddress(opAccount)})
	if err != nil {
		return nil, err
	}
	err = ks.Unlock(acc, password)
	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), acc.Address)
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &StakeHubAddress,
		Value:    value,
		Gas:      DefaultGasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	signedTx, err := ks.SignTx(acc, tx, chainId)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func SignTxByLedger(client *ethclient.Client, opAccount string, index uint32, data []byte, value *big.Int) (*types.Transaction, error) {
	wallet, acc, err := OpenLedgerAccount(index)
	if err != nil {
		return nil, err
	}
	defer wallet.Close()

	if acc.Address.Hex() != opAccount {
		return nil, errors.New("account does match, please set the correct account index")
	}

	nonce, err := client.PendingNonceAt(context.Background(), acc.Address)
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &StakeHubAddress,
		Value:    value,
		Gas:      DefaultGasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	signedTx, err := wallet.SignTx(acc, tx, chainId)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
