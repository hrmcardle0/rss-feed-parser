package main

import (
	"github.com/mmcdole/gofeed"
)

func GenerateClient() *gofeed.Parser {
	fp := gofeed.NewParser()
	return fp
}
