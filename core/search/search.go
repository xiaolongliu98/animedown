package search

import (
	"animedown/util"
	"io"
	"strings"
)

const (
	FilterPostDate = "張貼日期"
	FilterCategory = "分類"
	FilterTitle    = "標題"
	FilterMagnet   = "磁鏈"
	FilterSize     = "大小"
	FilterSeed     = "種子"
	FilterDownload = "下載"
	FilterComplete = "完成"
	FilterUploader = "發佈人"
)

var (
	AllFilters = []string{
		FilterPostDate, FilterCategory, FilterTitle, FilterMagnet, FilterSize, FilterSeed, FilterDownload, FilterComplete, FilterUploader,
	}
	ShowFilters = []string{
		FilterPostDate, FilterCategory, FilterTitle, FilterSize,
	}
)

func AllFilterRowFilter(allRow []string, filters []string) []string {
	var (
		filterRow []string
	)
	for _, filter := range filters {
		filterRow = append(filterRow, allRow[GetAllFilterIndex(filter)])
	}
	return filterRow
}

func GetAllFilterIndex(filter string) int {
	for i, f := range AllFilters {
		if f == filter {
			return i
		}
	}
	return -1
}

// RawSearch returns the raw search result
//
// @return: theadSlice, rowSlice, headIdxMap, err
func RawSearch(keyword string, opts ...SearchOption) (r1 []string, r2 [][]string, r3 map[string]int, err error) {
	defer func() {
		if r := recover(); r != nil {
			r1 = make([]string, 0)
			r2 = make([][]string, 0)
			r3 = make(map[string]int)
			err = nil
			//err = errors.New("页面内容为空")
		}
	}()
	client := util.GetHTTPClient(util.ClashProxyURL)
	req := GetSearchRequest(keyword, opts...)

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, nil, err
	}
	content := string(body)
	const (
		leftUniqueTag  = `class="tablesorter"` // elem: table
		rightUniqueTag = "</table>"
	)
	// clip left content
	leftIndex := strings.Index(content, leftUniqueTag)

	content = content[leftIndex:]
	// clip right content
	content = content[:strings.Index(content, rightUniqueTag)]

	// get thead content
	thead := content[strings.Index(content, "<thead>")+len("<thead>"):]
	thead = thead[:strings.Index(thead, "</thead>")]
	thead = strings.TrimSpace(thead)
	thead = thead[len("<tr>") : len(thead)-len("</tr>")] // remove <tr> </tr>
	thead = strings.TrimSpace(thead)
	// get tbody content
	tbody := content[strings.Index(content, "<tbody>")+len("<tbody>"):]
	tbody = tbody[:strings.Index(tbody, "</tbody>")]
	tbody = strings.TrimSpace(tbody)

	// extract thead struct
	var headIdxMap = make(map[string]int)
	var theadSlice []string
	theadLines := strings.Split(thead, "</th>")
	for i, line := range theadLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// remove <th ...>
		line = line[strings.Index(line, ">")+1:]
		line = strings.TrimSpace(util.RemoveHTMLElem(line))
		theadSlice = append(theadSlice, line)
		headIdxMap[line] = i
	}

	// extract tbody struct
	// tbody one line like:
	// ```
	// <tr ...>
	//   <td ...>
	//   ...
	// </tr>
	// ```
	// totally len(theadSlice) x <td>...</td>

	// firstly, we get each line
	var rowSlice [][]string
	tbodyLines := strings.Split(tbody, "</tr>")
	for _, line := range tbodyLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// remove <tr ...>
		line = line[strings.Index(line, ">")+1:]
		line = strings.TrimSpace(line)
		// for each line/table-row, we get each <td>...</td>
		var row []string
		tdLines := strings.Split(line, "</td>")
		for i, tdLine := range tdLines {
			tdLine = strings.TrimSpace(tdLine)
			if tdLine == "" {
				continue
			}
			// remove <td ...>
			tdLine = tdLine[strings.Index(tdLine, ">")+1:]
			// extract content
			extractor := GetExtractor(theadSlice[i])
			tdLine = extractor(tdLine)

			row = append(row, tdLine)
		}
		rowSlice = append(rowSlice, row)
	}

	return theadSlice, rowSlice, headIdxMap, nil
}
