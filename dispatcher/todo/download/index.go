package download

import (
	"animedown/core/todolist"
	"animedown/dispatcher/common"
	"animedown/dispatcher/constants"
	"animedown/util/argparser"
	"animedown/util/terminal"
	"fmt"
	"strconv"
	"strings"
)

const (
	Usage       = "d"
	ExplainText = `download anime by index/row, default dir is current.`
	FormatText  = `d <index1> <index2> ... [-d:dir] [-a, download all]`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
	list := this.Get(constants.TodoListKey).(*todolist.TodoList)

	ctx, err := argparser.Parse(args, false, map[string]bool{"a": false, "d": true})
	if err != nil {
		return terminal.ExitCodeError, err
	}
	if _, ok := ctx.Flags["a"]; ok {
		// reset args
		if err := resetArgsAll(ctx, list); err != nil {
			return terminal.ExitCodeError, err
		}
	}

	if err := ctx.Check(1); err != nil {
		return terminal.ExitCodeError, err
	}

	dir, _ := this.Get(constants.DirKey).(string) // default dir is "", current dir
	if targetDir, ok := ctx.Flags["d"]; ok {
		dir = targetDir
	}

	for _, arg := range ctx.Args {
		index, err := strconv.Atoi(arg)
		if err != nil {
			return terminal.ExitCodeError, err
		}
		magnet := list.GetMagnet(index)
		if err = common.Download(magnet, dir); err != nil {
			if strings.Contains(err.Error(), "download cancelled") {
				return terminal.ExitCodeError, err
			}
			fmt.Println(err.Error())
		}
	}

	return terminal.ExitCodeOK, nil
}

func resetArgsAll(ctx *argparser.ArgContext, list *todolist.TodoList) error {
	args := ctx.Args
	m := make(map[int]struct{})
	for _, arg := range args {
		index, err := strconv.Atoi(arg)
		if err != nil {
			return err
		}
		m[index] = struct{}{}
	}

	for i := 0; i < list.Len(); i++ {
		if _, ok := m[i]; !ok {
			args = append(args, strconv.Itoa(i))
		}
	}
	ctx.Args = args
	return nil
}

func New() *terminal.TerminalStage {
	stage := terminal.NewTerminalStage(Usage, ExplainText, FormatText,
		terminal.WithInitFunc(initFunc),
		terminal.WithLeafStage(),
	)
	return stage
}
