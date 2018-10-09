package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
)

func completeAndBroadcastTxCli(txBldr authtxb.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) error {

	name, err := cliCtx.GetFromName()
	if err != nil {
		return err
	}

	passphrase, err := keys.GetPassphrase(name)
	if err != nil {
		return err
	}

	txBytes, err := txBldr.BuildAndSign(name, passphrase, msgs)
	if err != nil {
		return err
	}

	_, err = cliCtx.BroadcastTx(txBytes)
	return err
}
