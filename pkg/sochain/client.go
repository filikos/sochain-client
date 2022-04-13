package sochain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var apiURL = "https://sochain.com/api/v2"

type Sochain struct {
	Client  *http.Client
	baseUrl string
}

func NewSochain() Connector {
	return &Sochain{
		Client:  &http.Client{},
		baseUrl: apiURL,
	}
}

type Connector interface {
	NetworkInfo(networkID string) (*NetworkInfo, error)
	BlockHeight(networkID string, height int) (*Block, error)
	BlockHash(networkID, blockHash string) (*Block, error)
	Transaction(networkID, txHash string) (*Transaction, error)
}

func (c *Sochain) NetworkInfo(networkID string) (*NetworkInfo, error) {

	url := fmt.Sprintf("%s/get_info/%s", apiURL, networkID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewClientErr(fmt.Errorf("sochain response statuscode %d, networkID '%s'", resp.StatusCode, networkID), resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info NetworkInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *Sochain) BlockHeight(networkID string, height int) (*Block, error) {

	url := fmt.Sprintf("%s/get_block/%s/%d", apiURL, networkID, height)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewClientErr(fmt.Errorf("sochain response statuscode %d, height '%d'", resp.StatusCode, height), resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var b Block
	if err := json.Unmarshal(body, &b); err != nil {
		return nil, err
	}

	return &b, nil
}

func (c *Sochain) BlockHash(networkID, blockHash string) (*Block, error) {
	url := fmt.Sprintf("%s/get_block/%s/%s", apiURL, networkID, blockHash)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewClientErr(fmt.Errorf("sochain response statuscode %d, blockhash '%s'", resp.StatusCode, blockHash), resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var b Block
	if err := json.Unmarshal(body, &b); err != nil {
		return nil, err
	}

	return &b, nil
}

func (c *Sochain) Transaction(networkID, txHash string) (*Transaction, error) {

	url := fmt.Sprintf("%s/tx/%s/%s", apiURL, networkID, txHash)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewClientErr(fmt.Errorf("sochain response statuscode %d, txhash '%s'", resp.StatusCode, txHash), resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tx Transaction
	if err := json.Unmarshal(body, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}
