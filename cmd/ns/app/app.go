package app

import (
	"os"

	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/mossid/sdk-nameservice-example/x/faucet"
	"github.com/mossid/sdk-nameservice-example/x/nameservice"
)

const (
	appName = "NSApp"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.example_nscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.example_nsd")
)

type NSApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain    *sdk.KVStoreKey
	keyAccount *sdk.KVStoreKey
	keyNS      *sdk.KVStoreKey

	accountMapper auth.AccountMapper
	bankKeeper    bank.Keeper
	nsKeeper      nameservice.Keeper
}

func NewNSApp(logger log.Logger, db dbm.DB) *NSApp {
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	var app = &NSApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:    sdk.NewKVStoreKey("main"),
		keyAccount: sdk.NewKVStoreKey("acc"),
		keyNS:      sdk.NewKVStoreKey("ns"),
	}

	app.accountMapper = auth.NewAccountMapper(
		app.cdc,
		app.keyAccount,
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(app.accountMapper)

	app.nsKeeper = nameservice.NewKeeper(
		app.cdc,
		app.keyNS,
		app.bankKeeper,
	)

	app.Router().
		AddRoute("nameservice", nameservice.NewHandler(app.nsKeeper)).
		AddRoute("faucet", faucet.NewHandler(app.bankKeeper))

	app.MountStoresIAVL(
		app.keyMain,
		app.keyAccount,
		app.keyNS,
	)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	nameservice.RegisterCodec(cdc)
	faucet.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
