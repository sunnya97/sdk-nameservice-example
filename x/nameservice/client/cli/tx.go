package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"

	"github.com/sunnya97/sdk-nameservice-example/x/nameservice"
)

const (
	flagName   = "name"
	flagValue  = "value"
	flagAmount = "amount"
)

func GetCmdBuyName(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy-name [name] [amount]",
		Short: "bid for existing name or claim new name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			name := args[0]

			amount := args[1]
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			account, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := nameservice.MsgBuyName{
				NameID: name,
				Bid:    coins,
				Buyer:  account,
			}

			tx := auth.StdTx{
				Msgs: []sdk.Msg{msg},
			}

			bz := cdc.MustMarshalBinary(tx)

			_, err = cliCtx.BroadcastTx(bz)

			return err
		},
	}

	return cmd
}

func GetCmdSetName(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-name [name] [value]",
		Short: "set the value associated with a name that you own",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			name := args[0]
			value := args[1]

			account, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := nameservice.MsgSetName{
				NameID: name,
				Value:  value,
				Owner:  account,
			}

			tx := auth.StdTx{
				Msgs: []sdk.Msg{msg},
			}

			bz := cdc.MustMarshalBinary(tx)

			_, err = cliCtx.BroadcastTx(bz)

			return err
		},
	}

	return cmd
}
