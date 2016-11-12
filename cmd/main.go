package main

import (
	"fmt"
	"github.com/k0kubun/pp"

	"github.com/mwf/golidays/crawlers"
)

func main() {
	crawler := crawlers.NewConsultantRu()
	holidays, err := crawler.ScrapeYear(2016)
	fmt.Printf("err: %#v", err)
	pp.Print(holidays)
}
