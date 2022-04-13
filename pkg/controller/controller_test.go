package controller

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sochain-client/pkg/sochain"
	mock_client "sochain-client/pkg/sochain/mock"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandleGetBlock(t *testing.T) {

	unixTime := 1231455600
	timeRFC3339 := time.Unix(int64(unixTime), 0).Format(time.RFC3339)

	tests := []struct {
		title                  string
		gotPathNetworkIDExists bool
		gotPathNetworkID       string
		gotHeightQuery         string
		gotBlockhashQuery      string
		mock                   func(m *mock_client.MockConnector)
		wantError              bool
		wantCode               int
		want                   *sochain.BlockResponse
	}{
		{
			title:                  "Error: missing path network id",
			wantError:              true,
			gotPathNetworkIDExists: false,
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: wrong path network id",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "someID",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: networkinfo client error",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			wantCode:               http.StatusBadRequest,
			mock: func(m *mock_client.MockConnector) {
				m.EXPECT().NetworkInfo("btc").Return(nil, errors.New("some"))
			},
		},
		{
			title:                  "Error: blockheight client error",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			wantCode:               http.StatusInternalServerError,
			mock: func(m *mock_client.MockConnector) {
				info := sochain.NetworkInfo{
					Data: sochain.NetworkData{
						Blocks: 1,
					},
				}
				m.EXPECT().NetworkInfo("btc").Return(&info, nil)

				m.EXPECT().BlockHeight("btc", 1).Return(nil, errors.New("some"))
			},
		},
		{
			title:                  "Error: blockheight client custom client error, bad request",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			wantCode:               http.StatusBadRequest,
			mock: func(m *mock_client.MockConnector) {
				info := sochain.NetworkInfo{
					Data: sochain.NetworkData{
						Blocks: 1,
					},
				}
				m.EXPECT().NetworkInfo("btc").Return(&info, nil)

				m.EXPECT().BlockHeight("btc", 1).Return(nil, *sochain.NewClientErr(errors.New("some"), http.StatusBadRequest))
			},
		},
		{
			title:                  "Error: blockheight client custom client error, not found",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			wantCode:               http.StatusNotFound,
			mock: func(m *mock_client.MockConnector) {
				info := sochain.NetworkInfo{
					Data: sochain.NetworkData{
						Blocks: 1,
					},
				}
				m.EXPECT().NetworkInfo("btc").Return(&info, nil)

				m.EXPECT().BlockHeight("btc", 1).Return(nil, *sochain.NewClientErr(errors.New("some"), http.StatusNotFound))
			},
		},
		{
			title:                  "Success: no transaction fetching errors",
			wantError:              false,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			wantCode:               http.StatusOK,
			mock: func(m *mock_client.MockConnector) {
				info := sochain.NetworkInfo{
					Data: sochain.NetworkData{
						Blocks: 1,
					},
				}
				m.EXPECT().NetworkInfo("btc").Return(&info, nil)

				block := sochain.Block{
					Data: sochain.BlockData{
						BlockNo:           1,
						Time:              unixTime,
						Blockhash:         "1",
						PreviousBlockhash: "1",
						NextBlockhash:     "1",
						Size:              1,
						// the same amount of calls to client.Transaction have to be passed to mocker, represented below
						Txs: []string{"1", "2"},
					},
				}
				m.EXPECT().BlockHeight("btc", 1).Return(&block, nil)

				m.EXPECT().Transaction("btc", "1").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "1",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)

				m.EXPECT().Transaction("btc", "2").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "2",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)
			},

			want: &sochain.BlockResponse{
				Blocknumber:  1,
				Timestamp:    timeRFC3339,
				Size:         1,
				PreviousHash: "1",
				NextHash:     "1",
				Transactions: sochain.TransactionResponses{
					{
						TxID:      "1",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
					{
						TxID:      "2",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
				},
			},
		},
		{
			title:                  "Success: with transaction fetching errors",
			wantError:              false,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			wantCode:               http.StatusOK,
			mock: func(m *mock_client.MockConnector) {
				info := sochain.NetworkInfo{
					Data: sochain.NetworkData{
						Blocks: 1,
					},
				}
				m.EXPECT().NetworkInfo("btc").Return(&info, nil)

				block := sochain.Block{
					Data: sochain.BlockData{
						BlockNo:           1,
						Time:              unixTime,
						Blockhash:         "1",
						PreviousBlockhash: "1",
						NextBlockhash:     "1",
						Size:              1,
						// the same amount of calls to client.Transaction have to be passed to mocker, represented below
						Txs: []string{"1", "2", "3", "4"},
					},
				}
				m.EXPECT().BlockHeight("btc", 1).Return(&block, nil)

				m.EXPECT().Transaction("btc", "1").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "1",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)

				m.EXPECT().Transaction("btc", "2").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "2",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)

				m.EXPECT().Transaction("btc", "3").Return(nil, errors.New("some"))
				m.EXPECT().Transaction("btc", "4").Return(nil, sochain.NewClientErr(errors.New("some"), http.StatusInternalServerError))
			},

			want: &sochain.BlockResponse{
				Blocknumber:  1,
				Timestamp:    timeRFC3339,
				Size:         1,
				PreviousHash: "1",
				NextHash:     "1",
				Transactions: sochain.TransactionResponses{
					{
						TxID:      "1",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
					{
						TxID:      "2",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
				},
			},
		},
		//query height provided
		{
			title:                  "Error: query blockheight invalid",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotHeightQuery:         "0",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: query blockheight invalid",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotHeightQuery:         "-1",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: query blockheight invalid",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotHeightQuery:         "noInt",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: query blockheight: client custom client error, bad request",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotHeightQuery:         "1",
			wantCode:               http.StatusBadRequest,
			mock: func(m *mock_client.MockConnector) {
				m.EXPECT().BlockHeight("btc", 1).Return(nil, *sochain.NewClientErr(errors.New("some"), http.StatusBadRequest))
			},
		},
		{
			title:                  "Error: query blockheight: client custom client error, not found",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotHeightQuery:         "1",
			wantCode:               http.StatusNotFound,
			mock: func(m *mock_client.MockConnector) {
				m.EXPECT().BlockHeight("btc", 1).Return(nil, *sochain.NewClientErr(errors.New("some"), http.StatusNotFound))
			},
		},
		{
			title:                  "Error: query blockheight: client custom client error",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotHeightQuery:         "1",
			wantCode:               http.StatusInternalServerError,
			mock: func(m *mock_client.MockConnector) {
				m.EXPECT().BlockHeight("btc", 1).Return(nil, errors.New("some"))
			},
		},
		{
			title:                  "Success: no transaction fetching errors",
			wantError:              false,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotHeightQuery:         "1",
			wantCode:               http.StatusOK,
			mock: func(m *mock_client.MockConnector) {
				block := sochain.Block{
					Data: sochain.BlockData{
						BlockNo:           1,
						Time:              unixTime,
						Blockhash:         "1",
						PreviousBlockhash: "1",
						NextBlockhash:     "1",
						Size:              1,
						// the same amount of calls to client.Transaction have to be passed to mocker, represented below
						Txs: []string{"1", "2"},
					},
				}
				m.EXPECT().BlockHeight("btc", 1).Return(&block, nil)

				m.EXPECT().Transaction("btc", "1").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "1",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)

				m.EXPECT().Transaction("btc", "2").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "2",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)
			},

			want: &sochain.BlockResponse{
				Blocknumber:  1,
				Timestamp:    timeRFC3339,
				Size:         1,
				PreviousHash: "1",
				NextHash:     "1",
				Transactions: sochain.TransactionResponses{
					{
						TxID:      "1",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
					{
						TxID:      "2",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
				},
			},
		},
		{
			title:                  "Success: with transaction fetching errors",
			wantError:              false,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotHeightQuery:         "1",
			wantCode:               http.StatusOK,
			mock: func(m *mock_client.MockConnector) {
				block := sochain.Block{
					Data: sochain.BlockData{
						BlockNo:           1,
						Time:              unixTime,
						Blockhash:         "1",
						PreviousBlockhash: "1",
						NextBlockhash:     "1",
						Size:              1,
						// the same amount of calls to client.Transaction have to be passed to mocker, represented below
						Txs: []string{"1", "2", "3", "4"},
					},
				}
				m.EXPECT().BlockHeight("btc", 1).Return(&block, nil)

				m.EXPECT().Transaction("btc", "1").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "1",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)

				m.EXPECT().Transaction("btc", "2").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "2",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)

				m.EXPECT().Transaction("btc", "3").Return(nil, errors.New("some"))
				m.EXPECT().Transaction("btc", "4").Return(nil, sochain.NewClientErr(errors.New("some"), http.StatusInternalServerError))
			},

			want: &sochain.BlockResponse{
				Blocknumber:  1,
				Timestamp:    timeRFC3339,
				Size:         1,
				PreviousHash: "1",
				NextHash:     "1",
				Transactions: sochain.TransactionResponses{
					{
						TxID:      "1",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
					{
						TxID:      "2",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
				},
			},
		},
		//blockhash height provided
		{
			title:                  "Error: query blockhash invalid",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotBlockhashQuery:      "0",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: query blockhash to long",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotBlockhashQuery:      "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876a",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: query blockhash to short",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotBlockhashQuery:      "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa87",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: query blockhash invalid characters",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotBlockhashQuery:      "eb6f76a4390f4e3cdcf8d2a73fc99d&401965abca0372f350bf9317a34a1aa87",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: query blockheight: client custom client error, bad request",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotBlockhashQuery:      "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
			wantCode:               http.StatusBadRequest,
			mock: func(m *mock_client.MockConnector) {
				m.EXPECT().BlockHash("btc", "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876").Return(nil, *sochain.NewClientErr(errors.New("some"), http.StatusBadRequest))
			},
		},
		{
			title:                  "Error: query blockheight: client custom client error, not found",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotBlockhashQuery:      "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
			wantCode:               http.StatusNotFound,
			mock: func(m *mock_client.MockConnector) {
				m.EXPECT().BlockHash("btc", "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876").Return(nil, *sochain.NewClientErr(errors.New("some"), http.StatusNotFound))
			},
		},
		{
			title:                  "Error: query blockheight: client custom client error",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotBlockhashQuery:      "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
			wantCode:               http.StatusInternalServerError,
			mock: func(m *mock_client.MockConnector) {
				m.EXPECT().BlockHash("btc", "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876").Return(nil, errors.New("some"))
			},
		},
		{
			title:                  "Success: query blockheight: no transaction fetching errors",
			wantError:              false,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotBlockhashQuery:      "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
			wantCode:               http.StatusOK,
			mock: func(m *mock_client.MockConnector) {
				block := sochain.Block{
					Data: sochain.BlockData{
						BlockNo:           1,
						Time:              unixTime,
						Blockhash:         "1",
						PreviousBlockhash: "1",
						NextBlockhash:     "1",
						Size:              1,
						// the same amount of calls to client.Transaction have to be passed to mocker, represented below
						Txs: []string{"1", "2"},
					},
				}
				m.EXPECT().BlockHash("btc", "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876").Return(&block, nil)

				m.EXPECT().Transaction("btc", "1").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "1",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)

				m.EXPECT().Transaction("btc", "2").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "2",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)
			},

			want: &sochain.BlockResponse{
				Blocknumber:  1,
				Timestamp:    timeRFC3339,
				Size:         1,
				PreviousHash: "1",
				NextHash:     "1",
				Transactions: sochain.TransactionResponses{
					{
						TxID:      "1",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
					{
						TxID:      "2",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
				},
			},
		},
		{
			title:                  "Success: with transaction fetching errors",
			wantError:              false,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotBlockhashQuery:      "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
			wantCode:               http.StatusOK,
			mock: func(m *mock_client.MockConnector) {
				block := sochain.Block{
					Data: sochain.BlockData{
						BlockNo:           1,
						Time:              unixTime,
						Blockhash:         "1",
						PreviousBlockhash: "1",
						NextBlockhash:     "1",
						Size:              1,
						// the same amount of calls to client.Transaction have to be passed to mocker, represented below
						Txs: []string{"1", "2", "3", "4"},
					},
				}
				m.EXPECT().BlockHash("btc", "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876").Return(&block, nil)

				m.EXPECT().Transaction("btc", "1").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "1",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)

				m.EXPECT().Transaction("btc", "2").Return(&sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "2",
						Fee:       "1",
						SentValue: "1",
						Time:      unixTime,
					},
				}, nil)

				m.EXPECT().Transaction("btc", "3").Return(nil, errors.New("some"))
				m.EXPECT().Transaction("btc", "4").Return(nil, sochain.NewClientErr(errors.New("some"), http.StatusInternalServerError))
			},

			want: &sochain.BlockResponse{
				Blocknumber:  1,
				Timestamp:    timeRFC3339,
				Size:         1,
				PreviousHash: "1",
				NextHash:     "1",
				Transactions: sochain.TransactionResponses{
					{
						TxID:      "1",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
					{
						TxID:      "2",
						Fee:       "1",
						Timestamp: timeRFC3339,
						Value:     "1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			path := "http://localhost:8080/network/:id"
			mCtrl := gomock.NewController(t)
			defer mCtrl.Finish()

			mockConn := mock_client.NewMockConnector(mCtrl)
			if tt.mock != nil {
				tt.mock(mockConn)
			}

			handler := func(w http.ResponseWriter, r *http.Request) *gin.Engine {
				gin.SetMode(gin.TestMode)
				c, ginEngine := gin.CreateTestContext(w)

				if tt.gotPathNetworkIDExists {
					c.Params = append(c.Params, gin.Param{Key: "id", Value: tt.gotPathNetworkID})
				}

				if tt.gotHeightQuery != "" {
					path = path + "?height=" + tt.gotHeightQuery
				}

				if tt.gotBlockhashQuery != "" {
					path = path + "?blockhash=" + tt.gotBlockhashQuery
				}

				var err error
				c.Request, err = http.NewRequest("GET", path, r.Body)
				c.Request.URL.RawPath = path
				assert.Nil(t, err)

				controller := NewController(zap.NewNop(), mockConn)
				controller.HandleGetBlock(c)

				return ginEngine
			}

			request := httptest.NewRequest("GET", path, nil)
			httpRecorder := httptest.NewRecorder()

			handler(httpRecorder, request)

			data, err := ioutil.ReadAll(httpRecorder.Body)
			assert.Nil(t, err)

			assert.Equal(t, tt.wantCode, httpRecorder.Code)

			if !tt.wantError {
				var response sochain.BlockResponse
				assert.Nil(t, json.Unmarshal(data, &response))
				assert.True(t, reflect.DeepEqual(response, *tt.want))
			}
		})
	}
}

