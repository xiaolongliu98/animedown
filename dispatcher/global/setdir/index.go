package setdir

import (
	"animedown/dispatcher/constants"
	"animedown/util/argparser"
	"animedown/util/terminal"
	"path/filepath"
)

const (
	Usage       = "setdir"
	ExplainText = `set default download directory.`
	FormatText  = `setdir <dir>`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
	ctx, err := argparser.Parse(args, false)
	if err != nil {
		return terminal.ExitCodeError, err
	}
	if err := ctx.Check(1); err != nil {
		return terminal.ExitCodeError, err
	}

	dir := ctx.Args[0]
	dir, err = filepath.Abs(dir)
	if err != nil {
		return terminal.ExitCodeError, err
	}

	this.Set(constants.DirKey, dir)

	return terminal.ExitCodeOK, nil
}

func New() *terminal.TerminalStage {
	stage := terminal.NewTerminalStage(Usage, ExplainText, FormatText,
		terminal.WithInitFunc(initFunc),
		terminal.WithLeafStage(),
	)

	return stage
}
