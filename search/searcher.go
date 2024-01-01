package search

import (
	"animedown/group"
	"animedown/sortid"
	"fmt"
	"strings"
)

type Searcher struct {
	currentKeyword     string
	currentPage        int
	searchDone         bool
	group              string
	sort               string
	mustKeywordFilters []string

	// Result
	theadSlice []string
	rowSlice   [][]string
	headIdxMap map[string]int
}

func DefaultSearcher(keyword string) *Searcher {
	return &Searcher{
		currentKeyword:     keyword,
		currentPage:        1,
		searchDone:         false,
		group:              "",
		sort:               "",
		mustKeywordFilters: nil,
		theadSlice:         nil,
		rowSlice:           nil,
		headIdxMap:         nil,
	}
}

func NewSearcher(keyword string, opts ...SearcherOption) *Searcher {
	s := DefaultSearcher(keyword)
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type SearcherOption func(*Searcher)

func WithGroup(group string) SearcherOption {
	return func(s *Searcher) {
		s.group = group
	}
}

func WithSort(sort string) SearcherOption {
	return func(s *Searcher) {
		s.sort = sort
	}
}

func WithMustKeywordFilter(filters ...string) SearcherOption {
	return func(s *Searcher) {
		s.mustKeywordFilters = filters
	}
}

func WithPage(page int) SearcherOption {
	return func(s *Searcher) {
		s.currentPage = page
	}
}

func (s *Searcher) Search() error {
	var err error
	if s.searchDone {
		return nil
	}

	opts := []SearchOption{
		withPage(s.currentPage),
	}
	if s.group != "" {
		teamID, ok := group.Map()[s.group]
		if ok {
			opts = append(opts, withTeamID(teamID))
		}
	}
	if s.sort != "" {
		sortID, ok := sortid.Map()[s.sort]
		if ok {
			opts = append(opts, withSortID(sortID))
		}
	}
	s.theadSlice, s.rowSlice, s.headIdxMap, err = RawSearch(s.currentKeyword, opts...)
	if err != nil {
		return err
	}
	// apply must keyword filters
	s.doMustKeywordFilter()

	s.searchDone = true
	return nil
}

// doMustKeywordFilter filter title
func (s *Searcher) doMustKeywordFilter() {
	count0 := len(s.mustKeywordFilters)
	if count0 == 0 {
		return
	}

	rowSlice := s.rowSlice[:0]
	for _, row := range s.rowSlice {
		count1 := 0
		// filter title
		for _, filter := range s.mustKeywordFilters {
			if !strings.Contains(row[s.headIdxMap[FilterTitle]], filter) {
				break
			}
			count1++
		}

		if count1 == count0 {
			rowSlice = append(rowSlice, row)
		}
	}

	s.rowSlice = rowSlice
}

func (s *Searcher) NextPage() error {
	s.currentPage++
	s.searchDone = false
	return s.Search()
}

func (s *Searcher) PrevPage() error {
	if s.currentPage == 1 {
		return fmt.Errorf("already at the first page")
	}
	s.currentPage--
	s.searchDone = false
	return s.Search()
}

func (s *Searcher) JumpPage(page int) error {
	if page < 1 {
		return fmt.Errorf("page number must be greater than 0")
	}
	s.currentPage = page
	s.searchDone = false
	return s.Search()
}

func (s *Searcher) GetCurrentPage() int {
	return s.currentPage
}

func (s *Searcher) GetTheadSlice() []string {
	return s.theadSlice
}

func (s *Searcher) GetRowSlice(headFilters ...string) [][]string {
	if len(headFilters) == 0 {
		return s.rowSlice
	}

	rowSlice := make([][]string, len(s.rowSlice), len(s.rowSlice))
	for i, row := range s.rowSlice {
		rowSlice[i] = make([]string, len(headFilters), len(headFilters))
		for j, filter := range headFilters {
			rowSlice[i][j] = row[s.headIdxMap[filter]]
		}
	}
	return rowSlice
}

func (s *Searcher) GetHeadIdxMap() map[string]int {
	return s.headIdxMap
}

func (s *Searcher) GetKeywordFilters() []string {
	return s.mustKeywordFilters
}

func (s *Searcher) GetKeyword() string {
	return s.currentKeyword
}

func (s *Searcher) GetGroup() string {
	return s.group
}

func (s *Searcher) GetSort() string {
	return s.sort
}

func (s *Searcher) GetSearchDone() bool {
	return s.searchDone
}

func (s *Searcher) GetMagnetLink(row int) string {
	if row < 0 || row >= len(s.rowSlice) {
		return ""
	}
	return s.rowSlice[row][s.headIdxMap[FilterMagnet]]
}

// IsEmpty return true if no result
func (s *Searcher) IsEmpty() bool {
	return len(s.rowSlice) == 0
}

// Clear
func (s *Searcher) Clear() {
	s.currentKeyword = ""
	s.currentPage = 1
	s.searchDone = false
	s.group = ""
	s.sort = ""
	s.mustKeywordFilters = nil
	s.theadSlice = nil
	s.rowSlice = nil
	s.headIdxMap = nil
}
