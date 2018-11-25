package main

import (
	"flag"

	"github.com/d0zingcat/go-logger/logger"

	"github.com/d0zingcat/labs/m2scrapy/spider"
)

func main() {
	workNum := flag.Int("n", 50, "number of pages to be fetch once a time ")
	dir := flag.String("d", "tmp", "location to download the pics")
	flag.Parse()
	logger.Debug("start working")
	spider.Process(*workNum, *dir)
	logger.Debug("end work")
}
