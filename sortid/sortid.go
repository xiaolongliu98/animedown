package sortid

import "sync"

/**
<select name="sort_id" id="AdvSearchSort">
        <option value="0">全部</option>
        <option value="2" style="color: red">動畫</option>
        <option value="31" style="color: red">季度全集</option>
        <option value="3" style="color: green">漫畫</option>
        <option value="41" style="color: green">港台原版</option>
        <option value="42" style="color: green">日文原版</option>
        <option value="4" style="color: purple">音樂</option>
        <option value="43" style="color: purple">動漫音樂</option>
        <option value="44" style="color: purple">同人音樂</option>
        <option value="15" style="color: purple">流行音樂</option>
        <option value="6" style="color: blue">日劇</option>
        <option value="7" style="color: orange">ＲＡＷ</option>
        <option value="9" style="color: #0eb9e7">遊戲</option>
        <option value="17" style="color: #0eb9e7">電腦遊戲</option>
        <option value="18" style="color: #0eb9e7">電視遊戲</option>
        <option value="19" style="color: #0eb9e7">掌機遊戲</option>
        <option value="20" style="color: #0eb9e7">網絡遊戲</option>
        <option value="21" style="color: #0eb9e7">遊戲周邊</option>
        <option value="12" style="color: brown">特攝</option>
        <option value="1" style="color: black">其他</option>
    </select>
*/

const (
	SortV1All    = "全部"
	SortV1Anime  = "動畫"
	SortV1Season = "季度全集"
)

var (
	m    map[string]int
	once sync.Once
)

// InitSortIDMapV1 update to 2023-12-30
func initSortIDMapV1() map[string]int {
	m := map[string]int{
		"全部":   0,
		"動畫":   2,
		"季度全集": 31,
		"漫畫":   3,
		"港台原版": 41,
		"日文原版": 42,
		"音樂":   4,
		"動漫音樂": 43,
		"同人音樂": 44,
		"流行音樂": 15,
		"日劇":   6,
		"ＲＡＷ":  7,
		"遊戲":   9,
		"電腦遊戲": 17,
		"電視遊戲": 18,
		"掌機遊戲": 19,
		"網絡遊戲": 20,
		"遊戲周邊": 21,
		"特攝":   12,
		"其他":   1,
	}
	return m
}

func Map() map[string]int {
	once.Do(func() {
		m = initSortIDMapV1()
	})
	return m
}
