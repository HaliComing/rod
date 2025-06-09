// Package main ...
package main

import (
	"fmt"

	"github.com/halicoming/rod/lib/launcher"
	"github.com/halicoming/rod/lib/utils"
)

func main() {
	p, err := launcher.NewBrowser().Get()
	utils.E(err)

	fmt.Println(p)
}
