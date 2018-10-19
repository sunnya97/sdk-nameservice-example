### Buy Name

Great, now owners can set names!  But what if a name doesn't have an owner yet?  We need a way for people to buy names!

#### Msg

We define the Msg for buying names and add it to the `msgs.go` file:

```go
type MsgBuyName struct {
	NameID  string
	Bid   sdk.Coins
	Buyer sdk.AccAddress
}

func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		NameID:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

// Implements Msg.
func (msg MsgBuyName) Type() string { return "nameservice" }
func (msg MsgSetName) Name() string { return "buy_name"}

// Implements Msg.
func (msg MsgBuyName) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.NameID) == 0 {
		return sdk.ErrUnknownRequest("Name and Value cannot be empty")
	}
	if !msg.Bid.IsPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// Implements Msg.
func (msg MsgBuyName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
```

In the `handler.go` file, we add the `MsgBuyName` handler to the module router, so it now looks like this:
```go
// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
```

And we add the actual handle function to the `handler.go` file:
```go
// Handle MsgBuyName
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) sdk.Result {
	if keeper.GetPrice(ctx, msg.NameID).IsGTE(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.NameID) {
		_, err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.NameID), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.NameID, msg.Buyer)
	keeper.SetPrice(ctx, msg.NameID, msg.Bid)
	return sdk.Result{}
}
```
In this function, we check to make sure the bid is higher than the current price.  If it is, we check to see whether the name already has an owner.  If it does, they get transferred the money from the Buyer.  If it doesn't, the money just gets burned from the buyer.  If either `SubtractCoins` or `SendCoins` returns a non-nil error, the handler throws an error, reverting the transaction.  Otherwise, we set the buyer to the new owner, set the new price to be the current bid, and return.
