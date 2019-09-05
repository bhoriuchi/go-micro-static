# go-micro-static
Serve static files from go-micro api server index. The plugin was written to work with single page applications.

## Usage

Register the plugin before building Micro

```go
package main

import (
	"github.com/micro/micro/plugin"
	"github.com/bhoriuchi/go-micro-static"
)

func init() {
	plugin.Register(static.NewPlugin())
}
```

The static directory can be supplied by the `--static_dir` option. The default is `html`

```sh
micro --static_dir=dist
```

When using along side API or WEB services, supply the service names to omit from serving with the plugin

```sh
micro --static_service=foo --static_service=bar
```
