package main

import (
	"animedown/dispatcher/search"
	"animedown/terminal"
	"animedown/util"
	"fmt"
)

func main() {
	const intro = `*** Welcome to use AnimeDown! *** 
 - AnimeDown is designed to download anime easily from DongManHuaYuan.
 - Refer: https://share.dmhy.org/
`

	t := terminal.NewTerminal(intro)

	_ = t.AddStage(search.New())

	if err := t.Run(); err != nil {
		// Fatal error
		fmt.Println(err)

		fmt.Println("Press any key to exit...")
		util.ReadLine()
	}
}
