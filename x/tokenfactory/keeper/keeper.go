package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/tokenfactory/x/tokenfactory/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	store "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey store.StoreKey

		accountKeeper       types.AccountKeeper
		bankKeeper          types.BankKeeper
		communityPoolKeeper types.CommunityPoolKeeper

		enabledCapabilities []string

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

// NewKeeper returns a new instance of the x/tokenfactory keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey store.StoreKey,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	communityPoolKeeper types.CommunityPoolKeeper,
	enabledCapabilities []string,
	authority string,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,

		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		communityPoolKeeper: communityPoolKeeper,

		authority: authority,

		enabledCapabilities: enabledCapabilities,
	}
}

func DefaultIsAdminFunc(ctx context.Context, addr string) bool {
	return false
}

func DefaultExtraSudoAllowedCheckFunc(ctx context.Context) bool {
	return false
}

// GetAuthority returns the x/mint module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) GetEnabledCapabilities() []string {
	return k.enabledCapabilities
}

func (k *Keeper) SetEnabledCapabilities(ctx sdk.Context, newCapabilities []string) {
	k.enabledCapabilities = newCapabilities
}

// Logger returns a logger for the x/tokenfactory module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetDenomPrefixStore returns the substore for a specific denom
func (k Keeper) GetDenomPrefixStore(ctx sdk.Context, denom string) store.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetDenomPrefixStore(denom))
}

// GetCreatorPrefixStore returns the substore for a specific creator address
func (k Keeper) GetCreatorPrefixStore(ctx sdk.Context, creator string) store.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetCreatorPrefix(creator))
}

// GetCreatorsPrefixStore returns the substore that contains a list of creators
func (k Keeper) GetCreatorsPrefixStore(ctx sdk.Context) store.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.GetCreatorsPrefix())
}
