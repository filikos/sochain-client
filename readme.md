# Sochain Client

This Application provides a simple-to-use service which wraps the https://sochain.com/ API and serves Blockchain related Information about Blocks and Transactions from the Bitcoin, Litecoin & Dogecoin Network.

## Requirements
- Go >= v1.17
- Mockgen >= v1.6.0 (https://github.com/golang/mock)

## Quick Start

#### Create your .env config file
```bash
cp .env.example .env
```

##### CI_ENV Values (optional)
Logging from Debug level:
'development'
'staging'

Logging from Info level:
'production'

Log levels: (https://github.com/uber-go/zap/blob/master/level.go)

##### Start Application
```bash
make run
```

##### Run all Tests
```bash
make tests
```

##### Lint
```bash
make lint
```

## Endpoints

<details><summary>GET /network/{id} </summary>
<p>

### Description:

Returns the latest block of choosen network {id} including the last 10 transactions.

Optional a specific block can be fetched by providing either the **height** of a block or its **blockhash**.
NOTE: Timestamps are formatted in **RFC3339** for increased readability, unification & timezone informations. (https://datatracker.ietf.org/doc/html/rfc3339)

### Parameters:
Content-Type: **application/json**

**Path Param:**
*required*
Name: *id*
Type: string
Values: 'BTC', 'LTC', 'DOGE'

**Query:**
*optional*
Name: *blockhash*
Type: string
Desc: Has to be a valid SHA-256 blockhash of the corresponding network.
Example BTC Blockhash: "00000000000000000008fa3759141044ae3db1e6ec222e114651354f58d5cc42"

**Query:**
*optional*
Name: *height*
Type: int
Desc: Has to be a valid block height of the corresponding network.
Example: 729446


### Request example
curl --location --request GET 'http://localhost:8080/network/btc'

### Example Response Body:

```json
{
    "blocknumber": 729576,
    "timestamp": "2022-03-29T18:00:12+02:00",
    "previoushash": "00000000000000000002468013524b804a49edc02e2100772d046f010006699c",
    "nexthash": "",
    "size": 1385079,
    "transactions": [
        {
            "txid": "b09201c3df876de5e785ed8cec6b6ef83e9f00228959ecb015d3a0dfc48edf08",
            "time": "2022-03-29T18:00:12+02:00",
            "fee": "0.0",
            "sent_value": "6.32374561"
        },
        {
            "txid": "6f86fe618a57bcc98e89270b41a4bb9aa8ca32614eab768b153d558fecacf01f",
            "time": "2022-03-29T18:00:12+02:00",
            "fee": "0.00047750",
            "sent_value": "0.06531446"
        },
        {
            "txid": "1e1784fa064b42377c5c4a31437a162125b8fcfcd6beb08173b66a7346163fbe",
            "time": "2022-03-29T18:00:12+02:00",
            "fee": "0.00040200",
            "sent_value": "1.39669721"
        },
        {
            "txid": "5d8c46f7ff3957e20d4bc5359b03f87b9bb5e7cd032c17bbf4cf6f9043453fa6",
            "time": "2022-03-29T18:00:12+02:00",
            "fee": "0.00045000",
            "sent_value": "0.00345129"
        },
        {
            "txid": "87d67b541dea2d719e95402dc15e821fdd117c7c2a491b4572dddf54519cc321",
            "time": "2022-03-29T18:00:12+02:00",
            "fee": "0.00062400",
            "sent_value": "0.02982970"
        },
        {
            "txid": "346e30ac98e72bb7972e749e5b39a0b26200565246051ac170c98ced7b233f65",
            "time": "2022-03-29T18:00:12+02:00",
            "fee": "0.00045000",
            "sent_value": "0.02270601"
        },
        {
            "txid": "dacbe54480f426e0a92a774221849430061c64db9032c86b28aa661fc1127d39",
            "time": "2022-03-29T18:00:12+02:00",
            "fee": "0.00028800",
            "sent_value": "1.97494441"
        },
        {
            "txid": "13cb8247709d4c20d61db19efae6dd34aec325a809cbbd360854b5724440726f",
            "time": "2022-03-29T18:00:12+02:00",
            "fee": "0.00042485",
            "sent_value": "0.51118484"
        },
        {
            "txid": "e2207191530658599384af48071ed827654835fa0218293360aa5488805f0d29",
            "time": "2022-03-29T18:00:12+02:00",
            "fee": "0.00100000",
            "sent_value": "0.97692110"
        },
        {
            "txid": "b9b314f1504842e90c1a1678aa35b559dc18bf90668aa48157c592d0c92db839",
            "time": "2022-03-29T18:00:12+02:00",
            "fee": "0.00116768",
            "sent_value": "0.25535031"
        }
    ]
}
```

### Responses:
200 OK
400 Bad Request
404 Not Found
500 Internal Server Error

</p>
</details>

<details><summary>GET /network/{id}/tx/{txhash} </summary>
<p>

### Description:

Returns a specific transaction.

NOTE: Timestamps are formatted in **RFC3339** for increased readability, unification & timezone informations. (https://datatracker.ietf.org/doc/html/rfc3339)

### Parameters:
Content-Type: **application/json**

**Path Param:**
*required*
Name: *id*
Type: string
Values: 'BTC', 'LTC', 'DOGE'

**Path Param:**
*required*
Name: *txhash*
Type: string
Desc: Has to be a valid SHA-256 blockhash of the corresponding network.
Example BTC transaction hash: "7496d0464cc324467f16bdec3db1a088a609c500fec6b9d123c0a22813f9983c"

### Request example
curl --location --request GET 'http://localhost:8080/network/btc/tx/2b068b203412a81666d8fc9e662eac81bca9cc881b354d5164039f571a078ddd'

### Example Response Body:

```json
{
    "txid": "2b068b203412a81666d8fc9e662eac81bca9cc881b354d5164039f571a078ddd",
    "time": "2022-03-29T12:23:39+02:00",
    "fee": "0.00000501",
    "sent_value": "0.10939511"
}
```

### Responses:
200 OK
400 Bad Request
404 Not Found
500 Internal Server Error

</p>
</details>