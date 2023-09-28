package keeper_test

import (
	"encoding/hex"
	"os"

	"github.com/CosmWasm/wasmd/x/wasm/ioutils"
	"github.com/hyperledger/burrow/crypto"

	"github.com/sedaprotocol/seda-chain/x/wasm-storage/types"
)

func (s *KeeperTestSuite) TestStoreDataRequestWasm() {
	regWasm, err := os.ReadFile("test_utils/hello-world.wasm")
	s.Require().NoError(err)
	regWasmZipped, err := ioutils.GzipIt(regWasm)
	s.Require().NoError(err)

	oversizedWasm, err := os.ReadFile("test_utils/oversized.wasm")
	s.Require().NoError(err)
	oversizedWasmZipped, err := ioutils.GzipIt(oversizedWasm)
	s.Require().NoError(err)

	cases := []struct {
		name      string
		preRun    func()
		input     types.MsgStoreDataRequestWasm
		expErr    bool
		expErrMsg string
		expOutput types.MsgStoreDataRequestWasmResponse
	}{
		{
			name:   "happy path",
			preRun: func() {},
			input: types.MsgStoreDataRequestWasm{
				Sender:   s.authority,
				Wasm:     regWasmZipped,
				WasmType: types.WasmTypeDataRequest,
			},
			expErr: false,
			expOutput: types.MsgStoreDataRequestWasmResponse{
				Hash: hex.EncodeToString(crypto.Keccak256(regWasm)),
			},
		},
		{
			name: "Data Request wasm already exist",
			input: types.MsgStoreDataRequestWasm{
				Sender:   s.authority,
				Wasm:     regWasmZipped,
				WasmType: types.WasmTypeDataRequest,
			},
			preRun: func() {
				input := types.MsgStoreDataRequestWasm{
					Sender:   s.authority,
					Wasm:     regWasmZipped,
					WasmType: types.WasmTypeDataRequest,
				}
				_, err := s.msgSrvr.StoreDataRequestWasm(s.ctx, &input)
				s.Require().Nil(err)
			},
			expErr:    true,
			expErrMsg: "data Request Wasm with given hash already exists",
		},
		// TO-DO: Add after migrating ValidateBasic logic
		// {
		// 	name: "inconsistent Wasm type",
		// 	input: types.MsgStoreDataRequestWasm{
		// 		Sender:   s.authority,
		// 		Wasm:     regWasmZipped,
		// 		WasmType: types.WasmTypeRelayer,
		// 	},
		// 	preRun:    func() {},
		// 	expErr:    true,
		// 	expErrMsg: "not a Data Request Wasm",
		// },
		{
			name: "unzipped Wasm",
			input: types.MsgStoreDataRequestWasm{
				Sender:   s.authority,
				Wasm:     regWasm,
				WasmType: types.WasmTypeDataRequest,
			},
			preRun:    func() {},
			expErr:    true,
			expErrMsg: "wasm is not gzip compressed",
		},
		{
			name: "oversized Wasm",
			input: types.MsgStoreDataRequestWasm{
				Sender:   s.authority,
				Wasm:     oversizedWasmZipped,
				WasmType: types.WasmTypeDataRequest,
			},
			preRun:    func() {},
			expErr:    true,
			expErrMsg: "",
		},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.preRun()
			res, err := s.msgSrvr.StoreDataRequestWasm(s.ctx, &tc.input)
			if tc.expErr {
				s.Require().ErrorContains(err, tc.expErrMsg)
			} else {
				s.Require().Nil(err)
				s.Require().Equal(tc.expOutput, *res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestStoreOverlayWasm() {
	regWasm, err := os.ReadFile("test_utils/hello-world.wasm")
	s.Require().NoError(err)
	regWasmZipped, err := ioutils.GzipIt(regWasm)
	s.Require().NoError(err)

	oversizedWasm, err := os.ReadFile("test_utils/oversized.wasm")
	s.Require().NoError(err)
	oversizedWasmZipped, err := ioutils.GzipIt(oversizedWasm)
	s.Require().NoError(err)

	cases := []struct {
		name      string
		preRun    func()
		input     types.MsgStoreOverlayWasm
		expErr    bool
		expErrMsg string
		expOutput types.MsgStoreOverlayWasmResponse
	}{
		{
			name: "happy path",
			input: types.MsgStoreOverlayWasm{
				Sender:   s.authority,
				Wasm:     regWasmZipped,
				WasmType: types.WasmTypeRelayer,
			},
			preRun:    func() {},
			expErr:    false,
			expErrMsg: "",
			expOutput: types.MsgStoreOverlayWasmResponse{
				Hash: hex.EncodeToString(crypto.Keccak256(regWasm)),
			},
		},
		{
			name: "invalid wasm type",
			input: types.MsgStoreOverlayWasm{
				Sender:   s.authority,
				Wasm:     regWasm,
				WasmType: types.WasmTypeDataRequest,
			},
			preRun:    func() {},
			expErr:    true,
			expErrMsg: "Overlay Wasm type must be data-request-executor or relayer",
		},
		{
			name: "invalid authority",
			input: types.MsgStoreOverlayWasm{
				Sender:   "cosmos16wfryel63g7axeamw68630wglalcnk3l0zuadc",
				Wasm:     regWasmZipped,
				WasmType: types.WasmTypeRelayer,
			},
			preRun:    func() {},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "Overlay wasm already exist",
			input: types.MsgStoreOverlayWasm{
				Sender:   s.authority,
				Wasm:     regWasmZipped,
				WasmType: types.WasmTypeRelayer,
			},
			preRun: func() {
				input := types.MsgStoreOverlayWasm{
					Sender:   s.authority,
					Wasm:     regWasmZipped,
					WasmType: types.WasmTypeRelayer,
				}
				_, err := s.msgSrvr.StoreOverlayWasm(s.ctx, &input)
				s.Require().Nil(err)
			},
			expErr:    true,
			expErrMsg: "overlay Wasm with given hash already exists",
		},
		{
			name: "unzipped Wasm",
			input: types.MsgStoreOverlayWasm{
				Sender:   s.authority,
				Wasm:     regWasm,
				WasmType: types.WasmTypeRelayer,
			},
			preRun:    func() {},
			expErr:    true,
			expErrMsg: "wasm is not gzip compressed",
		},
		{
			name: "oversized Wasm",
			input: types.MsgStoreOverlayWasm{
				Sender:   s.authority,
				Wasm:     oversizedWasmZipped,
				WasmType: types.WasmTypeRelayer,
			},
			preRun:    func() {},
			expErr:    true,
			expErrMsg: "",
		},
	}
	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.preRun()
			res, err := s.msgSrvr.StoreOverlayWasm(s.ctx, &tc.input)
			if tc.expErr {
				s.Require().ErrorContains(err, tc.expErrMsg)
			} else {
				s.Require().Nil(err)
				s.Require().Equal(tc.expOutput, *res)
			}
		})
	}
}
