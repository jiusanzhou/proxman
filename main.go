package main

import (
	"fmt"

	"go.zoe.im/proxman/cmd"

	_ "go.zoe.im/proxman/cmd/serve"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}
