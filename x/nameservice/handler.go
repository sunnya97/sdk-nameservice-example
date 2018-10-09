package nameservice

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgBuyDomain:
			return handleMsgBuyDomain(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgBuyDomain
func handleMsgBuyDomain(ctx sdk.Context, keeper Keeper, msg MsgBuyDomain) sdk.Result {
	if keeper.GetPrice(ctx, msg.Domain).IsGTE(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}

	// Subtract coins from the buyer, send it to the previous owner if exists
	if keeper.HasOwner(ctx, msg.Domain) {
		_, err := keeper.bk.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Domain), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		_, _, err := keeper.bk.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}

	keeper.SetOwner(ctx, msg.Domain, msg.Buyer)
	keeper.SetPrice(ctx, msg.Domain, msg.Bid)
	keeper.SetValue(ctx, msg.Domain, msg.Value)

	return sdk.Result{}
}
