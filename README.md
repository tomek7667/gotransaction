# gotransaction

library with universal transactions that allow custom rollbacking.

Useful especially for situation where various clients already commited changes, whereas last ones errored.

installation:

```bash
go get github.com/tomek7667/gotransaction
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/tomek7667/gotransaction/transaction"
)


func main() {
	tx := transaction.New()
	err := tx.A(func() error {
		return ActionThatMightChangeStateOrError()
	}, func() error {
		return ActionThatRollsBackTheActionAbove()
	})
	if err != nil {
		originalErr, rollbackErr := tx.R(err)
		if rollbackErr != nil {
			panic(fmt.Errorf("rollbacking failed: %w", err))
		}
	}
}
```

For a detailed usage with a specific use case see [scenario 1](./examples/scenario1/main.go)
