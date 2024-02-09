package main

import (
	"fmt"
	"os"

	"github.com/reecepbcups/tokenfactory/app"
	"github.com/reecepbcups/tokenfactory/cmd/tokend/cmd"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "tokend", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
