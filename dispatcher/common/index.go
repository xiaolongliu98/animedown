package common

import (
	"animedown/search"
	"fmt"
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
