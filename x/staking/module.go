package staking

import (
	"context"
	"encoding/json"
	"fmt"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking"
	sdkcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	sdkkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/simulation"
	sdktypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/sedaprotocol/seda-chain/x/staking/keeper"
	"github.com/sedaprotocol/seda-chain/x/staking/types"
)

const (
	consensusVersion uint64 = 5
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic defines the basic application module used by the staking module.
type AppModuleBasic struct {
	cdc codec.Codec
	ak  types.AccountKeeper
}

// Name returns the staking module's name.
func (AppModuleBasic) Name() string {
	return sdktypes.ModuleName
}

// RegisterLegacyAminoCodec registers the staking module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	sdktypes.RegisterLegacyAminoCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (AppModuleBasic) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// TO-DO to be activated later
	// types.RegisterInterfaces(registry)
	// err := registry.EnsureRegistered(&types.MsgCreateValidatorWithVRF{})
	// if err != nil {
	// 	fmt.Println("error")
	// }

	sdktypes.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the staking
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(sdktypes.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the staking module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data sdktypes.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", sdktypes.ModuleName, err)
	}

	return staking.ValidateGenesis(&data)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the staking module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := sdktypes.RegisterQueryHandlerClient(context.Background(), mux, sdktypes.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns the root tx command for the staking module.
func (amb AppModuleBasic) GetTxCmd() *cobra.Command {
	// TO-DO to be activated later
	// return cli.NewTxCmd(amb.cdc.InterfaceRegistry().SigningContext().ValidatorAddressCodec(), amb.cdc.InterfaceRegistry().SigningContext().AddressCodec())

	return sdkcli.NewTxCmd(amb.cdc.InterfaceRegistry().SigningContext().ValidatorAddressCodec(), amb.cdc.InterfaceRegistry().SigningContext().AddressCodec())
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements an application module for the staking module.
type AppModule struct {
	AppModuleBasic

	keeper           *keeper.Keeper
	accountKeeper    types.AccountKeeper
	bankKeeper       sdktypes.BankKeeper
	randomnessKeeper types.RandomnessKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	cdc codec.Codec,
	keeper *keeper.Keeper,
	ak types.AccountKeeper,
	bk sdktypes.BankKeeper,
	rk types.RandomnessKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic:   AppModuleBasic{cdc: cdc, ak: ak},
		keeper:           keeper,
		accountKeeper:    ak,
		bankKeeper:       bk,
		randomnessKeeper: rk,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	sdkMsgServer := NewMsgServerImpl(am.keeper.Keeper, am.accountKeeper)
	sdktypes.RegisterMsgServer(cfg.MsgServer(), sdkMsgServer)

	// TO-DO to be activated later
	// types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(sdkMsgServer, am.keeper, am.accountKeeper, am.randomnessKeeper))

	querier := sdkkeeper.Querier{Keeper: am.keeper.Keeper}
	sdktypes.RegisterQueryServer(cfg.QueryServer(), querier)
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// RegisterInvariants registers the staking module invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	sdkkeeper.RegisterInvariants(ir, am.keeper.Keeper)
}

// InitGenesis performs genesis initialization for the staking module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState sdktypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	return am.keeper.InitGenesis(ctx, &genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the staking
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(am.keeper.ExportGenesis(ctx))
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return consensusVersion }

// BeginBlock returns the begin blocker for the staking module.
func (am AppModule) BeginBlock(ctx context.Context) error {
	return am.keeper.BeginBlocker(ctx)
}

// EndBlock returns the end blocker for the staking module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	return am.keeper.EndBlocker(ctx)
}

// AppModuleSimulation functions

// GenerateGenesisState creates a randomized GenState of the staking module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	simulation.RandomizedGenState(simState)
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return simulation.ProposalMsgs()
}

// RegisterStoreDecoder registers a decoder for staking module's types
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {
	sdr[sdktypes.StoreKey] = simulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the staking module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	return simulation.WeightedOperations(
		simState.AppParams, simState.Cdc, simState.TxConfig,
		am.accountKeeper, am.bankKeeper, am.keeper.Keeper,
	)
}
