package main

import (
	"animedown/dispatcher/search"
	"animedown/terminal"
	"animedown/util"
	"fmt"
)

const (
	GuideText = `Welcome to use AnimeDown! 
Now, you can use the following commands to download anime:
  1. Search anime by name: "s [key1] [key2] ..."
  2. Next page: "np"
  3. Previous page: "pp"
  4. Download anime by magnet: "d [index] [OPTIONAL, target_dir, empty is '.']", 
     index is the row of the current search result.

  5. Exit: "q"
  6. Clear screen: "cls"

Please enter your command:`
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
