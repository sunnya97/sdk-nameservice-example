# Example SDK Module Tutorial

In this tutorial series, we are going to build a simplistic but functional module using the Cosmos SDK and learn the basics so that you can get started building your own modules and decentralized applications.  In this tutorial we will build a "nameservice", a mapping of strings to other strings (similar to Namecoin, ENS, or Handshake), in which to buy the name, the buyer has to pay the current owner more than the current owner paid to buy it!

All of the final source code for this tutorial project is in this directory, however, it is highly recommended that you follow along manually and try building the project yourself!

Docker image is available in `mossid/sdk-example:0.2`.

## The Keeper

The main core of a Cosmos SDK module is a piece called the Keeper. It is what handles interaction with the store, has references to other keepers, and often contains most of the core functionality of a module.  To begin, let's create a file called `keeper.go` and place it in a folder called `nameservice` that will hold our module.

### Keeper Struct

In this file, let's start by placing the following code.

```go
package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper - handlers sets/gets of custom variables for your module
type Keeper struct {
	coinKeeper bank.Keeper

	namesStoreKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
    ownersStoreKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.
    priceStoreKey sdk.StoreKey // The (unexposed) key used to access the store from the Context.

    cdc *wire.Codec // The wire codec for binary encoding/decoding.
}
```

Let's break this down.  The package name `nameservice` is the name of the package that this file is part of.  In Go, all code has to be part of a package.

Next we import the main `SDK` package and the `bank` module from the cosmos-sdk repository.

Next, we create the Keeper struct itself.  In this keeper there are a couple of key pieces:
- `bank.Keeper` - This is a reference to the Keeper from the module.  This allows code in this module to be able to call functions from the bank module.
- `sdk.StoreKey` - The SDK uses an object capabilities approach to accessing parts of the sections of the application state.  This is to allow developers to employ a least authority approach limiting the capabilities of a faulty or malicious module from affecting parts of state it doesn't need access to.  In this module, we will use three stores:
    - `namesStoreKey` - This is the main store that stores the value string that the name points to (i.e. The mapping from domain name -> IP Address)
    - `ownersStoreKey` - This store contains the current owner of this name
    - `priceStoreKey` - This store contains the price that the current owner paid. And buying of this name must spend more than the current owner.
- `*wire.Codec` - This is a pointer to the codec that is used by Amino to encode and decode binary structs.

### Getters and Setters

First let's add a function to set the string the a name resolves to.

```go
// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	store := ctx.KVStore(k.namesStoreKey)
	store.Set([]byte(name), []byte(value))
}
```

In this method on the Keeper, we first get the store object for the name resolutions using the the `namesStoreKey` from the Keeper.

Next, we insert the `<name, value>` pair into the store using its `.Set([]byte, []byte)` method.  As the store only takes `[]byte` while we have `string`s, we first need to cast the `string`s to `[]byte` and the use them as parameters into the `Set` method.

Next, let's add a method to actually resolve the names.

```go
// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	store := ctx.KVStore(k.namesStoreKey)
	bz := store.Get([]byte(name))
	return string(bz)
}
```

Here, like in the SetName method, we first get the store using the `StoreKey`.  Next, instead of using the `Set` method on the store key, we use the `.Get([]byte) []byte` method. As the parameter into the function, we pass the key, which is the `name` string casted to `[]byte`, and get back the result in the form of `[]byte`.  We cast this to a `string` and return the result.

We now add similar functions for Getting and Setting Owners.

```go
    // GetOwner - get the current owner of a name
    func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
        store := ctx.KVStore(k.ownersStoreKey)
        bz := store.Get([]byte(name))
        return bz
    }

    // SetOwner - sets the current owner of a name
    func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
        store := ctx.KVStore(k.ownersStoreKey)
        store.Set([]byte(name), owner)
    }

    // HasOwner - returns whether or not the name already has an owner
    func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
        store := ctx.KVStore(k.ownersStoreKey)
        bz := store.Get([]byte(name))
        return bz != nil
    }
```
Note that now, instead of accessing the the data from the `namesStoreKey` store, we now get it from the `ownersStoreKey` store.  Because sdk.AccAddress is a type alias for `[]byte`, we can natively cast to it.  We also added an extra function `HasOwner` that tells us whether a name already has an owner or not.

Finally, we will add a getter and setter for the Price of a name.

