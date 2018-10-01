package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterWire(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgBuyDomain{}, "nameservice/BuyDomain", nil)
}
