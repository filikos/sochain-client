package sochain

import "time"

type NetworkInfos []NetworkInfo
type NetworkInfo struct {
	Status string      `json:"status"`
	Data   NetworkData `json:"data"`
}

type NetworkData struct {
	Name             string `json:"name"`
	Acronym          string `json:"acronym"`
	Network          string `json:"network"`
	SymbolHtmlcode   string `json:"symbol_htmlcode"`
	URL              string `json:"url"`
	MiningDifficulty string `json:"mining_difficulty"`
	UnconfirmedTxs   int    `json:"unconfirmed_txs"`
	Blocks           int    `json:"blocks"`
	Price            string `json:"price"`
	PriceBase        string `json:"price_base"`
	PriceUpdateTime  int    `json:"price_update_time"`
	Hashrate         string `json:"hashrate"`
}

type Blocks []Block
type Block struct {
	Status string    `json:"status"`
	Data   BlockData `json:"data"`
}

type BlockData struct {
	Network           string   `json:"network"`
	Blockhash         string   `json:"blockhash"`
	BlockNo           int      `json:"block_no"`
	MiningDifficulty  string   `json:"mining_difficulty"`
	Time              int      `json:"time"`
	Confirmations     int      `json:"confirmations"`
	IsOrphan          bool     `json:"is_orphan"`
	Txs               []string `json:"txs"`
	Merkleroot        string   `json:"merkleroot"`
	PreviousBlockhash string   `json:"previous_blockhash"`
	NextBlockhash     string   `json:"next_blockhash"`
	Size              int      `json:"size"`
}

type BlockResponse struct {
	Blocknumber  int                  `json:"blocknumber"`
	Timestamp    string               `json:"timestamp"`
	PreviousHash string               `json:"previoushash"`
	NextHash     string               `json:"nexthash"`
	Size         int                  `json:"size"`
	Transactions TransactionResponses `json:"transactions"`
}

func (b *Block) Response() BlockResponse {
	return BlockResponse{
		Blocknumber:  b.Data.BlockNo,
		Timestamp:    time.Unix(int64(b.Data.Time), 0).Format(time.RFC3339),
		PreviousHash: b.Data.PreviousBlockhash,
		NextHash:     b.Data.NextBlockhash,
		Size:         b.Data.Size,
	}
}

type Transactions []Transaction
type Transaction struct {
	Status  string          `json:"status"`
	Data    TransactionData `json:"data"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
}

type TransactionData struct {
	Network       string  `json:"network"`
	Txid          string  `json:"txid"`
	Blockhash     string  `json:"blockhash"`
	BlockNo       int     `json:"block_no"`
	Confirmations int     `json:"confirmations"`
	Time          int     `json:"time"`
	Size          int     `json:"size"`
	Vsize         int     `json:"vsize"`
	Version       int     `json:"version"`
	Locktime      int     `json:"locktime"`
	SentValue     string  `json:"sent_value"`
	Fee           string  `json:"fee"`
	Inputs        Inputs  `json:"inputs"`
	Outputs       Outputs `json:"outputs"`
	TxHex         string  `json:"tx_hex"`
}

type TransactionResponses []TransactionResponse

func (t Transactions) Response() TransactionResponses {
	r := make(TransactionResponses, len(t))
	for i := 0; i < len(t); i++ {
		r[i] = t[i].Response()
	}

	return r
}

func (t Transaction) Response() TransactionResponse {
	return TransactionResponse{
		TxID:      t.Data.Txid,
		Timestamp: time.Unix(int64(t.Data.Time), 0).Format(time.RFC3339),
		Fee:       t.Data.Fee,
		Value:     t.Data.SentValue,
	}
}

type TransactionResponse struct {
	TxID      string `json:"txid,omitempty"`
	Timestamp string `json:"time,omitempty"`
	Fee       string `json:"fee,omitempty"`
	Value     string `json:"sent_value,omitempty"`
}

type Inputs []Input
type Input struct {
	InputNo      int         `json:"input_no"`
	Address      string      `json:"address"`
	Value        string      `json:"value"`
	ReceivedFrom interface{} `json:"received_from"`
	ScriptAsm    string      `json:"script_asm"`
	ScriptHex    interface{} `json:"script_hex"`
	Witness      []string    `json:"witness"`
}

type Outputs []Output
type Output struct {
	OutputNo  int         `json:"output_no"`
	Address   string      `json:"address"`
	Value     string      `json:"value"`
	Type      string      `json:"type"`
	ReqSigs   interface{} `json:"req_sigs"`
	Spent     interface{} `json:"spent"`
	ScriptAsm string      `json:"script_asm"`
	ScriptHex string      `json:"script_hex"`
}