func TestHandleGetTransaction(t *testing.T) {

	unixTime := 1231455600
	timeRFC3339 := time.Unix(int64(unixTime), 0).Format(time.RFC3339)

	tests := []struct {
		title                  string
		gotPathNetworkIDExists bool
		gotPathTxHashExists    bool
		gotPathNetworkID       string
		gotPathTxHash          string
		gotHeightQuery         string
		gotBlockhashQuery      string
		mock                   func(m *mock_client.MockConnector)
		wantError              bool
		wantCode               int
		want                   *sochain.TransactionResponse
	}{
		{
			title:                  "Error: path param networkID missing",
			wantError:              true,
			gotPathNetworkIDExists: false,
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: path param networkID invalid",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "eth",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: path param txhash missing",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotPathTxHashExists:    false,
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: txhash to long",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotPathTxHashExists:    true,
			gotPathTxHash:          "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876a",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: txhash to short",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotPathTxHashExists:    true,
			gotPathTxHash:          "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa87",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: txhash invalid characters",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotPathTxHashExists:    true,
			gotPathTxHash:          "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876&",
			wantCode:               http.StatusBadRequest,
		},
		{
			title:                  "Error: tx not found",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotPathTxHashExists:    true,
			gotPathTxHash:          "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
			wantCode:               http.StatusNotFound,
			mock: func(m *mock_client.MockConnector) {
				m.EXPECT().Transaction("btc", "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876").Return(nil, *sochain.NewClientErr(errors.New("some"), http.StatusNotFound))
			},
		},
		{
			title:                  "Error: bad request",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotPathTxHashExists:    true,
			gotPathTxHash:          "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
			wantCode:               http.StatusBadRequest,
			mock: func(m *mock_client.MockConnector) {
				m.EXPECT().Transaction("btc", "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876").Return(nil, *sochain.NewClientErr(errors.New("some"), http.StatusBadRequest))
			},
		},
		{
			title:                  "Error: some error ",
			wantError:              true,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotPathTxHashExists:    true,
			gotPathTxHash:          "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
			wantCode:               http.StatusInternalServerError,
			mock: func(m *mock_client.MockConnector) {
				m.EXPECT().Transaction("btc", "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876").Return(nil, errors.New("some"))
			},
		},
		{
			title:                  "Success",
			wantError:              false,
			gotPathNetworkIDExists: true,
			gotPathNetworkID:       "btc",
			gotPathTxHashExists:    true,
			gotPathTxHash:          "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
			wantCode:               http.StatusOK,
			mock: func(m *mock_client.MockConnector) {

				tx := sochain.Transaction{
					Data: sochain.TransactionData{
						Txid:      "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
						Time:      unixTime,
						Fee:       "1",
						SentValue: "1",
					},
				}

				m.EXPECT().Transaction("btc", "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876").Return(&tx, nil)
			},
			want: &sochain.TransactionResponse{
				TxID:      "eb6f76a4390f4e3cdcf8d2a73fc99d401965abca0372f350bf9317a34a1aa876",
				Timestamp: timeRFC3339,
				Fee:       "1",
				Value:     "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			path := "http://localhost:8080/network/:id/tx/:txhash"
			mCtrl := gomock.NewController(t)
			defer mCtrl.Finish()

			mockConn := mock_client.NewMockConnector(mCtrl)
			if tt.mock != nil {
				tt.mock(mockConn)
			}

			handler := func(w http.ResponseWriter, r *http.Request) *gin.Engine {
				gin.SetMode(gin.TestMode)
				c, ginEngine := gin.CreateTestContext(w)

				if tt.gotPathNetworkIDExists {
					c.Params = append(c.Params, gin.Param{Key: "id", Value: tt.gotPathNetworkID})
				}

				if tt.gotPathTxHashExists {
					c.Params = append(c.Params, gin.Param{Key: "txhash", Value: tt.gotPathTxHash})
				}

				if tt.gotHeightQuery != "" {
					path = path + "?height=" + tt.gotHeightQuery
				}

				if tt.gotBlockhashQuery != "" {
					path = path + "?blockhash=" + tt.gotBlockhashQuery
				}

				var err error
				c.Request, err = http.NewRequest("GET", path, r.Body)
				c.Request.URL.RawPath = path
				assert.Nil(t, err)

				controller := NewController(zap.NewNop(), mockConn)
				controller.HandleGetTransaction(c)

				return ginEngine
			}

			request := httptest.NewRequest("GET", path, nil)
			httpRecorder := httptest.NewRecorder()

			handler(httpRecorder, request)

			data, err := ioutil.ReadAll(httpRecorder.Body)
			assert.Nil(t, err)

			assert.Equal(t, tt.wantCode, httpRecorder.Code)

			if !tt.wantError {
				var response sochain.TransactionResponse
				assert.Nil(t, json.Unmarshal(data, &response))
				assert.True(t, reflect.DeepEqual(response, *tt.want))
			}
		})
	}
}
