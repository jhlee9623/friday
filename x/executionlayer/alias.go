package executionlayer

import (
	"github.com/hdac-io/friday/x/executionlayer/types"
)

const (
	ModuleName      = types.ModuleName
	RouterKey       = types.RouterKey
	HashMapStoreKey = types.HashMapStoreKey
)

var (
	// function aliases
	NewMsgExecute  = types.NewMsgExecute
	NewMsgTransfer = types.NewMsgTransfer
	NewMsgBond     = types.NewMsgBond
	NewMsgUnBond   = types.NewMsgUnBond
	RegisterCodec  = types.RegisterCodec
	NewUnitHashMap = types.NewUnitHashMap

	// variable aliases
	ModuleCdc    = types.ModuleCdc
	ValidatorKey = types.ValidatorKey
)

type (
	MsgExecute                = types.MsgExecute
	MsgBond                   = types.MsgBond
	MsgUnBond                 = types.MsgUnBond
	MsgCreateValidator        = types.MsgCreateValidator
	MsgEditValidator          = types.MsgEditValidator
	QueryExecutionLayer       = types.QueryExecutionLayer
	UnitHashMap               = types.UnitHashMap
	QueryExecutionLayerResp   = types.QueryExecutionLayerResp
	QueryExecutionLayerDetail = types.QueryExecutionLayerDetail
	QueryGetBalance           = types.QueryGetBalance
	QueryGetBalanceDetail     = types.QueryGetBalanceDetail
)
