# Msgs and Handlers

Now that you have the `Keeper` setup, it's time to built the `Msgs` and `Handlers` that actually allow users to buy and set names:

### `Msgs`

Represent state transitions that clients are submitting to the network. You can think of `Msgs` as similar to `Txns` in other blockchain systems. In the CosmosSDK `Msgs` must satisfy the below interface:

```go
// Transactions messages must fulfill the Msg
type Msg interface {
	// Return the message type.
	// Must be alphanumeric or empty.
	Type() string

	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Name() string

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

### `Handlers`

Define the action that needs to be taken (which stores need to get updated, how, and under what conditions) when a given `Msg` is received.

In this application you will have two types of `Msgs` that users can send to interact with the application state. They will each have an associated `Handler`:

* [`SetName`](./tutorial/set-name.md) - used to set the `value` of a `name` that is owned by the user
* [`BuyName`](./tutorial/buy-name.md) - used to buy a new `name` with `sdk.Coins` from an `address`

This completes the core logic of your `nameservice`. The next sections will focus on the rest of the core module. After that we focus on building an app using that module. In most cases, you'll be using similar boilerplate as well.
