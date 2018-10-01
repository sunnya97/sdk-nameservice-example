package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/sdk-nameservice-example/x/nameservice"
)

type QueryResult struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// GetCmdQueryDomain queries information about a domain
func GetCmdQueryDomain(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain [domain]",
		Short: "Query domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			value, err := cliCtx.QueryStore(nameservice.ValueKey(domain), storeName)
			if err != nil {
				return err
			}

			owner, err := cliCtx.QueryStore(nameservice.OwnerKey(domain), storeName)
			if err != nil {
				return err
			}

			pricebz, err := cliCtx.QueryStore(nameservice.PriceKey(domain), storeName)
			if err != nil {
				return err
			}
			var price sdk.Coins
			cdc.MustUnmarshalBinary(pricebz, &price)

			result := QueryResult{
				Value: string(value),
				Owner: owner,
				Price: price,
			}

			output, err := codec.MarshalJSONIndent(cdc, result)
			if err != nil {
				return err
			}

			fmt.Println(string(output))

			return nil
		},
	}

	return cmd
}
