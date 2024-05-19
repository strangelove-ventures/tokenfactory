package decorators

import (
	"fmt"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/strangelove-ventures/tokenfactory/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/strangelove-ventures/tokenfactory/x/tokenfactory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgBeforeSendHook struct {
	TokenFactoryKeeper keeper.Keeper
}

func NewMsgBeforeSendHook(k keeper.Keeper) MsgBeforeSendHook {
	return MsgBeforeSendHook{
		TokenFactoryKeeper: k,
	}
}

func (mfd MsgBeforeSendHook) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	mfd.applyDenomHookIfApplicable(tx.GetMsgs())
	return next(ctx, tx, simulate)
}

func (mfd MsgBeforeSendHook) applyDenomHookIfApplicable(msgs []sdk.Msg) {
	for _, msg := range msgs {
		// TODO: authz nested check here

		if m, ok := msg.(*banktypes.MsgSend); ok {
			mfd.performAction(m)
		}
	}
}

func (mfd MsgBeforeSendHook) performAction(msg *banktypes.MsgSend) {
	sender := msg.FromAddress
	recipient := msg.ToAddress
	coins := msg.Amount

	for _, coin := range coins {
		coin := coin

		// validate is a tokenfactory denom
		_, _, err := tokenfactorytypes.DeconstructDenom(coin.Denom)
		if err != nil {
			continue
		}

		// This is a validate tokenfactory denom being sent, check if it is registered for hooks
		fireEvent(sender, recipient, coin)
	}
}

func fireEvent(sender, recipient string, coin sdk.Coin) error {
	// Perform the SudoMsg execute if it is registered.
	fmt.Printf("Fire event for %s -> %s: %s\n", sender, recipient, coin)
	return nil
}
