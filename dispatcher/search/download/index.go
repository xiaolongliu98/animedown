package download

import (
	"animedown/argparser"
	"animedown/dispatcher/constants"
	"animedown/downloader"
	"animedown/search"
	"animedown/terminal"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const (
	Usage       = "d"
	ExplainText = `download anime by index/row, default dir is current.`
	FormatText  = `d <index/row> [dir]`
)

func initFunc(this *terminal.TerminalStage, args []string) (terminal.ExitCode, error) {
	ctx, err := argparser.Parse(args, false)
	if err != nil {
		return terminal.ExitCodeError, err
	}
	if err := ctx.Check(1); err != nil {
		return terminal.ExitCodeError, err
	}

	dir := ""
	if len(ctx.Args) >= 2 {
		dir = ctx.Args[1]
	}

	index, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		return terminal.ExitCodeError, err
	}
	// get magnet from search result
	s := this.Get(constants.SearcherKey).(*search.Searcher)
	magnet := s.GetMagnetLink(index)

	signalCtx, cancel := context.WithCancel(context.Background())
	// catch signal(Ctrl+C)
	// 创建一个信号通道
	signalCh := make(chan os.Signal, 1)
	doneCh := make(chan struct{}, 1)

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-signalCh:
			cancel()
		case <-doneCh:
			return
		}
	}()

	// do download
	err = downloader.DownloadBlocked(&downloader.DownloadConfig{
		Dir:       dir,
		Magnet:    magnet,
		Ctx:       signalCtx,
		Obs:       Observer,
		ObsPeriod: 450,
	})
	doneCh <- struct{}{}

	if err != nil {
		return terminal.ExitCodeError, err
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

func Observer(stage downloader.Stage, readBytes int64, totalBytes int64, activePeers, totalPeers int, downloadRate string, percentage float64) {
	switch stage {
	case downloader.StagePrepare:
		fmt.Print("[download] Get torrent metadata...\r")
	case downloader.StageDownload:
		totalMiBytes := float64(totalBytes) / 1024 / 1024
		readMiBytes := min(totalMiBytes, float64(readBytes)/1024/1024)
		progress := fmt.Sprintf("%.2f/%.2fMiB, %s, %d peers, %.2f%%", readMiBytes, totalMiBytes, downloadRate, activePeers, percentage)
		fmt.Printf("\r%s", strings.Repeat(" ", len(progress)+16)) // 清除当前行
		fmt.Printf("\r[download] %s", progress)
	case downloader.StageSuccess:
		fmt.Println("\n[download] Download success!")
	case downloader.StageCancel:
		fmt.Println("\n[download] Download cancel.")
	}
}
