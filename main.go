package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bnb-chain/bc-migration-tool/cmd"
)

func main() {
	rootCmd := &cobra.Command{Use: "bc-migration-tool"}

	cmd.AddCreateCmd(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
