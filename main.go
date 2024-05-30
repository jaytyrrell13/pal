package main

import (
	"github.com/jaytyrrell13/pal/cmd"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	err := cmd.Execute(version)
	cobra.CheckErr(err)
}
