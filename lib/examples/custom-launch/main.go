// Package main ...
package main

import (
	"fmt"

	"github.com/halicoming/rod"
	"github.com/halicoming/rod/lib/launcher"
)

func main() {
	l := launcher.New()

	// For more info: https://pkg.go.dev/github.com/halicoming/rod/lib/launcher
	u := l.MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()

	page := browser.MustPage("http://example.com").MustWaitStable()

	fmt.Println(page.MustInfo().Title)
}
