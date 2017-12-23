# affiliate
affiliate code system for spaco

```
1. teller 发送对应的币给用户
2. affiliate 通过接口获取已经打币的交易信息
3. affiliate 计算出奖励，发送命令给teller，teller按照指定发送奖励
```

总体原则：
teller负责打币
affiliate负责计算奖励
affilicate 通过接口获知 有哪些地址已经存入，并且被打币了。
