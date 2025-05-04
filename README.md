# stdio-logging-proxy

## 概要

`stdio-logging-proxy` は、指定されたコマンドの標準入力と標準出力をプロキシし、すべての通信内容をログファイルに記録します。これにより、MCP サーバーとの通信をデバッグする際に、送受信されたメッセージを簡単に確認できます。

## 機能

*   **標準入出力のプロキシ:** 指定されたコマンドの標準入力と標準出力をプロキシします。
*   **通信内容のログ記録:** すべての通信内容をログファイルに記録します。
*   **環境変数による設定:** ログファイルのパスやプロキシ対象のコマンドを環境変数で設定できます。

## 使い方

### インストール

```bash
go install github.com/usk6666/stdio-logging-proxy
```

### 実行例

```bash
export PROXY_COMMAND="/bin/bash"
stdio-logging-proxy
```

## 設定

### 環境変数

*   `LOG_FILE`: ログファイルのパスを指定します (デフォルト: `stdio-logging-proxy.log`)。
*   `PROXY_COMMAND`: プロキシするコマンドを指定します (デフォルト: `/bin/bash`)。

### Roo CodeへのMCPサーバ設定例

- .roo/mcp.json
```json
{
  "mcpServers": {
    "example": {
      "name": "example-mcp",
      "description": "example mcp server",
      "command": "/path/to/stdio-logging-proxy",
      "env": {
        "PROXY_COMMAND": "/path/to/example-mcp-server",
        "LOG_FILE": "/path/to/logfile.txt"
      },
      "timeout": 30,
      "alwaysAllow": [],
      "disabled": false
    }
  }
}
```

## ログ

ログファイルには、標準入力、標準出力、標準エラー出力の内容が記録されます。

- ログイメージ
```log
2025/05/04 04:17:00 main.go:42: Starting proxy...
2025/05/04 04:17:00 main.go:92: stdin: {"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"Roo Code","version":"3.15.3"}},"jsonrpc":"2.0","id":0}
2025/05/04 04:17:00 main.go:117: stdout: {"id":0,"jsonrpc":"2.0","result":{"capabilities":{"prompts":{"listChanged":false},"resources":{"listChanged":false},"tools":{"listChanged":false}},"protocolVersion":"2024-11-05","serverInfo":{"name":"","version":""}}}
2025/05/04 04:17:00 main.go:92: stdin: {"method":"notifications/initialized","jsonrpc":"2.0"}
2025/05/04 04:17:00 main.go:92: stdin: {"method":"tools/list","jsonrpc":"2.0","id":1}
2025/05/04 04:17:00 main.go:117: stdout: {"id":1,"jsonrpc":"2.0","result":{"tools":[{"description":"example mcp description","inputSchema":{"$schema":"https://json-schema.org/draft/2020-12/schema","properties":{"name":{"type":"string","description":"The name of example"}},"type":"object","required":["name"]},"name":"getExampleInformation"}]}}
2025/05/04 04:17:00 main.go:92: stdin: {"method":"resources/list","jsonrpc":"2.0","id":2}
2025/05/04 04:17:00 main.go:117: stdout: {"id":2,"jsonrpc":"2.0","result":{"resources":[]}}
2025/05/04 04:17:00 main.go:92: stdin: {"method":"resources/templates/list","jsonrpc":"2.0","id":3}
2025/05/04 04:17:00 main.go:117: stdout: {"id":3,"jsonrpc":"2.0","result":{"resourceTemplates":[]}}
```
