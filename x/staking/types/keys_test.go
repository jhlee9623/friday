package types

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/hdac-io/tendermint/crypto/ed25519"

	sdk "github.com/hdac-io/friday/types"
)

var (
	keysPK1   = ed25519.GenPrivKeyFromSecret([]byte{1}).PubKey()
	keysPK2   = ed25519.GenPrivKeyFromSecret([]byte{2}).PubKey()
	keysPK3   = ed25519.GenPrivKeyFromSecret([]byte{3}).PubKey()
	keysAddr1 = keysPK1.Address()
	keysAddr2 = keysPK2.Address()
	keysAddr3 = keysPK3.Address()
)

func TestGetValidatorPowerRank(t *testing.T) {
	valAddr1 := sdk.ValAddress(keysAddr1)
	emptyDesc := Description{}
	val1 := NewValidator(valAddr1, keysPK1, emptyDesc)
	val1.Tokens = sdk.ZeroInt()
	val2, val3, val4 := val1, val1, val1
	val2.Tokens = sdk.TokensFromConsensusPower(1)
	val3.Tokens = sdk.TokensFromConsensusPower(10)
	x := new(big.Int).Exp(big.NewInt(2), big.NewInt(40), big.NewInt(0))
	val4.Tokens = sdk.TokensFromConsensusPower(x.Int64())

	tests := []struct {
		validator Validator
		wantHex   string
	}{
		{val1, "2300000000000000009c288ede7df62742fc3b7d0962045a8cef0f79f64642574ef021eaeba07e5398"},
		{val2, "2300000000000000019c288ede7df62742fc3b7d0962045a8cef0f79f64642574ef021eaeba07e5398"},
		{val3, "23000000000000000a9c288ede7df62742fc3b7d0962045a8cef0f79f64642574ef021eaeba07e5398"},
		{val4, "2300000100000000009c288ede7df62742fc3b7d0962045a8cef0f79f64642574ef021eaeba07e5398"},
	}
	for i, tt := range tests {
		got := hex.EncodeToString(getValidatorPowerRank(tt.validator))

		assert.Equal(t, tt.wantHex, got, "Keys did not match on test case %d", i)
	}
}

func TestGetREDByValDstIndexKey(t *testing.T) {
	tests := []struct {
		delAddr    sdk.AccAddress
		valSrcAddr sdk.ValAddress
		valDstAddr sdk.ValAddress
		wantHex    string
	}{
		{sdk.AccAddress(keysAddr1), sdk.ValAddress(keysAddr1), sdk.ValAddress(keysAddr1),
			"3663d771218209d8bd03c482f69dfba57310f08609b9bda8b10fde15145f81ac6763d771218209d8bd03c482f69dfba57310f08609b9bda8b10fde15145f81ac6763d771218209d8bd03c482f69dfba57310f08609b9bda8b10fde15145f81ac67"},
		{sdk.AccAddress(keysAddr1), sdk.ValAddress(keysAddr2), sdk.ValAddress(keysAddr3),
			"363ab62f0d93849be495e21e3e9013a517038f45bd20a77c87c18f39006eac915e63d771218209d8bd03c482f69dfba57310f08609b9bda8b10fde15145f81ac675ef3b5f25c54946d4a89fc0d09d2f126614540f247e01a01d9c19bfd853455f1"},
		{sdk.AccAddress(keysAddr2), sdk.ValAddress(keysAddr1), sdk.ValAddress(keysAddr3),
			"363ab62f0d93849be495e21e3e9013a517038f45bd20a77c87c18f39006eac915e5ef3b5f25c54946d4a89fc0d09d2f126614540f247e01a01d9c19bfd853455f163d771218209d8bd03c482f69dfba57310f08609b9bda8b10fde15145f81ac67"},
	}
	for i, tt := range tests {
		got := hex.EncodeToString(GetREDByValDstIndexKey(tt.delAddr, tt.valSrcAddr, tt.valDstAddr))

		assert.Equal(t, tt.wantHex, got, "Keys did not match on test case %d", i)
	}
}

func TestGetREDByValSrcIndexKey(t *testing.T) {
	tests := []struct {
		delAddr    sdk.AccAddress
		valSrcAddr sdk.ValAddress
		valDstAddr sdk.ValAddress
		wantHex    string
	}{
		{sdk.AccAddress(keysAddr1), sdk.ValAddress(keysAddr1), sdk.ValAddress(keysAddr1),
			"3563d771218209d8bd03c482f69dfba57310f08609b9bda8b10fde15145f81ac6763d771218209d8bd03c482f69dfba57310f08609b9bda8b10fde15145f81ac6763d771218209d8bd03c482f69dfba57310f08609b9bda8b10fde15145f81ac67"},
		{sdk.AccAddress(keysAddr1), sdk.ValAddress(keysAddr2), sdk.ValAddress(keysAddr3),
			"355ef3b5f25c54946d4a89fc0d09d2f126614540f247e01a01d9c19bfd853455f163d771218209d8bd03c482f69dfba57310f08609b9bda8b10fde15145f81ac673ab62f0d93849be495e21e3e9013a517038f45bd20a77c87c18f39006eac915e"},
		{sdk.AccAddress(keysAddr2), sdk.ValAddress(keysAddr1), sdk.ValAddress(keysAddr3),
			"3563d771218209d8bd03c482f69dfba57310f08609b9bda8b10fde15145f81ac675ef3b5f25c54946d4a89fc0d09d2f126614540f247e01a01d9c19bfd853455f13ab62f0d93849be495e21e3e9013a517038f45bd20a77c87c18f39006eac915e"},
	}
	for i, tt := range tests {
		got := hex.EncodeToString(GetREDByValSrcIndexKey(tt.delAddr, tt.valSrcAddr, tt.valDstAddr))

		assert.Equal(t, tt.wantHex, got, "Keys did not match on test case %d", i)
	}
}
