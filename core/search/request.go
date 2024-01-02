package search

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// GetSearchRequest 获取请求
//
// @page 从1开始
func GetSearchRequest(keyword string, opts ...SearchOption) *http.Request {
	// advance args: &sort_id=0&team_id=669&order=date-desc
	// 创建一个url qury参数
	queries := url.Values{}
	queries.Add("keyword", keyword)
	queries.Add("order", "date-desc")
	page := 1
	for _, opt := range opts {
		opt(&queries)
	}
	if queries.Has("page") {
		page, _ = strconv.Atoi(queries.Get("page"))
		if page < 1 {
			page = 1
		}
		queries.Del("page")
	}

	targetURL := fmt.Sprintf("https://share.dmhy.org/topics/list/page/%d?%s",
		page, queries.Encode())
	req, _ := http.NewRequest("GET", targetURL, nil)
	// req.Header.Set("Content-Type", "application/json")
	// 设置一些默认浏览器参数（例如谷歌浏览器的User-Agent等）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	return req
}

type SearchOption func(values *url.Values)

func withPage(page int) SearchOption {
	return func(values *url.Values) {
		values.Set("page", strconv.Itoa(page))
	}
}

// withTeamID 设置team_id, TeamID等同于GroupID
func withTeamID(teamID int) SearchOption {
	return func(values *url.Values) {
		values.Set("team_id", strconv.Itoa(teamID))
	}
}

func withSortID(sortID int) SearchOption {
	return func(values *url.Values) {
		values.Set("sort_id", strconv.Itoa(sortID))
	}
}
