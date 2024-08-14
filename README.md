# respparser-nabil

Reading Resp string and returning it as an array

## Installation

As a library

```shell
go get github.com/codescalersinternships/respparser-nabil.git/pkg
```

## Usage

in your Go app you can do something like

```go
package main

import (
	"fmt"

	respparser "github.com/codescalersinternships/respparser-nabil.git/pkg"
)

func main() {
	str := ",123.45\r\n"

	x,err := respparser.Parser(str);
	if err != nil {
		fmt.Println(err)
		return
	}
	for _,ele := range x{

		fmt.Printf("(%v, %T)\n",ele,ele)
	}
}

```

## Testing

```shell
make test
```