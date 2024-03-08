package keeper

import (
	"context"

	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"

	"github.com/cosmos/gogoproto/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAuthorityMetadata returns the authority metadata for a specific denom
func (k Keeper) GetAuthorityMetadata(ctx context.Context, denom string) (types.DenomAuthorityMetadata, error) {
	bz := k.GetDenomPrefixStore(sdk.UnwrapSDKContext(ctx), denom).Get([]byte(types.DenomAuthorityMetadataKey))

	metadata := types.DenomAuthorityMetadata{}
	err := proto.Unmarshal(bz, &metadata)
	if err != nil {
		return types.DenomAuthorityMetadata{}, err
	}
	return metadata, nil
}

// setAuthorityMetadata stores authority metadata for a specific denom
func (k Keeper) setAuthorityMetadata(ctx context.Context, denom string, metadata types.DenomAuthorityMetadata) error {
	err := metadata.Validate()
	if err != nil {
		return err
	}

	store := k.GetDenomPrefixStore(sdk.UnwrapSDKContext(ctx), denom)

	bz, err := proto.Marshal(&metadata)
	if err != nil {
		return err
	}

	store.Set([]byte(types.DenomAuthorityMetadataKey), bz)
	return nil
}

func (k Keeper) setAdmin(ctx context.Context, denom string, admin string) error {
	metadata, err := k.GetAuthorityMetadata(sdk.UnwrapSDKContext(ctx), denom)
	if err != nil {
		return err
	}

	metadata.Admin = admin

	return k.setAuthorityMetadata(ctx, denom, metadata)
}
