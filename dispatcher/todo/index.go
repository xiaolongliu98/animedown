package todo

import (
	"animedown/core/todolist"
	"animedown/dispatcher/common"
	"animedown/dispatcher/constants"
	"animedown/dispatcher/todo/del"
	"animedown/dispatcher/todo/download"
	"animedown/util/terminal"
)

const (
	Usage       = "todo"
	ExplainText = `enter TODO list.`
	FormatText  = `todo`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
	// Get TodoList
	list := this.Get(constants.TodoListKey).(*todolist.TodoList)
	common.PrintTodoList(list)
	this.PrintChildrenUsage()
	return terminal.ExitCodeOK, nil
}

func New() *terminal.TerminalStage {
	stage := terminal.NewTerminalStage(Usage, ExplainText, FormatText,
		terminal.WithInitFunc(initFunc),
		terminal.WithNoEntryGuide(),
		terminal.WithNoPrintDefaultGuideFuncGuideText(),
		terminal.WithNoPrintDefaultCMDUsage(),
	)

	err := stage.AddChild(
		download.New(),
		del.New(),
	)
	if err != nil {
		panic(err)
	}

	return stage
}
