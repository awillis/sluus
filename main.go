// go:generate goversioninfo
package main

import (
	"github.com/awillis/sluus/cmd"
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	cmd.Execute()
}
