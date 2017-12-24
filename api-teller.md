统一说明
返回json object结构统一为：
成功：{"code":0, "data":object}
失败：{"code":1,"errmsg":"errmsg","data":object}

1.1 /api/all-cryptocurrency-type/
无参数
返回json object, data属性为对象数组（即map）, key为加密数字货币简称，value为完整名称

1.2 /api/generate-address/
使用application/x-www-form-urlencoded方式提交数据，两个参数，currencyType是数字货币简称，size是生成地址数量
返回json object, data为字符串数组

1.3 /api/query-address-balance/
使用application/json方式提交数据，提交json object两个属性，currencyType是数字货币简称，address是字符串数组
返回json object, data为对象，结构参考以下golang结构(未来可能微调)：
type AddressBalance struct{
	Address string
	Balance float64
	LastUpdatedTimestamp uint64
	LastUpdatedTransactionId string		
}

1.4 /api/send-coin/
使用application/json方式提交数据，提交json object为以下golang结构数组：

```bash
Method: POST
Accept: application/json
Content-Type: application/json
URI: /api/send-coin
Request Body: [{
    "address": "2AzuN3aqF53vUC2yHqdfMKnw4i8eRrwye71"
    "amount":"1234"}
    ]
```
Example:

```bash
curl -X POST -H "Content-Type:application/json" -d '[{"address":"2AzuN3aqF53vUC2yHqdfMKnw4i8eRrwye71","amount":"1234"}]' http://localhost:7071/api/send-coin
```

response:

```bash
{
    "errmsg": "",
    "code": 0,
    "data": "done"
}

## Teller Service Api

### Bind

```bash
Method: POST
Accept: application/json
Content-Type: application/json
URI: /api/bind
Request Body: {
    "address": "..."
    "tokenType":"[bitcoin|ethcoin|skycoin|xmrcoin]"
}
```

Binds a spocoin address to a deposit address. A spocoin address can be bound to
multiple deposit addresses.  The default maximum number of bound addresses is 5.

Example:

```bash
curl -X POST -H "Content-Type:application/json" -d '{"address":"2AzuN3aqF53vUC2yHqdfMKnw4i8eRrwye71","coin_type":"skycoin"}' http://localhost:7071/api/bind/
```

response:

```bash
{
    "errmsg": "",
    "code": 0,
    "data": {
        "address": "2do3K1YLMy3Aq6EcPMdncEurP5BfAUdFPJj",
        "coin_type": "skycoin"
    }
}
```

### Status

```bash
Method: GET
Content-Type: application/json
URI: /api/status
Query Args: address
```

Returns statuses of a spocoin address.

Possible statuses are:

* `waiting_deposit` - Spocoin address is bound, no deposit seen on address yet
* `waiting_send` - deposit detected, waiting to send spocoin out
* `waiting_confirm` - Spocoin sent out, waiting to confirm the spocoin transaction
* `done` - Skycoin transaction confirmed

Example:

```bash
curl http://localhost:7071/api/status?address=2AzuN3aqF53vUC2yHqdfMKnw4i8eRrwye71\&coin_type=bitcoin
```

response:

```bash
{
    "statuses": [
        {
            "seq": 1,
            "update_at": 1501137828,
            "address":"ZJHwZfwXrqq49bEKmXXCqjcMTzF8RkQSXm",
            "tokenType":"bitcoin"
            "status": "done"
        },
        {
            "seq": 2,
            "update_at": 1501128062,
            "address":"ZJHwZfwXrqq49bEKmXXCqjcMTzF8RkQSXm",
            "tokenType":"bitcoin"
            "status": "waiting_deposit"
        },
        {
            "seq": 3,
            "update_at": 1501128063,
            "address":"ZJHwZfwXrqq49bEKmXXCqjcMTzF8RkQSXm",
            "tokenType":"bitcoin"
            "status": "waiting_deposit"
        },
    ]
}
```

### Rate

```bash
Method: GET
Content-Type: application/json
URI: /api/rate
Query Args: tokenType=[all|skycoin|bitcoin|ethcoin|xmrcoin]
```

Returns exchange rate

Example:

```bash
curl "http://127.0.0.1:7071/api/rate"
```

Response:
```
{
    "errmsg": "",
    "code": 0,
    "data": {
        "tokenType": "all",
        "rate": 0,
        "allcoin": [
            {
                "coin_name": "bitcoin",
                "coin_code": "BTC",
                "coin_rate": 22200
            },
            {
                "coin_name": "ethcoin",
                "coin_code": "ETH",
                "coin_rate": 1200
            },
            {
                "coin_name": "skycoin",
                "coin_code": "SKY",
                "coin_rate": 200
            },
            {
                "coin_name": "xmrcoin",
                "coin_code": "XMR",
                "coin_rate": 300
            }
        ]
    }
}
```
----

新的接口如下：

```bash
Method: GET
Content-Type: application/json
URI: /api/deposite
Query Args: req
```
如果nextseq 不为空， 则继续请求，新请求的req=nextseq， 如果为空，说明没有新的存入，等待间隔后继续请求
如果req指定为0，则从第一个存入返回
Response:
```
{"code":0, "data":{nextseq:5, deposit:[depositValue1, depositValue2]}
depositValue:
{"Seq":3,"UpdatedAt":1513210524, "SpoAddress":"6v7gu8WP2V9aggo","depositAddress":"5fa2f213f18690bc","CoinType":"bitcoin", "Txid":"3486ca63d6169536c4552bm "SkySent":12000000,"SkyBtcRate":25, "depositValue":0.48,"Height":105948}
```


--------以下是golang接口版本-----
//AllCryptocurrencyType give all Cryptocurrency Type, return map key is short name, map value is full name
func AllCryptocurrencyType() map[string]string

//GenerateAddress generate a batch of digital currency address, currencyType is BTC, ETH etc, size is batch size
func GenerateAddress(currencyType String, size uint32) []string

//AddressBalance BTC, ETH etc account balance and last updated timestamp and transaction id
type AddressBalance struct{
	Address string
	Balance float64
	LastUpdatedTimestamp uint64
	LastUpdatedTransactionId string		
}

//QueryAddressBalance get BTC, ETH etc account balance
func QueryAddressBalance(currencyType String, address ...string) []AddressBalance

type SendCoinInfo struct{
	Address string
	Amount uint64
}

//SendCoin transfer coin and reward to address
func SendCoin(addrAndAmount []SendCoinInfo)
