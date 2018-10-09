package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"

	"github.com/mossid/sdk-nameservice-example/x/nameservice"
)

const (
	flagDomain = "domain"
	flagValue  = "value"
	flagAmount = "amount"
)

func GetCmdBuyDomain(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy-domain",
		Short: "bid for existing domain or claim new domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			domain := viper.GetString(flagDomain)
			value := viper.GetString(flagValue)

			amount := viper.GetString(flagAmount)
			coins, err := sdk.ParseCoins(amount)
			if err != nil {
				return err
			}

			account, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := nameservice.MsgBuyDomain{
				Domain: domain,
				Value:  value,
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

	cmd.Flags().String(flagDomain, "", "Domain to claim")
	cmd.Flags().String(flagValue, "", "Value to associate with the domain")
	cmd.Flags().String(flagAmount, "", "Coins willing to pay for the domain")

	return cmd
}
