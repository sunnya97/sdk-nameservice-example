package nameservice

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgBuyDomain struct {
	Domain string
	Value  string
	Bid    sdk.Coins
	Buyer  sdk.AccAddress
}

func NewMsgBuyDomain(domain string, value string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyDomain {
	return MsgBuyDomain{
		Domain: domain,
		Value:  value,
		Bid:    bid,
		Buyer:  buyer,
	}
}

// Implements Msg.
func (msg MsgBuyDomain) Name() string { return "buy_domain" }

// Implements Msg.
func (msg MsgBuyDomain) Type() string { return "nameservice" }

// Implements Msg.
func (msg MsgBuyDomain) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Domain) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Domain and Value cannot be empty")
	}
	if !msg.Bid.IsPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// Implements Msg.
func (msg MsgBuyDomain) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgBuyDomain) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
