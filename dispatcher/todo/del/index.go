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
	FormatText  = `del <index/row>`
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
	// get TodoList
	list := this.Get(constants.TodoListKey).(*todolist.TodoList)
	// del
	list.Delete(index)
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
