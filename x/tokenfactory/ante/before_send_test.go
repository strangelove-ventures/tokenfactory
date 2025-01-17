package decorators_test

import (
	"fmt"
	"os"
	"testing"

	sdkmath "cosmossdk.io/math"

	"github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/strangelove-ventures/tokenfactory/app"
	"github.com/stretchr/testify/require"

	ante "github.com/strangelove-ventures/tokenfactory/x/tokenfactory/ante"
)

var (
	EmptyAnte = func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
		return ctx, nil
	}
)

func TestBeforeSend(t *testing.T) {
	fromAddr := createAccount()
	toAddr := createAccount()

	// Setup chain and single ante
	ctx, chain := app.Setup(t)
	anteHandler := ante.NewMsgBeforeSendHook(chain.TokenFactoryKeeper)

	// Create new token
	token, err := chain.TokenFactoryKeeper.CreateDenom(ctx, fromAddr.String(), "bitcoin")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	fmt.Println("token", token)

	// Mint new coins to the fromAddr
	initCoins := newCoins(token, 100_000)
	err = chain.BankKeeper.MintCoins(ctx, "mint", initCoins)
	require.NoError(t, err)
	err = chain.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", fromAddr, initCoins)
	require.NoError(t, err)

	// Setup the hooks contract
	// TODO: change me (and the readme) to proper contract sudo bindings
	codeID := storeHooksContract(t, ctx, chain, fromAddr)
	cAddr := instantiateBeforeSendContract(t, ctx, chain, fromAddr, codeID)
	fmt.Println(cAddr.String())

	// TODO: register contract -> token (match osmosis spec)

	// TODO: validate execution
	bankSendMsg := banktypes.MsgSend{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      newCoins(token, 117),
	}

	if _, err := anteHandler.AnteHandle(ctx, NewMockTx(&bankSendMsg), false, EmptyAnte); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

}

func storeHooksContract(t *testing.T, ctx sdk.Context, app *app.TokenFactoryApp, addr sdk.AccAddress) uint64 {
	wasmCode, err := os.ReadFile("../bindings/testdata/hooks.wasm")
	require.NoError(t, err)

	contractKeeper := keeper.NewDefaultPermissionKeeper(app.WasmKeeper)
	codeID, _, err := contractKeeper.Create(ctx, addr, wasmCode, nil)
	require.NoError(t, err)

	return codeID
}

func instantiateBeforeSendContract(t *testing.T, ctx sdk.Context, app *app.TokenFactoryApp, funder sdk.AccAddress, codeID uint64) sdk.AccAddress {
	initMsgBz := []byte("{}")
	contractKeeper := keeper.NewDefaultPermissionKeeper(app.WasmKeeper)

	addr, _, err := contractKeeper.Instantiate(ctx, codeID, funder, funder, initMsgBz, "hooks contract", nil)
	require.NoError(t, err)

	return addr
}

func createAccount() sdk.AccAddress {
	pk1 := ed25519.GenPrivKey().PubKey()
	return sdk.AccAddress(pk1.Address())
}

func newCoins(token string, amt int64) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(token, sdkmath.NewInt(amt)))
}
