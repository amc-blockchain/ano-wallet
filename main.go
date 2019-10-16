package main

import (
	//_ "blockChainWallet/database"

	_ "blockChainWallet/purseInterface"

	"blockChainWallet/router"

	"github.com/henrylee2cn/faygo"
	//	"runtime"
)

func main() {

	//	runtime.GOMAXPROCS(4)

	router.Route(faygo.New("blockChainWallet", "v1.0"))
	faygo.Run()
}
