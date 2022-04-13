package controller

import (
	"log"
	"net/http"
	"sochain-client/pkg/sochain"
	"sochain-client/pkg/util"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
	client sochain.Connector
}

func NewController(l *zap.Logger, client sochain.Connector) *Controller {
	return &Controller{
		logger: l,
		client: client,
	}
}

const maxTxPerBlock = 10

// BTC, LTC & DOGE use SHA-256 for blocks & tx hashes
var HashSHA256Regex *regexp.Regexp

func init() {
	var err error
	HashSHA256Regex, err = regexp.Compile("^[A-Fa-f0-9]{64}$")
	if err != nil {
		log.Fatal(err)
	}
}

// Rturns latest block of network including transactions. Specific block can be choosen optional by providing blockcounter or blockhash
func (c *Controller) HandleGetBlock(ctx *gin.Context) {
	networkID, err := util.GetParamNetwork(ctx, "id")
	if err != nil {
		c.logger.Info("invalid path param network 'id'", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, "path param: network 'id' can only be 'btc', 'eth' or 'ltc'")
		return
	}

	networkID = strings.ToLower(networkID)

	var wg sync.WaitGroup
	results := make(chan TxChanResp)

	height := ctx.Query("height")
	blockHash := ctx.Query("blockhash")
	if height == "" && blockHash == "" {
		info, err := c.client.NetworkInfo(networkID)
		if err != nil {
			c.logger.Info("missing query params 'height' or 'blockhash'", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, "one of query params 'height' or 'blockhash' is mandatory")
			return
		}

		block, err := c.client.BlockHeight(strings.ToLower(networkID), info.Data.Blocks)
		if err != nil {
			if cErr, ok := err.(sochain.ClientError); ok {
				switch cErr.Code() {
				case http.StatusNotFound:
					c.logger.Info("unable to fetch transaction", zap.Error(cErr))
					ctx.JSON(http.StatusNotFound, "unable to find block for given height")
					return
				case http.StatusBadRequest:
					c.logger.Info("unable to fetch transaction", zap.Error(cErr))
					ctx.JSON(http.StatusBadRequest, "bad request")
					return
				}
			}

			c.logger.Info("unable to receive block by height", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, "unable to fetch block for given height")
			return
		}

		f := func(networkID, txHash string) (*sochain.Transaction, error) {
			return c.client.Transaction(networkID, txHash)
		}

		if len(block.Data.Txs) >= maxTxPerBlock {
			FetchTxAsync(f, &wg, results, networkID, block.Data.Txs[:maxTxPerBlock])
		} else {
			FetchTxAsync(f, &wg, results, networkID, block.Data.Txs)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		transactions := make(sochain.Transactions, 0)
		for v := range results {
			if v.err != nil {
				if cErr, ok := err.(sochain.ClientError); ok {
					c.logger.Warn("unable to fetch transaction", zap.String("txhash", v.hash), zap.Int("statuscode", cErr.Code()), zap.Error(cErr))
					continue
				}

				c.logger.Info("unable to fetch transaction", zap.String("txhash", v.hash), zap.Error(v.err))
			} else {
				transactions = append(transactions, *v.tx)
			}
		}

		bResp := block.Response()
		bResp.Transactions = transactions.Response()

		ctx.JSON(http.StatusOK, bResp)
		return
	}

	if height != "" {
		heightInt, err := strconv.Atoi(height)
		if err != nil || heightInt <= 0 {
			c.logger.Info("invalid query param 'height'", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, "query param 'height' is no integer number, zero or negative")
			return
		}

		block, err := c.client.BlockHeight(strings.ToLower(networkID), heightInt)
		if err != nil {
			if cErr, ok := err.(sochain.ClientError); ok {
				switch cErr.Code() {
				case http.StatusNotFound:
					c.logger.Info("unable to fetch transaction", zap.Error(cErr))
					ctx.JSON(http.StatusNotFound, "unable to find block for given height")
					return
				case http.StatusBadRequest:
					c.logger.Info("unable to fetch transaction", zap.Error(cErr))
					ctx.JSON(http.StatusBadRequest, "bad request")
					return
				}
			}

			c.logger.Info("unable to receive block by height", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, "unable to fetch block for given height")
			return
		}

		f := func(networkID, txHash string) (*sochain.Transaction, error) {
			return c.client.Transaction(networkID, txHash)
		}

		if len(block.Data.Txs) >= maxTxPerBlock {
			FetchTxAsync(f, &wg, results, networkID, block.Data.Txs[:maxTxPerBlock])
		} else {
			FetchTxAsync(f, &wg, results, networkID, block.Data.Txs)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		transactions := make(sochain.Transactions, 0)
		for v := range results {
			if v.err != nil {
				if cErr, ok := err.(sochain.ClientError); ok {
					c.logger.Warn("unable to fetch transaction", zap.String("txhash", v.hash), zap.Int("statuscode", cErr.Code()), zap.Error(cErr))
					continue
				}

				c.logger.Info("unable to fetch transaction", zap.String("txhash", v.hash), zap.Error(v.err))
			} else {
				transactions = append(transactions, *v.tx)
			}
		}

		bResp := block.Response()
		bResp.Transactions = transactions.Response()

		ctx.JSON(http.StatusOK, bResp)
		return
	}

	if blockHash != "" {
		if !HashSHA256Regex.MatchString(blockHash) {
			c.logger.Info("provided blockhash is not a valid SHA-256 hash", zap.String("blockhash", blockHash))
			ctx.JSON(http.StatusBadRequest, "provided blockhash is not a valid SHA-256 hash")
			return
		}

		block, err := c.client.BlockHash(networkID, blockHash)
		if err != nil {
			if cErr, ok := err.(sochain.ClientError); ok {
				switch cErr.Code() {
				case http.StatusNotFound:
					c.logger.Info("unable to fetch transaction", zap.Error(cErr))
					ctx.JSON(http.StatusNotFound, "unable to find block for given hash")
					return
				case http.StatusBadRequest:
					c.logger.Info("unable to fetch transaction", zap.Error(cErr))
					ctx.JSON(http.StatusBadRequest, "bad request block for given hash")
					return
				}
			}

			c.logger.Info("unable to fetch block by hash", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, "unable to fetch block by hash")
			return
		}

		f := func(networkID, txHash string) (*sochain.Transaction, error) {
			return c.client.Transaction(networkID, txHash)
		}

		if len(block.Data.Txs) >= maxTxPerBlock {
			FetchTxAsync(f, &wg, results, networkID, block.Data.Txs[:maxTxPerBlock])
		} else {
			FetchTxAsync(f, &wg, results, networkID, block.Data.Txs)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		transactions := make(sochain.Transactions, 0)
		for v := range results {
			if v.err != nil {
				if cErr, ok := err.(sochain.ClientError); ok {
					c.logger.Warn("unable to fetch transaction", zap.String("txhash", v.hash), zap.Int("statuscode", cErr.Code()), zap.Error(cErr))
					continue
				}

				c.logger.Warn("unable to fetch transaction", zap.String("txhash", v.hash), zap.Error(v.err))
			} else {
				transactions = append(transactions, *v.tx)
			}
		}

		bResp := block.Response()
		bResp.Transactions = transactions.Response()

		ctx.JSON(http.StatusOK, bResp)
		return
	}
}

type TxChanResp struct {
	tx   *sochain.Transaction
	hash string
	err  error
}

func FetchTxAsync(f func(networkID, blockHash string) (*sochain.Transaction, error),
	wg *sync.WaitGroup, ch chan TxChanResp, networkID string, txHashList []string) {

	for _, v := range txHashList {
		wg.Add(1)

		go func(networkID, hash string) {
			defer wg.Done()
			tx, err := f(networkID, hash)
			resp := TxChanResp{
				tx:  tx,
				err: err,
			}
			ch <- resp
		}(networkID, v)
	}
}

// Returns details of specific transaction
func (c *Controller) HandleGetTransaction(ctx *gin.Context) {
	networkID, err := util.GetParamNetwork(ctx, "id")
	if err != nil {
		c.logger.Info("invalid path param network'id'", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, "path param: Network 'id' can only be 'btc', 'eth' or 'ltc'")
		return
	}

	txHash := ctx.Param("txhash")
	if txHash == "" {
		c.logger.Info("path param 'txhash' missing")
		ctx.JSON(http.StatusBadRequest, "path param: missing 'txhash'")
		return
	}

	if !HashSHA256Regex.MatchString(txHash) {
		c.logger.Info("path param: 'txhash' is not a valid SHA-256 hash")
		ctx.JSON(http.StatusBadRequest, "path param: 'txhash' is not a valid SHA-256 hash")
		return
	}

	tx, err := c.client.Transaction(networkID, txHash)
	if err != nil {
		if cErr, ok := err.(sochain.ClientError); ok {
			switch cErr.Code() {
			case http.StatusNotFound:
				c.logger.Info("unable to fetch transaction", zap.Error(cErr))
				ctx.JSON(http.StatusNotFound, "unable to find tx for given hash")
				return
			case http.StatusBadRequest:
				c.logger.Info("unable to fetch transaction", zap.Error(cErr))
				ctx.JSON(http.StatusBadRequest, "bad request tx for given hash")
				return
			}
		}

		c.logger.Info("unable to fetch transaction", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, "unable to fetch transaction")
		return
	}

	ctx.JSON(http.StatusOK, tx.Response())
}
