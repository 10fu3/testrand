## testrand
### 概要
This program is interpretter, that is interpretation S-Expression, and can distribute computing.

これはS式を解釈し、計算を分散させることができるインタプリタです。

This main.go file is a simulation that multi-machine movement.

main.goファイルは複数マシンでの分散挙動をシミュレーションするものです。

This client.go file is the role of lightweight computing to run and heavyClient.go file is the role of heavyweight computing to run.

client.goは軽量な計算を行うために起動させるファイルで、heavyClient.goは重たい計算をさせるために起動させるファイルです。

Light clients and Heavy clients do messaging for calculating by JSON and HTTP servers.

軽量計算用マシンと重量計算用マシンはJSONとHTTPサーバーによって計算のためのメッセージングをします。

When sending S-Expression, this system is serialize S-Expression to String and writes on JSON, and sends to another system HTTP endpoint.

もしS式を送る場合、文字列にS式をシリアライズし、JSONにパックして相手のHTTPエンドポイントに送り付けます
### 対応構文
```
"and"
"or"
"if" 
"eq?"
"not"
"define"
"set"
"loop"
"wait"
"+"
"begin"
"lambda"
"quote"
```
