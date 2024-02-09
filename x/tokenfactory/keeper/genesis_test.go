package keeper_test

import (
	"github.com/cosmos/tokenfactory/x/tokenfactory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (suite *KeeperTestSuite) TestGenesis() {
	genesisState := types.GenesisState{
		FactoryDenoms: []types.GenesisDenom{
			{
				Denom: "factory/tokenfactory1t7egva48prqmzl59x5ngv4zx0dtrwewcvz09h0/bitcoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "tokenfactory1t7egva48prqmzl59x5ngv4zx0dtrwewcvz09h0",
				},
			},
			{
				Denom: "factory/tokenfactory1t7egva48prqmzl59x5ngv4zx0dtrwewcvz09h0/diff-admin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "tokenfactory15czt5nhlnvayqq37xun9s9yus0d6y26d8uh5qf",
				},
			},
			{
				Denom: "factory/tokenfactory1t7egva48prqmzl59x5ngv4zx0dtrwewcvz09h0/litecoin",
				AuthorityMetadata: types.DenomAuthorityMetadata{
					Admin: "tokenfactory1t7egva48prqmzl59x5ngv4zx0dtrwewcvz09h0",
				},
			},
		},
	}

	suite.SetupTestForInitGenesis()
	app := suite.App

	// Test both with bank denom metadata set, and not set.
	for i, denom := range genesisState.FactoryDenoms {
		// hacky, sets bank metadata to exist if i != 0, to cover both cases.
		if i != 0 {
			app.BankKeeper.SetDenomMetaData(suite.Ctx, banktypes.Metadata{Base: denom.GetDenom()})
		}
	}

	if err := app.TokenFactoryKeeper.SetParams(suite.Ctx, types.Params{DenomCreationFee: sdk.Coins{sdk.NewInt64Coin("stake", 100)}}); err != nil {
		panic(err)
	}
	app.TokenFactoryKeeper.InitGenesis(suite.Ctx, genesisState)

	exportedGenesis := app.TokenFactoryKeeper.ExportGenesis(suite.Ctx)
	suite.Require().NotNil(exportedGenesis)
	suite.Require().Equal(genesisState, *exportedGenesis)
}
