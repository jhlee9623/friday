package executionlayer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hdac-io/casperlabs-ee-grpc-go-util/grpc"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/casper/consensus/state"
	"github.com/hdac-io/casperlabs-ee-grpc-go-util/protobuf/io/casperlabs/ipc"

	"github.com/hdac-io/friday/codec"
	sdk "github.com/hdac-io/friday/types"
	"github.com/hdac-io/friday/x/executionlayer/types"
)

type ExecutionLayerKeeper struct {
	HashMapStoreKey sdk.StoreKey
	client          ipc.ExecutionEngineServiceClient
	cdc             *codec.Codec
}

func NewExecutionLayerKeeper(
	cdc *codec.Codec, hashMapStoreKey sdk.StoreKey, path string) ExecutionLayerKeeper {
	return ExecutionLayerKeeper{
		HashMapStoreKey: hashMapStoreKey,
		client:          grpc.Connect(path),
		cdc:             cdc,
	}
}

func (k ExecutionLayerKeeper) MustGetProtocolVersion(ctx sdk.Context) state.ProtocolVersion {
	genesisConf := k.GetGenesisConf(ctx)
	pv, err := types.ToProtocolVersion(genesisConf.Genesis.ProtocolVersion)
	if err != nil {
		panic(fmt.Errorf("System has invalid protocol version: %v", err))
	}
	return *pv
}

// SetUnitHashMap map unitHash to blockHash
func (k ExecutionLayerKeeper) SetUnitHashMap(ctx sdk.Context, blockHash []byte, unitHash UnitHashMap) bool {
	if bytes.Equal(blockHash, []byte{}) {
		blockHash = []byte("genesis")
	}
	if bytes.Equal(unitHash.EEState, []byte{}) || len(unitHash.EEState) != 32 {
		return false
	}

	unitBytes, err := k.cdc.MarshalBinaryBare(unitHash)
	if err != nil {
		return false
	}

	store := ctx.KVStore(k.HashMapStoreKey)
	store.Set(blockHash, unitBytes)

	return true
}

// GetUnitHashMap returns a UnitHashMap for blockHash
func (k ExecutionLayerKeeper) GetUnitHashMap(ctx sdk.Context, blockHash []byte) UnitHashMap {
	if bytes.Equal(blockHash, []byte{}) {
		blockHash = []byte("genesis")
	}
	store := ctx.KVStore(k.HashMapStoreKey)
	unitBytes := store.Get(blockHash)
	var unit UnitHashMap
	k.cdc.UnmarshalBinaryBare(unitBytes, &unit)
	return unit
}

// SetEEState map eeState to blockHash
func (k ExecutionLayerKeeper) SetEEState(ctx sdk.Context, blockHash []byte, eeState []byte) bool {
	if bytes.Equal(blockHash, []byte{}) {
		blockHash = []byte("genesis")
	}
	if bytes.Equal(eeState, []byte{}) || len(eeState) != 32 {
		return false
	}

	unit := UnitHashMap{
		EEState: eeState,
	}

	unitBytes, err := k.cdc.MarshalBinaryBare(unit)
	if err != nil {
		return false
	}

	store := ctx.KVStore(k.HashMapStoreKey)
	store.Set(blockHash, unitBytes)

	return true
}

// GetEEState returns a eeState for blockHash
func (k ExecutionLayerKeeper) GetEEState(ctx sdk.Context, blockHash []byte) []byte {
	if bytes.Equal(blockHash, []byte{}) {
		blockHash = []byte("genesis")
	}
	store := ctx.KVStore(k.HashMapStoreKey)
	unitBytes := store.Get(blockHash)
	var unit UnitHashMap
	k.cdc.UnmarshalBinaryBare(unitBytes, &unit)
	return unit.EEState
}

