
统一说明 返回json object结构统一为： 成功：{"code":0, "data":object} 失败：{"code":1,"errmsg":"errmsg","data":object}


1. /api/deposit/
发送存入信息
Method: POST
Content-Type: application/json
Request Body:
[{"seq":3,"updated_at":1513210524, "spo_address":"6v7gu8WP2V9aggo","deposit_address":"5fa2f213f18690bc","coin_type":"bitcoin", "txid":"3486ca63d6169536c4552bm "sent":12000000,"rate":25, "deposit_value":0.48,"height":105948},{"seq":4,"updated_at":1513220524, "spo_address":"6v7gu8WP2V9aggp","deposit_address":"5fa2f213f18690bd","coin_type":"bitcoin", "txid":"3486ca63d6169536c4552bn "sent":2000000,"rate":25, "deposit_value":0.48,"height":105949}]
Response:
成功：{"code":0}
失败：{"code":1,"errmsg":"errmsg"}
要么整体全部失败，要么整体全部成功

2. /api/reward/
获取奖励信息
Method: GET
Request Body:
成功：{"code":0,"data": [{"id":29,"address":"nF5xC41ZBh7vXQqMLSnLPRc68jDXMy9GL6","amount":1000000},{"id":30,"address":"k9QgadMDxisLfj2CLgNrwgzZZSEmZMeUpK","amount":1000000}]}
失败：{"code":1,"errmsg":"errmsg"}

3. /api/reward-status/
更新奖励完成状态，接收到的id将都被标记为已打币状态
Method: POST
Content-Type: application/json
Request Body:
[29,30]
Response:
成功：{"code":0}
失败：{"code":1,"errmsg":"errmsg"}
要么整体全部失败，要么整体全部成功
