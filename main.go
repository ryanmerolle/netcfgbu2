package main

import (
	"github.com/ryanmerolle/netcfgbu2/cmd"
	"github.com/ryanmerolle/netcfgbu2/utils"
)

func main() {
	utils.InitLogger()
	cmd.Execute()
}
