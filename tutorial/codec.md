## Codec File

In your module's folder create a `codec.go` file.  This allows Amino to register the `MsgSetName` and `MsgBuyName`.

```go
package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetName{}, "nameservice/SetName", nil)
	cdc.RegisterConcrete(MsgBuyName{}, "nameservice/BuyName", nil)
}
```
