package main

import (
	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"

	"github.com/mossid/sdk-nameservice-example/cmd/ns/app"
	faucetcmd "github.com/mossid/sdk-nameservice-example/x/faucet/client/cli"
	nscmd "github.com/mossid/sdk-nameservice-example/x/nameservice/client/cli"
)

const storeAcc = "acc"
const storeNS = "ns"

var (
	rootCmd = &cobra.Command{
		Use:   "nscli",
		Short: "Namespace App Client",
	}
)

func main() {
	cobra.EnableCommandSorting = false
	cdc := app.MakeCodec()

	rootCmd.AddCommand(client.ConfigCmd())
	rpc.AddCommands(rootCmd)

	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(),
	)
	tx.AddCommands(queryCmd, cdc)
	queryCmd.AddCommand(client.LineBreak)
	queryCmd.AddCommand(client.GetCommands(
		authcmd.GetAccountCmd(storeAcc, cdc, authcmd.GetAccountDecoder(cdc)),
		nscmd.GetCmdQueryDomain(storeNS, cdc),
	)...)

	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		nscmd.GetCmdBuyDomain(cdc),
		faucetcmd.GetCmdRequestCoins(cdc),
	)...)

	rootCmd.AddCommand(
		queryCmd,
		txCmd,
		client.LineBreak,
	)

	rootCmd.AddCommand(
		keys.Commands(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "NS", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
