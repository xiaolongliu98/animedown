package download

import (
	"animedown/core/search"
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
	FormatText  = `d <index1> <index2> ... [-d:dir]`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
	searcher := this.Get(constants.SearcherKey).(*search.Searcher)

	ctx, err := argparser.Parse(args, false, map[string]bool{"d": true})
	if err != nil {
		return terminal.ExitCodeError, err
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
		// get magnet from searcher
		magnet := searcher.GetMagnetLink(index)
		if err = common.Download(magnet, dir); err != nil {
			if strings.Contains(err.Error(), "download cancelled") {
				return terminal.ExitCodeError, err
			}
			fmt.Println(err)
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
