## App.go

To create your application start a new file: `./app.go`. To get started add the dependencies you will need:

> _*NOTE*_: Your application needs to import the code you just wrote. Here the import path is set to this repository (`github.com/jackzampolin/sdk-nameservice-example/x/nameservice`). If you are following along in your own repo you will need to change the import path to reflect that (`github.com/{{ .Username }}/{{ .Project.Repo }}/x/nameservice`).

```go
package app

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/jackzampolin/sdk-nameservice-example/x/nameservice"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	faucet "github.com/sunnya97/sdk-faucet-module"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
)
```

Links to godocs for each module:
- [`codec`](https://godoc.org/github.com/cosmos/cosmos-sdk/codec): Functions for working with amino
- [`auth`](https://godoc.org/github.com/cosmos/cosmos-sdk/x/auth)
- [`bank`](https://godoc.org/github.com/cosmos/cosmos-sdk/x/bank)
- [`baseapp`](https://godoc.org/github.com/cosmos/cosmos-sdk): This module helps developers bootstrap CosmosSDK applications
- [`types`](https://godoc.org/github.com/cosmos/cosmos-sdk): Common types for working with SDK applications
- [`abci`](https://godoc.org/github.com/tendermint/tendermint/abci/types): Similar to the `sdk/types` module, but for Tendermint
- [`cmn`](https://godoc.org/github.com/tendermint/tendermint/libs/common): Code for working with Tendermint applications
- [`dbm`](https://godoc.org/github.com/tendermint/tendermint/libs/db): Code for working with the Tendermint database

Start by declaring the name and struct for our app.  In this tutorial the app is called `nameservice`.

```go
const (
	appName = "nameservice"
)

type nameserviceApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain     *sdk.KVStoreKey
	keyAccount  *sdk.KVStoreKey
	keyNSnames  *sdk.KVStoreKey
	keyNSowners *sdk.KVStoreKey
	keyNSprices *sdk.KVStoreKey

	accountMapper auth.AccountMapper
	bankKeeper    bank.Keeper
	nsKeeper      nameservice.Keeper
}
```

Next, create a constructor for a new `nameserviceApp`.  In this function your application:
- Generates `storeKeys`
- Creates `Keepers`
- Registers `Handler`s
- Registers `Querier`s
- Mounts `KVStore`s
- Sets the `initChainer`

```go
func NewnameserviceApp(logger log.Logger, db dbm.DB) *nameserviceApp {
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	var app = &nameserviceApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:     sdk.NewKVStoreKey("main"),
		keyAccount:  sdk.NewKVStoreKey("acc"),
		keyNSnames:  sdk.NewKVStoreKey("ns_names"),
		keyNSowners: sdk.NewKVStoreKey("ns_owners"),
		keyNSprices: sdk.NewKVStoreKey("ns_prices"),
	}

	app.accountMapper = auth.NewAccountMapper(
		app.cdc,
		app.keyAccount,
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(app.accountMapper)

	app.nsKeeper = nameservice.NewKeeper(
		app.bankKeeper,
		app.keyNSnames,
		app.keyNSowners,
		app.keyNSprices,
		app.cdc,
	)

	app.Router().
		AddRoute("nameservice", nameservice.NewHandler(app.nsKeeper)).
		AddRoute("faucet", faucet.NewHandler(app.bankKeeper))

	app.QueryRouter().
		AddRoute("nameservice", nameservice.NewQuerier(app.nsKeeper))

	app.SetInitChainer(app.initChainer)

	app.MountStoresIAVL(
		app.keyMain,
		app.keyAccount,
		app.keyNSnames,
		app.keyNSowners,
		app.keyNSprices,
	)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}
```

The `initChainer` defines how accounts in `genesis.json` are mapped into the application state on initial chain start. The constructor registers the `initChainer` function, but it isn't defined yet. Go ahead and create it:

```go
// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState struct {
	Accounts []auth.BaseAccount `json:"accounts"`
}

func (app *nameserviceApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, &acc)
	}

	return abci.ResponseInitChain{}
}

```

Finally add a helper function to generate an amino [`*codec.Codec`](https://godoc.org/github.com/cosmos/cosmos-sdk/codec#Codec) that properly registers all of the modules used in your application:

```go
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	nameservice.RegisterCodec(cdc)
	faucet.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
```

### Now that you have created an application that includes your module, it's time to [build your entrypoints](./entrypoint.md)!
