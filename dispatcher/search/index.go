package search

import (
	"animedown/argparser"
	"animedown/dispatcher/common"
	"animedown/dispatcher/constants"
	"animedown/dispatcher/search/download"
	"animedown/dispatcher/search/next"
	"animedown/dispatcher/search/prev"
	"animedown/search"
	"animedown/terminal"
	"fmt"
	"strings"
)

const (
	Usage       = "s"
	ExplainText = `search anime by name, and then show the result.`
	FormatText  = `s <key1> <key2> ...`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
	ctx, err := argparser.Parse(args, false)
	if err != nil {
		return terminal.ExitCodeError, err
	}
	if err := ctx.Check(1); err != nil {
		return terminal.ExitCodeError, err
	}

	// reshow args
	fmt.Println(strings.Join(append([]string{this.Usage}, args...), " "))

	// init searcher
	keyword := strings.Join(ctx.Args, " ")
	s := search.NewSearcher(keyword, search.WithMustKeywordFilter(ctx.Args...))
	if err := s.Search(); err != nil {
		return terminal.ExitCodeError, err
	}
	// show
	common.ShowSearchResult(s)

	// show usage
	this.RunDefaultGuideFunc()

	// save to context
	this.Set(constants.SearcherKey, s)
	return terminal.ExitCodeOK, nil
}

func exitFunc(this *terminal.TerminalStage, exitCode terminal.ExitCode, err error) (terminal.ExitCode, error) {
	// free searcher
	this.Delete(constants.SearcherKey)
	return exitCode, err
}

func New() *terminal.TerminalStage {
	stage := terminal.NewTerminalStage(Usage, ExplainText, FormatText,
		terminal.WithInitFunc(initFunc),
		terminal.WithExitFunc(exitFunc),
		terminal.WithNoEntryGuide(),
		terminal.WithNoPrintDefaultGuideText(),
	)

	err := stage.AddChild(
		download.New(),
		next.New(),
		prev.New(),
	)
	if err != nil {
		panic(err)
	}
	if err := stage.AddRecallSelfChild(); err != nil {
		panic(err)
	}

	return stage
}