// TODO: change KeyType from string to typed enum.
func toBytes(keyType string, key string) ([]byte, error) {
	var bytes []byte
	var err error = nil

	switch keyType {
	case "address":
		bech32addr, err := sdk.AccAddressFromBech32(key)
		if err != nil {
			return nil, err
		}
		bytes = types.ToPublicKey(bech32addr)
	case "uref", "local", "hash":
		bytes, err = hex.DecodeString(key)
	default:
		err = fmt.Errorf("Unknown QueryKey type: %v", keyType)
	}

	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GetQueryResult queries with whole parameters
func (k ExecutionLayerKeeper) GetQueryResult(ctx sdk.Context,
	stateHash []byte, keyType string, keyData string, path string) (state.Value, error) {
	arrPath := strings.Split(path, "/")

	protocolVersion := k.MustGetProtocolVersion(ctx)
	keyDataBytes, err := toBytes(keyType, keyData)
	if err != nil {
		return state.Value{}, err
	}
	res, errstr := grpc.Query(k.client, stateHash, keyType, keyDataBytes, arrPath, &protocolVersion)
	if errstr != "" {
		return state.Value{}, fmt.Errorf(errstr)
	}

	return *res, nil
}

// GetQueryResultSimple queries without state hash.
// State hash comes from Tendermint block state - EE state mapping DB
func (k ExecutionLayerKeeper) GetQueryResultSimple(ctx sdk.Context,
	keyType string, keyData string, path string) (state.Value, error) {
	unitHash := k.GetUnitHashMap(ctx, k.GetCurrentBlockHash(ctx))
	arrPath := strings.Split(path, "/")

	keyDataBytes, err := toBytes(keyType, keyData)
	if err != nil {
		return state.Value{}, err
	}

	protocolVersion := k.MustGetProtocolVersion(ctx)
	res, errstr := grpc.Query(k.client, unitHash.EEState, keyType, keyDataBytes, arrPath, &protocolVersion)
	if errstr != "" {
		return state.Value{}, fmt.Errorf(errstr)
	}

	return *res, nil
}

// GetQueryBalanceResult queries with whole parameters
func (k ExecutionLayerKeeper) GetQueryBalanceResult(ctx sdk.Context, stateHash []byte, address types.PublicKey) (string, error) {
	protocolVersion := k.MustGetProtocolVersion(ctx)
	res, err := grpc.QueryBalance(k.client, stateHash, address, &protocolVersion)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return res, nil
}

// GetQueryBalanceResultSimple queries with whole parameters
func (k ExecutionLayerKeeper) GetQueryBalanceResultSimple(ctx sdk.Context, address types.PublicKey) (string, error) {
	unitHash := k.GetUnitHashMap(ctx, k.GetCurrentBlockHash(ctx))
	protocolVersion := k.MustGetProtocolVersion(ctx)
	res, err := grpc.QueryBalance(k.client, unitHash.EEState, address, &protocolVersion)
	if err != "" {
		return "", fmt.Errorf(err)
	}

	return res, nil
}

// GetGenesisConf retrieves GenesisConf from sdk store
func (k ExecutionLayerKeeper) GetGenesisConf(ctx sdk.Context) types.GenesisConf {
	store := ctx.KVStore(k.HashMapStoreKey)
	genesisConfBytes := store.Get([]byte("genesisconf"))

	var genesisConf types.GenesisConf
	k.cdc.UnmarshalBinaryBare(genesisConfBytes, &genesisConf)
	return genesisConf
}

// SetGenesisConf saves GenesisConf in sdk store
func (k ExecutionLayerKeeper) SetGenesisConf(ctx sdk.Context, genesisConf types.GenesisConf) {
	store := ctx.KVStore(k.HashMapStoreKey)
	genesisConfBytes := k.cdc.MustMarshalBinaryBare(genesisConf)
	store.Set([]byte("genesisconf"), genesisConfBytes)
}

// GetGenesisAccounts retrieves GenesisAccounts in sdk store
func (k ExecutionLayerKeeper) GetGenesisAccounts(ctx sdk.Context) []types.Account {
	store := ctx.KVStore(k.HashMapStoreKey)
	genesisAccountsBytes := store.Get([]byte("genesisaccounts"))
	if genesisAccountsBytes == nil {
		return nil
	}
	var genesisAccounts []types.Account
	k.cdc.UnmarshalBinaryBare(genesisAccountsBytes, &genesisAccounts)
	return genesisAccounts
}

// SetGenesisAccounts saves GenesisAccounts in sdk store
func (k ExecutionLayerKeeper) SetGenesisAccounts(ctx sdk.Context, accounts []types.Account) {
	if accounts == nil {
		panic(fmt.Errorf("Nil is not allowed for GenesisAccounts"))
	}
	store := ctx.KVStore(k.HashMapStoreKey)
	genesisAccountsBytes := k.cdc.MustMarshalBinaryBare(accounts)
	store.Set([]byte("genesisaccounts"), genesisAccountsBytes)
}

// GetChainName retrieves ChainName in sdk store
func (k ExecutionLayerKeeper) GetChainName(ctx sdk.Context) string {
	store := ctx.KVStore(k.HashMapStoreKey)
	chainNameBytes := store.Get([]byte("chainname"))
	return string(chainNameBytes)
}

// SetChainName saves ChainName in sdk store
func (k ExecutionLayerKeeper) SetChainName(ctx sdk.Context, chainName string) {
	if chainName == "" {
		panic(fmt.Errorf("Empty string is not allowed for ChainName"))
	}
	store := ctx.KVStore(k.HashMapStoreKey)
	store.Set([]byte("chainname"), []byte(chainName))
}

// GetCurrentBlockHash returns current block hash
func (k ExecutionLayerKeeper) GetCurrentBlockHash(ctx sdk.Context) []byte {
	store := ctx.KVStore(k.HashMapStoreKey)
	blockHash := store.Get([]byte("currentblockhash"))

	return blockHash
}

// SetCurrentBlockHash saves current block hash
func (k ExecutionLayerKeeper) SetCurrentBlockHash(ctx sdk.Context, blockHash []byte) {
	store := ctx.KVStore(k.HashMapStoreKey)
	store.Set([]byte("currentblockhash"), blockHash)
}
