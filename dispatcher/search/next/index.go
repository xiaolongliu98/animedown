package next

import (
	"animedown/core/search"
	"animedown/dispatcher/common"
	"animedown/dispatcher/constants"
	"animedown/util/argparser"
	terminal2 "animedown/util/terminal"
)

const (
	Usage       = "np"
	ExplainText = `next page.`
	FormatText  = `np`
)

func initFunc(this *terminal2.TerminalStage, args []string) (terminal2.ExitCode, error) {
	ctx, err := argparser.Parse(args, false)
	if err != nil {
		return terminal2.ExitCodeError, err
	}
	if err := ctx.Check(0); err != nil {
		return terminal2.ExitCodeError, err
	}

	// get searcher
	s := this.Get(constants.SearcherKey).(*search.Searcher)
	// next page
	err = s.NextPage()
	if err != nil {
		return terminal2.ExitCodeError, err
	}
	// show
	common.ShowSearchResult(s)
	// show usage
	this.Parent.RunDefaultGuideFunc()

	return terminal2.ExitCodeOK, nil
}

func New() *terminal2.TerminalStage {
	stage := terminal2.NewTerminalStage(Usage, ExplainText, FormatText,
		terminal2.WithInitFunc(initFunc),
		terminal2.WithLeafStage(),
	)

	return stage
}
