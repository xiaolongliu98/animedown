package prev

import (
	"animedown/argparser"
	"animedown/dispatcher/common"
	"animedown/dispatcher/constants"
	"animedown/search"
	"animedown/terminal"
)

const (
	Usage       = "pp"
	ExplainText = `previous page.`
	FormatText  = `pp`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
	ctx, err := argparser.Parse(args, false)
	if err != nil {
		return terminal.ExitCodeError, err
	}
	if err := ctx.Check(0); err != nil {
		return terminal.ExitCodeError, err
	}

	// get searcher
	s := this.Get(constants.SearcherKey).(*search.Searcher)
	// next page
	err = s.PrevPage()
	if err != nil {
		return terminal.ExitCodeError, err
	}
	// show
	common.ShowSearchResult(s)
	// show usage
	this.Parent.RunDefaultGuideFunc()

	return terminal.ExitCodeOK, nil
}

func New() *terminal.TerminalStage {
	stage := terminal.NewTerminalStage(Usage, ExplainText, FormatText,
		terminal.WithInitFunc(initFunc),
		terminal.WithLeafStage(),
	)

	return stage
}
