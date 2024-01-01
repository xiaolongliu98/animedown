package main

import (
	"animedown/downloader"
	"animedown/search"
	"animedown/util"
	"context"
	"fmt"
	"strconv"
	"strings"
)

var (
	s *search.Searcher
)

func showSearchResult() {
	heads := search.ShowFilters
	fmt.Printf("Index\t")
	for _, head := range heads {
		fmt.Printf("%s\t\t", head)
	}
	fmt.Println()
	for index, row := range s.GetRowSlice(heads...) {
		fmt.Printf("%d\t", index)
		for _, elem := range row {
			fmt.Printf("%s\t\t", elem)
		}
		fmt.Println()
	}

	// show page
	fmt.Printf("Page: %d, pp <- . -> np\n", s.GetCurrentPage())
}

func CMDSearch(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please enter keyword")
	}

	keyword := strings.Join(args, " ")

	s = search.NewSearcher(keyword, search.WithMustKeywordFilter(args...))
	if err := s.Search(); err != nil {
		return err
	}
	showSearchResult()
	return nil
}

func CMDNextPage(args []string) error {
	if s == nil {
		return fmt.Errorf("please search first")
	}
	if err := s.NextPage(); err != nil {
		return err
	}
	util.ClearScreen()
	showSearchResult()
	return nil
}

func CMDPrevPage(args []string) error {
	if s == nil {
		return fmt.Errorf("please search first")
	}
	if err := s.PrevPage(); err != nil {
		return err
	}
	util.ClearScreen()
	showSearchResult()
	return nil
}

func CMDDownload(args []string, ctx ...context.Context) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify the index(row number)")
	}
	if len(ctx) == 0 {
		ctx = append(ctx, context.Background())
	}
	dir := "."
	if len(args) > 1 {
		dir = args[1]
	}

	if s == nil {
		return fmt.Errorf("please search first")
	}
	done := s.GetSearchDone()
	if !done {
		return fmt.Errorf("please search first")
	}

	if s.IsEmpty() {
		return fmt.Errorf("no result")
	}

	// idx starts from 0
	idx, ok := strconv.Atoi(args[0])
	if ok != nil {
		return fmt.Errorf("invalid index")
	}
	if idx < 0 || idx >= len(s.GetRowSlice()) {
		return fmt.Errorf("index out of range")
	}

	link := s.GetMagnetLink(idx)
	return downloader.Download(ctx[0], dir, "", link)
}

func CMDClearScreen(args []string) {
	s = nil
}
