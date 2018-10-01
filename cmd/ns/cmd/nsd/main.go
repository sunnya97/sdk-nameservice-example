package main

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/server"

	"github.com/mossid/sdk-nameservice-example/cmd/ns/app"
)

var appInit = server.AppInit{
	AppGenState: server.SimpleAppGenState,
	AppGenTx:    server.SimpleAppGenTx,
}

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "nsd",
		Short:             "Nameservice App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, appInit,
		server.ConstructAppCreator(newApp, "ns"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "ns"))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "NS", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewNSApp(logger, db)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	return nil, nil, nil
}
