package common

import (
	"animedown/core/downloader"
	"animedown/core/search"
	"animedown/core/todolist"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func ShowSearchResult(s *search.Searcher) {
	heads := search.ShowFilters
	fmt.Printf("Index\t")
	for _, head := range heads {
		fmt.Printf("%s\t\t", head)
	}
	fmt.Println()
	rows := s.GetRowSlice(heads...)
	for index, row := range rows {
		fmt.Printf("%d\t", index)
		for _, elem := range row {
			fmt.Printf("%s\t\t", elem)
		}
		fmt.Println()
	}
	if len(rows) == 0 {
		fmt.Println("No result found")
		fmt.Println()
	}

	// show page
	fmt.Printf("Page: %d\n", s.GetCurrentPage())
}

func PrintTodoList(list *todolist.TodoList) {
	heads := search.ShowFilters
	fmt.Printf("Index\t")
	for _, head := range heads {
		fmt.Printf("%s\t\t", head)
	}
	fmt.Println()

	rows := list.GetList()

	for index, row := range rows {
		row = search.AllFilterRowFilter(row, search.ShowFilters)
		fmt.Printf("%d\t", index)
		for _, elem := range row {
			fmt.Printf("%s\t\t", elem)
		}
		fmt.Println()
	}
	if len(rows) == 0 {
		fmt.Println("No item.")
		fmt.Println()
	}

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

func Download(magnet, dir string) error {
	signalCtx, cancel := context.WithCancel(context.Background())
	// catch signal(Ctrl+C)
	// 创建一个信号通道
	signalCh := make(chan os.Signal, 1)
	doneCh := make(chan struct{}, 1)

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	cancelled := false
	go func() {
		select {
		case <-signalCh:
			cancel()
			cancelled = true
		case <-doneCh:
			return
		}
	}()

	// do download
	err := downloader.DownloadBlocked(&downloader.DownloadConfig{
		Dir:       dir,
		Magnet:    magnet,
		Ctx:       signalCtx,
		Obs:       Observer,
		ObsPeriod: 450,
	})
	doneCh <- struct{}{}

	if cancelled {
		return fmt.Errorf("download cancelled")
	}
	return err
}
