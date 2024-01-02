package add

import (
	"animedown/core/search"
	"animedown/core/todolist"
	"animedown/dispatcher/constants"
	"animedown/util/argparser"
	"animedown/util/terminal"
	"strconv"
)

const (
	Usage       = "add"
	ExplainText = `add anime to temporary TODO list.`
	FormatText  = `add <index/row>`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
	ctx, err := argparser.Parse(args, false)
	if err != nil {
		return terminal.ExitCodeError, err
	}
	if err := ctx.Check(1); err != nil {
		return terminal.ExitCodeError, err
	}

	index, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return terminal.ExitCodeError, err
	}
	// get searcher
	s := this.Get(constants.SearcherKey).(*search.Searcher)
	// get TodoList
	list := this.Get(constants.TodoListKey).(*todolist.TodoList)
	// add
	err = list.Add(s.GetRowSlice()[index])
	if err != nil {
		return terminal.ExitCodeError, err
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
