package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bnb-chain/bc-migration-tool/utils"
)

func AddVerifyLedgerCmd(rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "verify-ledger",
		Short: "Verify Ledger account index",
		RunE: func(cmd *cobra.Command, _ []string) error {
			opAccount, err := cmd.Flags().GetString(FlagOperatorAccount)
			if err != nil {
				return err
			}

			usingLedger, err := cmd.Flags().GetBool(FlagLedger)
			if err != nil {
				return err
			}
			if !usingLedger {
				return errors.New("only for ledger")
			}

			index, err := cmd.Flags().GetUint32(FlagIndex)
			if err != nil {
				return err
			}

			wallet, ledgerAccount, err := utils.OpenLedgerAccount(index)
			if err != nil {
				return err
			}
			defer wallet.Close()

			fmt.Println("Connected Ledger Account", ledgerAccount.Address)
			if ledgerAccount.Address.Hex() != opAccount {
				return errors.New("account does not match")
			}

			fmt.Println("the account index matches operator address")

			return nil
		},
	}

	cmd.Flags().String(FlagOperatorAccount, "", "operator account address")
	cmd.Flags().Bool(FlagLedger, false, "whether the operator account is a ledger account")
	cmd.Flags().Uint32(FlagIndex, 0, "ledger account index")

	rootCmd.AddCommand(cmd)
}
