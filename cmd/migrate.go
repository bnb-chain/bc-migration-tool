package cmd

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	validatorpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/iface"
	"github.com/prysmaticlabs/prysm/v4/validator/accounts/wallet"
	"github.com/prysmaticlabs/prysm/v4/validator/keymanager"
	"github.com/spf13/cobra"

	"github.com/bnb-chain/bc-migration-tool/abi"
	"github.com/bnb-chain/bc-migration-tool/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func AddCreateCmd(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "Create a validator on BSC",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := utils.NewConfig()
			if err != nil {
				return err
			}

			blsPassword, err := cmd.Flags().GetString(FlagBlsPassword)
			if err != nil {
				return err
			}

			blsPubkeyStr, err := cmd.Flags().GetString(FlagBlsPubkey)
			if err != nil {
				return err
			}
			if blsPubkeyStr[:2] == "0x" {
				blsPubkeyStr = blsPubkeyStr[2:]
			}
			blsPubkey, err := hex.DecodeString(blsPubkeyStr)
			if err != nil {
				return err
			}

			blsKm, err := getBlsKeymanager(cfg.BlsDataDir, blsPassword)
			if err != nil {
				return err
			}

			pubkeys, err := blsKm.FetchValidatingPublicKeys(context.Background())
			if err != nil {
				return err
			}
			var found bool
			for _, pubkey := range pubkeys {
				if bytes.Equal(pubkey[:], blsPubkey) {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("bls pubkey not found")
			}

			opAccount, err := cmd.Flags().GetString(FlagOperatorAccount)
			if err != nil {
				return err
			}
			usingLedger, err := cmd.Flags().GetBool(FlagLedger)
			if err != nil {
				return err
			}
			privateKey, err := cmd.Flags().GetString(FlagPrivateKey)
			if err != nil {
				return err
			}
			keystorePath, err := cmd.Flags().GetString(FlagKeystorePath)
			if err != nil {
				return err
			}
			if !usingLedger && privateKey == "" && keystorePath == "" {
				return fmt.Errorf("ledger or keystore or private key is required")
			}

			client, err := ethclient.Dial(cfg.BscRpcUrl)
			if err != nil {
				return err
			}

			consensusAddr := common.HexToAddress(cfg.ValidatorInfo.ConsensusAddress)
			delegation, ok := new(big.Int).SetString(cfg.ValidatorInfo.Delegation, 10)
			if !ok {
				return fmt.Errorf("invalid delegation amount")
			}
			description := abi.StakeHubDescription{
				Moniker:  cfg.ValidatorInfo.Description.Moniker,
				Identity: cfg.ValidatorInfo.Description.Identity,
				Website:  cfg.ValidatorInfo.Description.Website,
				Details:  cfg.ValidatorInfo.Description.Details,
			}
			commission := abi.StakeHubCommission{
				Rate:          cfg.ValidatorInfo.Commission.Rate,
				MaxRate:       cfg.ValidatorInfo.Commission.MaxRate,
				MaxChangeRate: cfg.ValidatorInfo.Commission.MaxChangeRate,
			}

			chainId, err := client.ChainID(context.Background())
			if err != nil {
				return err
			}
			paddedChainIdBytes := make([]byte, 32)
			copy(paddedChainIdBytes[32-len(chainId.Bytes()):], chainId.Bytes())
			msgHash := crypto.Keccak256(append(blsPubkey, paddedChainIdBytes...))
			req := validatorpb.SignRequest{
				PublicKey:   blsPubkey,
				SigningRoot: msgHash,
			}
			proof, err := blsKm.Sign(context.Background(), &req)

			stakeHubAbi, err := abi.StakeHubMetaData.GetAbi()
			if err != nil {
				return err
			}
			method := "createValidator"
			data, err := stakeHubAbi.Pack(method, consensusAddr, blsPubkey, proof.Marshal(), commission, description)
			if err != nil {
				return err
			}

			var signedTx *types.Transaction
			if usingLedger {
				signedTx, err = utils.SignTxByLedger(client, opAccount, data, delegation)
			} else if privateKey != "" {
				signedTx, err = utils.SignTxByPrivateKey(client, privateKey, opAccount, data, delegation)
			} else {
				password, err := cmd.Flags().GetString(FlagPassword)
				if err != nil {
					return err
				}
				signedTx, err = utils.SignTxByKeystore(client, keystorePath, password, opAccount, data, delegation)
			}
			if err != nil {
				return err
			}

			err = client.SendTransaction(context.Background(), signedTx)
			if err != nil {
				return err
			}

			var receipt *types.Receipt
			for {
				receipt, err = client.TransactionReceipt(context.Background(), signedTx.Hash())
				if err != nil {
					if err.Error() == "not found" {
						time.Sleep(3 * time.Second)
						continue
					}
					return err
				}
				if receipt != nil {
					break
				}
			}

			if receipt.Status == 0 {
				return fmt.Errorf("transaction failed. tx hash: 0x%x", signedTx.Hash())
			} else {
				fmt.Println("transaction success")
			}

			return nil
		},
	}

	cmd.Flags().String(FlagBlsPassword, "", "bls wallet password")
	cmd.Flags().String(FlagBlsPubkey, "", "bls pubkey")
	cmd.Flags().Bool(FlagLedger, false, "whether the operator account is a ledger account")
	cmd.Flags().String(FlagPrivateKey, "", "private key of operator account")
	cmd.Flags().String(FlagKeystorePath, "", "keystore path of operator account")
	cmd.Flags().String(FlagPassword, "", "password of the keystore")
	cmd.Flags().String(FlagOperatorAccount, "", "operator account address")

	rootCmd.AddCommand(cmd)
}

func getBlsKeymanager(walletPath, password string) (keymanager.IKeymanager, error) {
	w, err := wallet.OpenWallet(context.Background(), &wallet.Config{
		WalletDir:      walletPath,
		WalletPassword: password,
	})
	if err != nil {
		return nil, err
	}

	km, err := w.InitializeKeymanager(context.Background(), iface.InitKeymanagerConfig{ListenForChanges: false})
	if err != nil {
		return nil, err
	}

	return km, nil
}
