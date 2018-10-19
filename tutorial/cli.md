## Nameservice Module CLI

Next, in the module's folder, create two files:
 - `./client/cli/query.go`
 - `./client/cli/tx.go`

These will enable our cli to understand our module.

query.go
```go
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryResult struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// GetCmdResolveName queries information about a name
func GetCmdResolveName(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve [name]",
		Short: "resolve name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/proposals/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not resolve name - %s \n", string(name))
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}

// GetCmdWhois queries information about a domain
func GetCmdWhois(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whois [name]",
		Short: "Query whois info of name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not resolve whois - %s \n", string(name))
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}

	return cmd
}
```

tx.go
```go
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

```
