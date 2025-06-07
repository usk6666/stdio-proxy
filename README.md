# stdio-proxy

## Overview

`stdio-logging-proxy` proxies the standard input and standard output of a specified command and logs all communication to a log file. This allows you to easily check the messages sent and received when debugging communication with the MCP server.

## Features

*   **Standard Input/Output Proxy:** Proxies the standard input and standard output of the specified command.
*   **Communication Logging:** Logs all communication to a log file.
*   **Configuration via Command-Line Arguments:** You can configure the log file path and the command to be proxied using command-line arguments.

## Usage

### Build

```bash
git clone https://github.com/usk6666/stdio-proxy.git
cd stdio-proxy
go build
```

### Example Execution

```bash
./stdio-proxy exec /path/to/mcp-server mcp-args
./stdio-proxy shell /path/to/mcp-server mcp-args
```

## Command-Line Arguments

*   `--output string`: Path to the output log file (default: `stdio-proxy.log`).

## Configuration

### Environment Variables

N/A

### Example MCP Server Configuration for Roo Code

- .roo/mcp.json
```json
{
  "mcpServers": {
    "example": {
      "name": "example-mcp",
      "description": "example mcp server",
      "command": "/path/to/stdio-proxy /path/to/mcp-server args",
      "timeout": 30,
      "alwaysAllow": [],
      "disabled": false
    }
  }
}
```

## Logging

The log file records the contents of standard input, standard output, and standard error output.

- Log Image
```log
2025/06/07 11:57:39 stdio-proxy-msg: Starting proxy...
2025/06/07 11:57:39 stdio-proxy-msg: Name: /path/to/mcp-server
2025/06/07 11:57:39 stdio-proxy-msg: Args: [args]
2025/05/04 04:17:00 stdin: {"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"Roo Code","version":"3.15.3"}},"jsonrpc":"2.0","id":0}
2025/05/04 04:17:00 stdout: {"id":0,"jsonrpc":"2.0","result":{"capabilities":{"prompts":{"listChanged":false},"resources":{"listChanged":false},"tools":{"listChanged":false}},"protocolVersion":"2024-11-05","serverInfo":{"name":"","version":""}}}
2025/05/04 04:17:00 stdin: {"method":"notifications/initialized","jsonrpc":"2.0"}
2025/05/04 04:17:00 stdin: {"method":"tools/list","jsonrpc":"2.0","id":1}
2025/05/04 04:17:00 stdout: {"id":1,"jsonrpc":"2.0","result":{"tools":[{"description":"example mcp description","inputSchema":{"$schema":"https://json-schema.org/draft/2020-12/schema","properties":{"name":{"type":"string","description":"The name of example"}},"type":"object","required":["name"]},"name":"getExampleInformation"}]}}
2025/05/04 04:17:00 stdin: {"method":"resources/list","jsonrpc":"2.0","id":2}
2025/05/04 04:17:00 stdout: {"id":2,"jsonrpc":"2.0","result":{"resources":[]}}
2025/05/04 04:17:00 stdin: {"method":"resources/templates/list","jsonrpc":"2.0","id":3}
2025/05/04 04:17:00 stdout: {"id":3,"jsonrpc":"2.0","result":{"resourceTemplates":[]}}
```
