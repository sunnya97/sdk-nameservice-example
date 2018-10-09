package app

import (
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	faucet "github.com/sunnya97/sdk-faucet-module"
	"github.com/sunnya97/sdk-nameservice-example/x/nameservice"
)

const (
	appName = "Nameshake"
)

type NameshakeApp struct {
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

func NewNameshakeApp(logger log.Logger, db dbm.DB) *NameshakeApp {
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	var app = &NameshakeApp{
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

type GenesisState struct {
	Accounts []auth.BaseAccount `json:"accounts"`
}

func (app *NameshakeApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
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
