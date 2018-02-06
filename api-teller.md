## Unified Description

Response format is json as below

```
success {"code":0, "data":object}
failed {"code":1,"errmsg":"errmsg","data":object}

```
Request Validation

```

there are header `timestamp` and `auth` in request，auth value is sha256(timestamp+secretToken)
first step: timestamp must be within specified period of validity
second step: auth must pass validation

go code as :

timestamp := fmt.Sprintf("%d", time.Now().UTC().Unix())
hash := hmac.New(sha256.New, []byte(token))
hash.Write([]byte(timestamp))
req.Header.Set("Content-Type", "application/json")
req.Header.Set("timestamp", timestamp)
req.Header.Set("auth", hex.EncodeToString(hash.Sum(nil)))

```

1.1 /api/config
args: none
```
{
    "enabled": true,
    "max_bound_addrs": 7,
    "max_decimals": 3,
    "allcoins": {
        "BTC": {
            "coin_name": "BTC",
            "rate": "200",
            "enabled": false,
            "unit": 100000000,
            "confirmations_required": 1
        },
        "ETH": {
            "coin_name": "ETH",
            "rate": "100",
            "enabled": false,
            "unit": 1000000000,
            "confirmations_required": 1
        },
        "SKY": {
            "coin_name": "SKY",
            "rate": "70",
            "enabled": true,
            "unit": 1000000,
            "confirmations_required": 0
        },
        "XMR": {
            "coin_name": "XMR",
            "rate": "10",
            "enabled": false,
            "unit": 1000000000000,
            "confirmations_required": 1
        }
    }
}
```
1.4 /api/depoist/

```bash
Method: POST
Accept: application/json
Content-Type: application/json
URI: /api/deposit
Request Body:
[{"seq":11,"update_at":1517134128,"coin_type":"SKY","address":"25PN9qx8NKga2RFqHNv5xm9UBuowk5gi9pv","deposit_address":"LAWbVXeTL82vxjh21TNv6ALnMv2CT1mjL4","txid":"f37d9e96b84c5a7451993e5252da91d84c857a406caa4e22eb783e21eb8907a8","rate":"70","sent":119000000,"deposit_value":1700000,"height":12643}]
```
Example:

```bash
curl -X POST -H "Content-Type:application/json" -d '[deposit]' http://affiliate:port/api/deposit/
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
HEADER: affilaite=true
Request Body: {
    "address": "..."
    "tokenType":"[BTC|ETH|SKY|XMR]"
}
```

Binds a spocoin address to a deposit address. A spocoin address can be bound to
multiple deposit addresses.  The default maximum number of bound addresses is 5.

Example:

```bash
curl -X POST -H "Content-Type:application/json" -H "affiliate:true" -d '{"address":"2AzuN3aqF53vUC2yHqdfMKnw4i8eRrwye71","tokenType":"SKY"}' http://localhost:7071/api/bind/
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
    "errmsg": "",
    "code": 0,
    "data": {
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

```
----

### get reward

```bash
Method: GET
URI: /api/reward/
```
Response:
```
[{Addr:nF5xC41ZBh7vXQqMLSnLPRc68jDXMy9GL6 Coins:1000000 ID:32} {Addr:k9QgadMDxisLfj2CLgNrwgzZZSEmZMeUpK Coins:1000000 ID:33}]
```

### post reward-status

```bash
Method: POST
Content-Type: application/json
URI: /api/reward-status/
ARGS:
    [10, 30]
```

