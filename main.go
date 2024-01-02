package main

import (
	"animedown/core/todolist"
	"animedown/dispatcher/constants"
	"animedown/dispatcher/global/setdir"
	"animedown/dispatcher/search"
	"animedown/dispatcher/todo"
	"animedown/util"
	"animedown/util/terminal"
	"fmt"
)

func main() {
	const intro = `*** Welcome to use AnimeDown! *** 
 - AnimeDown is designed to download anime easily from DongManHuaYuan.
 - Refer: https://share.dmhy.org/
`

	t := terminal.NewTerminal(intro, func(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
		// init to-do list
		list := todolist.New()
		this.Set(constants.TodoListKey, list)
		return terminal.ExitCodeOK, nil
	})

	_ = t.AddStage(search.New())
	_ = t.AddStage(todo.New())

	t.AddGlobalStage(setdir.Usage, setdir.New)

	if err := t.Run(); err != nil {
		// Fatal error
		fmt.Println(err)

		fmt.Println("Press any key to exit...")
		util.ReadLine()
	}
}
