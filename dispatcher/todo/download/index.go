package download

import (
	"animedown/core/todolist"
	"animedown/dispatcher/common"
	"animedown/dispatcher/constants"
	"animedown/util/argparser"
	"animedown/util/terminal"
	"strconv"
)

const (
	Usage       = "d"
	ExplainText = `download anime by index/row, default dir is current.`
	FormatText  = `d <index/row> [dir] [-a: all]`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
	ctx, err := argparser.Parse(args, false, map[string]bool{"a": false})
	if err != nil {
		return terminal.ExitCodeError, err
	}
	if _, ok := ctx.Flags["a"]; ok {
		return downloadAll(this)
	}
	if err := ctx.Check(1); err != nil {
		return terminal.ExitCodeError, err
	}

	dir, _ := this.Get(constants.DirKey).(string) // default dir is "", current dir
	if len(ctx.Args) >= 2 {
		dir = ctx.Args[1]
	}

	index, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return terminal.ExitCodeError, err
	}
	// get magnet from to-do list

	list := this.Get(constants.TodoListKey).(*todolist.TodoList)
	magnet := list.GetMagnet(index)

	err = common.Download(magnet, dir)

	if err != nil {
		return terminal.ExitCodeError, err
	}
	return terminal.ExitCodeOK, nil
}

func downloadAll(this *terminal.TerminalStage) (terminal.ExitCode, error) {
	list := this.Get(constants.TodoListKey).(*todolist.TodoList)
	dir, _ := this.Get(constants.DirKey).(string) // default dir is "", current dir
	for i := 0; i < list.Len(); i++ {
		magnet := list.GetMagnet(i)
		err := common.Download(magnet, dir)
		if err != nil {
			return terminal.ExitCodeError, err
		}
	}
	return terminal.ExitCodeOK, nil
}

func New() *terminal.TerminalStage {
	stage := terminal.NewTerminalStage(Usage, ExplainText, FormatText,
		terminal.WithInitFunc(initFunc),
		terminal.WithLeafStage(),
	)

	return stage
}
