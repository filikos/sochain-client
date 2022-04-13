package sochain

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_NetworkInfo_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	want := NetworkInfo{
		Status: "test",
		Data: NetworkData{
			Name:    "test",
			Acronym: "test",
			Network: "btc",
		},
	}

	gotNetwork := "btc"

	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/get_info/"+gotNetwork,
		httpmock.NewJsonResponderOrPanic(200, want))

	s := NewSochain()
	got, err := s.NetworkInfo(gotNetwork)
	assert.Nil(t, err)

	assert.True(t, reflect.DeepEqual(*got, want))
}

func Test_NetworkInfo_Error_StatusCode(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	gotNetwork := "btc"
	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/get_info/"+gotNetwork,
		httpmock.NewStringResponder(500, ""))

	s := NewSochain()
	_, err := s.NetworkInfo(gotNetwork)
	assert.NotNil(t, err)
}

func Test_NetworkInfo_Error_Unmarshal(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	gotNetwork := "btc"

	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/get_info/"+gotNetwork,
		httpmock.NewBytesResponder(200, nil))

	s := NewSochain()
	_, err := s.NetworkInfo(gotNetwork)
	assert.NotNil(t, err)
}

func Test_BlockHeight_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	want := Block{
		Status: "test",
		Data: BlockData{
			Network:   "btc",
			Blockhash: "somehash",
			BlockNo:   200000,
		},
	}

	gotNetwork := "btc"
	gotHeight := 200000
	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/get_block/"+gotNetwork+"/"+strconv.Itoa(gotHeight),
		httpmock.NewJsonResponderOrPanic(200, want))

	s := NewSochain()
	got, err := s.BlockHeight(gotNetwork, gotHeight)
	assert.Nil(t, err)

	assert.True(t, reflect.DeepEqual(*got, want))
}

func Test_BlockHeight_Error_StatusCode(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	gotNetwork := "btc"
	gotHeight := 200000
	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/get_block/"+gotNetwork+"/"+strconv.Itoa(gotHeight),
		httpmock.NewStringResponder(500, ""))

	s := NewSochain()
	_, err := s.BlockHeight(gotNetwork, gotHeight)
	assert.NotNil(t, err)
}

func Test_BlockHeight_Error_Unmarshal(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	gotNetwork := "btc"
	gotHeight := 200000
	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/get_block/"+gotNetwork+"/"+strconv.Itoa(gotHeight),
		httpmock.NewBytesResponder(200, nil))

	s := NewSochain()
	_, err := s.BlockHeight(gotNetwork, gotHeight)
	assert.NotNil(t, err)
}

func Test_BlockHash_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	want := Block{
		Status: "test",
		Data: BlockData{
			Network:   "btc",
			Blockhash: "somehash",
			BlockNo:   200000,
		},
	}

	gotNetwork := "btc"
	gotBlockhash := "200000"
	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/get_block/"+gotNetwork+"/"+gotBlockhash,
		httpmock.NewJsonResponderOrPanic(200, want))

	s := NewSochain()
	got, err := s.BlockHash(gotNetwork, gotBlockhash)
	assert.Nil(t, err)

	assert.True(t, reflect.DeepEqual(*got, want))
}

func Test_BlockHash_Error_StatusCode(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	gotNetwork := "btc"
	gotBlockhash := "200000"
	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/get_block/"+gotNetwork+"/"+gotBlockhash,
		httpmock.NewStringResponder(500, ""))

	s := NewSochain()
	_, err := s.BlockHash(gotNetwork, gotBlockhash)
	assert.NotNil(t, err)
}

func Test_BlockHash_Error_Unmarshal(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	gotNetwork := "btc"
	gotBlockhash := "200000"
	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/get_block/"+gotNetwork+"/"+gotBlockhash,
		httpmock.NewBytesResponder(200, nil))

	s := NewSochain()
	_, err := s.BlockHash(gotNetwork, gotBlockhash)
	assert.NotNil(t, err)
}

func Test_Transaction_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	want := Transaction{
		Status: "test",
		Code:   1,
		Data: TransactionData{
			Network: "btc",
			Txid:    "txid",
			Inputs: Inputs{
				{
					InputNo: 1,
					Address: "0x0000",
				},
			},
			Outputs: Outputs{
				{
					OutputNo: 1,
					Address:  "0x000",
				},
			},
		},
	}

	gotNetwork := "btc"
	gotTxHash := "200000"
	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/tx/"+gotNetwork+"/"+gotTxHash,
		httpmock.NewJsonResponderOrPanic(200, want))

	s := NewSochain()
	got, err := s.Transaction(gotNetwork, gotTxHash)
	assert.Nil(t, err)

	assert.True(t, reflect.DeepEqual(*got, want))
}

func Test_Transaction_Error_StatusCode(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	gotNetwork := "btc"
	gotTxHash := "200000"
	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/tx/"+gotNetwork+"/"+gotTxHash,
		httpmock.NewStringResponder(500, ""))

	s := NewSochain()
	_, err := s.Transaction(gotNetwork, gotTxHash)
	assert.NotNil(t, err)
}

func Test_Transaction_Error_Unmarshal(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	gotNetwork := "btc"
	gotTxHash := "200000"
	httpmock.RegisterResponder("GET", "https://sochain.com/api/v2/tx/"+gotNetwork+"/"+gotTxHash,
		httpmock.NewBytesResponder(200, nil))

	s := NewSochain()
	_, err := s.Transaction(gotNetwork, gotTxHash)
	assert.NotNil(t, err)
}
