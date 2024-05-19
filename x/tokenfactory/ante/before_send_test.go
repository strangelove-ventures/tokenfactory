package decorators_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"

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
	// generate a private/public key pair and get the respective address
	pk1 := ed25519.GenPrivKey().PubKey()
	fromAddr := sdk.AccAddress(pk1.Address())

	pk2 := ed25519.GenPrivKey().PubKey()
	toAddr := sdk.AccAddress(pk2.Address())

	// app, ctx := bindings.SetupCustomApp(t, addr1)
	ctx, chain := app.Setup(t)

	token, err := chain.TokenFactoryKeeper.CreateDenom(ctx, fromAddr.String(), "bitcoin")
	require.NoError(t, err)

	fmt.Println(token)

	// Create change rate decorator
	anteHandler := ante.NewMsgBeforeSendHook(chain.TokenFactoryKeeper)

	bankSendMsg := banktypes.MsgSend{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(token, sdkmath.NewInt(100))),
	}

	if _, err := anteHandler.AnteHandle(ctx, NewMockTx(&bankSendMsg), false, EmptyAnte); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// change this to a hooks contract (on call increase, ref clock)
	// reflect := instantiateReflectContract(t, ctx, app, actor)
	// require.NotEmpty(t, reflect)

}
