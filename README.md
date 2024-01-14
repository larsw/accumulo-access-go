# Accumulo Access Expressions for Go

## Introduction

This package provides a simple way to parse and evaluate Accumulo access expressions in Go, based on the [AccessExpression specification](https://github.com/apache/accumulo-access/blob/main/SPECIFICATION.md).

## Usage

```go
package main

import (
	"fmt"
	accumulo "github.com/larsw/accumulo-access-go/pkg"
)

func main() {
	res, err := accumulo.CheckAuthorization("A & B & (C | D)", "A,B,C")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	// Print the result
	fmt.Printf("%v\n", res)
}
```

* Lars Wilhelmsen (https://github.com/larsw/)

## License

Licensed under the MIT License [LICENSE](LICENSE).

