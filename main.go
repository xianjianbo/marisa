package main

import (
	"github.com/xianjianbo/marisa/bootstrap"
	"github.com/xianjianbo/marisa/router"
)

func main() {
	bootstrap.Init()
	router.InitRouter()
}
