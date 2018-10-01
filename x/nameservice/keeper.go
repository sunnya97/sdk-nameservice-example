package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper - handlers sets/gets of custom variables for your module
type Keeper struct {
	bk bank.Keeper

	key sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc *codec.Codec // The codec codec for binary encoding/decoding.
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk bank.Keeper) Keeper {
	return Keeper{
		bk:  bk,
		key: key,
		cdc: cdc,
	}
}

// GetTrend - returns the current cool trend
func (k Keeper) GetValue(ctx sdk.Context, domain string) string {
	store := ctx.KVStore(k.key)
	bz := store.Get(ValueKey(domain))

	return string(bz)
}

func (k Keeper) SetValue(ctx sdk.Context, domain string, value string) {
	store := ctx.KVStore(k.key)
	store.Set(ValueKey(domain), []byte(value))
}

// HasOwner - returns whether or not the domain already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, domain string) bool {
	store := ctx.KVStore(k.key)
	return store.Has(OwnerKey(domain))
}

// GetOwner - get the current owner of a domain
func (k Keeper) GetOwner(ctx sdk.Context, domain string) sdk.AccAddress {
	store := ctx.KVStore(k.key)
	bz := store.Get(OwnerKey(domain))
	return bz
}

// SetOwner - sets the current owner of a domain
func (k Keeper) SetOwner(ctx sdk.Context, domain string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.key)
	store.Set(OwnerKey(domain), owner)
}

// GetPrice - gets the current price of a domain.  If price doesn't exist yet, set to 1steak.
func (k Keeper) GetPrice(ctx sdk.Context, domain string) (price sdk.Coins) {
	store := ctx.KVStore(k.key)
	bz := store.Get(PriceKey(domain))
	if bz == nil {
		return
	}
	k.cdc.MustUnmarshalBinary(bz, &price)
	return
}

// SetPrice - sets the current price of a domain
func (k Keeper) SetPrice(ctx sdk.Context, domain string, price sdk.Coins) {
	store := ctx.KVStore(k.key)
	store.Set([]byte(domain), k.cdc.MustMarshalBinary(price))
}