```go
    // GetPrice - gets the current price of a name.  If price doesn't exist yet, set to 1steak.
    func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
        if !k.HasOwner(ctx, name) {
            return sdk.Coins{sdk.NewInt64Coin("steak", 1)}
        }
        store := ctx.KVStore(k.priceStoreKey)
        bz := store.Get([]byte(name))
        var price sdk.Coins
        k.cdc.MustUnmarshalBinary(bz, &price)
        return price
    }

    // SetPrice - sets the current price of a name
    func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
        store := ctx.KVStore(k.priceStoreKey)
        store.Set([]byte(name), k.cdc.MustMarshalBinary(price))
    }
```
We put this data in the `priceStoreKey` store.  Note that `sdk.Coins` does not have it's own Bytes encoding, and so, to marshal and unmarshal the price for inserting and removing from store, we use Amino (read more about Amino here:  https://github.com/tendermint/go-amino/)

When getting the price for a name that has no owner (and thus no price), we will return 1steak as the price.


## Msgs and Handlers

Now that we have the keeper setup, it is time to built the msgs and handlers that actually allow users to buy and set names.

### Set Name

#### Msg

Let's first setup the different messages that a user can use to interact with this module.  The Cosmos SDK define a standard interface that all Msgs must satisfy:

```go
// Transactions messages must fulfill the Msg
type Msg interface {
	// Return the message type.
	// Must be alphanumeric or empty.
	Type() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() Error

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigners() []AccAddress
}
```

We'll start by defining `MsgSetName` in a new file called `msgs.go` in the `nameservice` package, a Msg that allows owner of an address to set the result of resolving a name.

```go
type MsgSetName struct {
	Name  string
	Value string
	Owner sdk.AccAddress
}

func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgBuyName{
		Name:  name,
		Value: value,
		Owner: owner,
	}
}
```
The `MsgSetName` has three attributes:
- `name` - The name trying to be set
- `value` - What the name resolves to
- `owner` - The owner of that name

```go
// Implements Msg.
func (msg MsgSetName) Type() string { return "nameservice" }
```
This is used by the SDK to route msgs to the proper module for handling.

```go
// Implements Msg.
func (msg MsgSetName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and Value cannot be empty")
	}
	return nil
}
```
This is used to provide some basic *stateless* checks on the validity of the msg.  In this case, we check that none of the attributes are empty.

```go
// Implements Msg.
func (msg MsgSetName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}
```
This defines how the Msg gets encoded for signing.  This should usually be in JSON and should not be modified in most cases.

```go
// Implements Msg.
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
```
This allows the Msg to define who's signature is required on a Tx in order for it to be valid.  In this case, for example, the `MsgSetName` requires that the `Owner` sign the transaction trying to reset what the name points to.

#### Handler

Now that we have the `MsgSetName` defined, we now have to define the handler that actually executes the Msg.

In a new file called `handler.go` in the `nameservice` package, we start off with:

```go
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
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
```

This is essentially a subrouter that directs messages coming into this module to the proper handler for the message.  At the moment, we only have one Msg/Handler.

In the same file, we define the function `handleMsgSetName`.

```go
// Handle MsgSetName
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	return sdk.Result{}                      // return
}
```
In this function we check to see if the Msg sender is actually the owner of the name (which we get using `keeper.GetOwner`).  If so, we let them set the name by calling the function on the keeper.  If not, we throw an error.

### Buy Name

Great, now owners can set names!  But what if a name doesn't have an owner yet?  We need a way for people to buy names!

#### Msg

We define the Msg for buying names and add it to the `msgs.go` file:

```go
type MsgBuyName struct {
	Name  string
	Bid   sdk.Coins
	Buyer sdk.AccAddress
}

func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

// Implements Msg.
func (msg MsgBuyName) Type() string { return "nameservice" }

// Implements Msg.
func (msg MsgBuyName) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
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
	if keeper.GetPrice(ctx, msg.Name).IsGTE(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.Name) {
		_, err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		_, _, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	return sdk.Result{}
}
```
In this function, we check to make sure the bid is higher than the current price.  If it is, we check to see whether the name already has an owner.  If it does, they get transferred the money from the Buyer.  If it doesn't, the money just gets burned from the buyer.  If either `SubtractCoins` or `SendCoins` returns a non-nil error, the handler throws an error, reverting the transaction.  Otherwise, we set the buyer to the new owner and return.


