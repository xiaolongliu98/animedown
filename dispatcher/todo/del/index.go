package del

import (
	"animedown/core/todolist"
	"animedown/dispatcher/common"
	"animedown/dispatcher/constants"
	"animedown/util/argparser"
	"animedown/util/terminal"
	"strconv"
)

const (
	Usage       = "del"
	ExplainText = `delete anime from temporary TODO list.`
	FormatText  = `del <index1> <index2> ...`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
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
		// del
		list.Delete(index)
	}
	// reshow
	common.PrintTodoList(list)
	this.Parent.PrintChildrenUsage()
	return terminal.ExitCodeOK, nil
}

func New() *terminal.TerminalStage {
	stage := terminal.NewTerminalStage(Usage, ExplainText, FormatText,
		terminal.WithInitFunc(initFunc),
		terminal.WithLeafStage(),
	)

	return stage
}
