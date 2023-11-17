package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sedaprotocol/seda-chain/x/randomness/types"
)

func CustomQuerier(randomnessKeeper *Querier) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		ctx.Logger().Info("Received request")
		ctx.Logger().Info(fmt.Sprintf("% x", string(request)))

		var contractQuery types.QuerySeedRequest
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, sdkerrors.Wrap(err, "seed query")
		}

		ctx.Logger().Info("Parsed json")

		seedQueryResponse, err := randomnessKeeper.Seed(ctx, &contractQuery)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "seed query")
		}

		ctx.Logger().Info("Get seed")
		ctx.Logger().Info(seedQueryResponse.Seed)

		bz, err := json.Marshal(seedQueryResponse)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "seed query response")
		}
		return bz, nil

	}
}