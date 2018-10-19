# Entrypoints

In golang the convention is to place files that compile to a binary in the `./cmd` folder of a project. For our application we have 2 binaries that we want to create:

- `nameserviced`: This binary is similar to `bitcoind` or other cryptocurrency daemons in that it maintains peer connections a propagates transactions through the network.
- `nameservicecli`: This binary provides commands that allow users to interact with your application.


To get started create two files in the root of the project directory that will instantiate these binaries:
- `./cmd/nameserviced/main.go`
- `./cmd/nameservicecli/main.go`

## `nameserviced`

Start by adding the following code to `nameserviced/main.go`:

> _*NOTE*_: Your application needs to import the code you just wrote. Here the import path is set to this repository (`github.com/jackzampolin/sdk-nameservice-example`). If you are following along in your own repo you will need to change the import path to reflect that (`github.com/{{ .Username }}/{{ .Project.Repo }}`).

```go
package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	app "github.com/jackzampolin/sdk-nameservice-example"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

// DefaultNodeHome sets the folder where the application data and configuration will be stored
var DefaultNodeHome = os.ExpandEnv("$HOME/.nameserviced")

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	appInit := server.AppInit{
		AppGenState: server.SimpleAppGenState,
		AppGenTx:    server.SimpleAppGenTx,
	}

	rootCmd := &cobra.Command{
		Use:               "nameserviced",
		Short:             "nameservice App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, appInit,
		server.ConstructAppCreator(newApp, "nameservice"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "nameservice"))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "NS", DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewnameserviceApp(logger, db)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	return nil, nil, nil
}
```

Notes on the above code:
- Most of the code above combines the CLI commands from 1. Tendermint, 2. CosmosSDK, and 3. Nameservice module
- The rest of the code helps the application generate genesis state from the configuration

## `nameservicecli`

Finish up by building the `nameservicecli` command:

> _*NOTE*_: Your application needs to import the code you just wrote. Here the import path is set to this repository (`github.com/jackzampolin/sdk-nameservice-example`). If you are following along in your own repo you will need to change the import path to reflect that (`github.com/{{ .Username }}/{{ .Project.Repo }}`).

```go
package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"

	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	app "github.com/jackzampolin/sdk-nameservice-example"
	nameservicecmd "github.com/jackzampolin/sdk-nameservice-example/x/nameservice/client/cli"
	faucetcmd "github.com/sunnya97/sdk-faucet-module/client/cli"
)

const storeAcc = "acc"

var (
	rootCmd = &cobra.Command{
		Use:   "nameservicecli",
		Short: "nameservice Client",
	}
	DefaultCLIHome = os.ExpandEnv("$HOME/.nameservicecli")
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
		nameservicecmd.GetCmdResolveName("nameservice", cdc),
		nameservicecmd.GetCmdWhois("nameservice", cdc),
	)...)

	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(client.PostCommands(
		nameservicecmd.GetCmdBuyName(cdc),
		nameservicecmd.GetCmdSetName(cdc),
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

	executor := cli.PrepareMainCmd(rootCmd, "NS", DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
```

Notes on the above code:
- Most of the code above combines the CLI commands from 1. Tendermint, 2. CosmosSDK, and 3. Nameservice module

### Now that you have your binaries defined its time to deal with [dependency management and build your app](./dep.md)!
