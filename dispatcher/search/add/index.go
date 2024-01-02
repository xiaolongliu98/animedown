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
	FormatText  = `add <index1> <index2> ...`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
	// get searcher
	s := this.Get(constants.SearcherKey).(*search.Searcher)
	// get TodoList
	list := this.Get(constants.TodoListKey).(*todolist.TodoList)

	ctx, err := argparser.Parse(args, false)
	if err != nil {
		return terminal.ExitCodeError, err
	}
	if err := ctx.Check(1); err != nil {
		return terminal.ExitCodeError, err
	}

	for _, arg := range ctx.Args {
		index, err := strconv.Atoi(arg)
		if err != nil {
			return terminal.ExitCodeError, err
		}
		// add
		err = list.Add(s.GetRowSlice()[index])
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
