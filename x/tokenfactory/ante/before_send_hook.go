package decorators

import (
	"fmt"

	authz "github.com/cosmos/cosmos-sdk/x/authz"
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
	if err := mfd.applyDenomHookIfApplicable(tx.GetMsgs()); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func (mfd MsgBeforeSendHook) applyDenomHookIfApplicable(msgs []sdk.Msg) error {
	for _, msg := range msgs {
		// Check if an authz message, loop through all inner messages, and recursively call this function
		if execMsg, ok := msg.(*authz.MsgExec); ok {
			msgs, err := execMsg.GetMessages()
			if err != nil {
				return err
			}

			// Recursively call this function with the inner messages
			if err = mfd.applyDenomHookIfApplicable(msgs); err != nil {
				return err
			}
		}

		// if it is a bank message, perform the action if it's a tokenfactory denom
		if m, ok := msg.(*banktypes.MsgSend); ok {
			// mfd.performAction(m)

			sender := m.FromAddress
			recipient := m.ToAddress
			coins := m.Amount

			for _, coin := range coins {
				coin := coin

				// validate is a tokenfactory denom
				_, _, err := tokenfactorytypes.DeconstructDenom(coin.Denom)
				if err != nil {
					continue
				}

				// TODO: see if hooks is registered

				// This is a validate tokenfactory denom being sent, check if it is registered for hooks
				fireEvent(sender, recipient, coin)
			}
		}
	}

	return nil
}

func fireEvent(sender, recipient string, coin sdk.Coin) error {
	// Perform the SudoMsg execute if it is registered.
	fmt.Printf("Fire event for %s -> %s: %s\n", sender, recipient, coin)
	return nil
}
