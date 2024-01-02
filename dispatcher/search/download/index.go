package download

import (
	"animedown/core/search"
	"animedown/dispatcher/common"
	"animedown/dispatcher/constants"
	"animedown/util/argparser"
	terminal2 "animedown/util/terminal"
	"strconv"
)

const (
	Usage       = "d"
	ExplainText = `download anime by index/row, default dir is current.`
	FormatText  = `d <index/row> [dir]`
)

func initFunc(this *terminal2.TerminalStage, args []string) (terminal2.ExitCode, error) {
	ctx, err := argparser.Parse(args, false)
	if err != nil {
		return terminal2.ExitCodeError, err
	}
	if err := ctx.Check(1); err != nil {
		return terminal2.ExitCodeError, err
	}

	dir, _ := this.Get(constants.DirKey).(string) // default dir is "", current dir
	if len(ctx.Args) >= 2 {
		dir = ctx.Args[1]
	}

	index, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return terminal2.ExitCodeError, err
	}
	// get magnet from search result
	s := this.Get(constants.SearcherKey).(*search.Searcher)
	magnet := s.GetMagnetLink(index)
	err = common.Download(magnet, dir)
	if err != nil {
		return terminal2.ExitCodeError, err
	}
	return terminal2.ExitCodeOK, nil
}

func New() *terminal2.TerminalStage {
	stage := terminal2.NewTerminalStage(Usage, ExplainText, FormatText,
		terminal2.WithInitFunc(initFunc),
		terminal2.WithLeafStage(),
	)

	return stage
}
