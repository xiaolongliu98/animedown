package group

import (
	"animedown/util"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// 常用的组别
const (
	GoupIDAll      = "全部"
	GroupDMHY      = "動漫花園"
	GroupMMNCW     = "喵萌奶茶屋"
	GroupLilith    = "Lilith-Raws"
	GroupYinDu     = "桜都字幕组"
	GroupAni       = "ANi"
	GroupQZGD      = "轻之国度"
	GroupSweetSub  = "SweetSub"
	GroupXingKong  = "星空字幕组"
	GroupZhuShen   = "诸神kamigami字幕组"
	GroupHuanYing  = "幻樱字幕组"
	GroupDMG       = "动漫国字幕组"
	GroupTSDM      = "天使动漫论坛"
	GroupLoliHouse = "LoliHouse"
	GroupJiYing    = "极影字幕社"
)

func TestGetGroupIDMap() {
	for k, v := range Map() {
		fmt.Printf("%s: %d\n", k, v)
	}
}

var (
	once sync.Once
	m    map[string]int
)

func Map() map[string]int {
	once.Do(func() {
		var err error
		m, err = initGroupIDMap()
		if err != nil {
			panic(err)
		}
	})
	return m
}

func initGroupIDMap() (map[string]int, error) {
	// build from "https://share.dmhy.org/topics/advanced-search?team_id=0&sort_id=0&orderby=date-desc"
	// <select name="team_id" id="AdvSearchTeam">
	//        <option value="0">全部</option>
	//        <option value="117">動漫花園</option>
	//        <option value="823">拨雪寻春</option>
	//        <option value="801">NC-Raws</option>
	//        <option value="669">喵萌奶茶屋</option>
	//        <option value="803">Lilith-Raws</option>
	//        ...
	//    </select>

	// <select name="sort_id" id="AdvSearchSort">
	//        <option value="0">全部</option>
	//        <option value="2" style="color: red">動畫</option>
	//        <option value="31" style="color: red">季度全集</option>
	//        <option value="3" style="color: green">漫畫</option>
	//        ...
	//        <option value="1" style="color: black">其他</option>
	// </select>

	const (
		left  = `<select name="team_id" id="AdvSearchTeam">`
		right = `</select>`
	)

	client := util.GetHTTPClient(util.ClashProxyURL)
	req := getGroupRequest()
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()
	// extract group id from resp
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	content := string(body)
	// clip left content
	content = content[strings.Index(content, left)+len(left):]
	// clip right content
	content = content[:strings.Index(content, right)]
	// regex extract
	reg := regexp.MustCompile(`<option value="(\d+)">([^<]+)</option>`)
	submatch := reg.FindAllStringSubmatch(content, -1)
	// build map
	groupIDMap := make(map[string]int)
	for _, match := range submatch {
		id, err := strconv.Atoi(match[1])
		if err != nil {
			fmt.Println(err)
			continue
		}
		groupIDMap[match[2]] = id
	}
	return groupIDMap, nil
}

func getGroupRequest() *http.Request {
	targetURL := "https://share.dmhy.org/topics/advanced-search?team_id=0&sort_id=0&orderby=date-desc"
	req, _ := http.NewRequest("GET", targetURL, nil)
	// req.Header.Set("Content-Type", "application/json")
	// 设置一些默认浏览器参数（例如谷歌浏览器的User-Agent等）
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	return req
}

/**
  <option value="0">全部</option>
  <option value="117">動漫花園</option>
  <option value="823">拨雪寻春</option>
  <option value="801">NC-Raws</option>
  <option value="669">喵萌奶茶屋</option>
  <option value="803">Lilith-Raws</option>
  <option value="648">魔星字幕团</option>
  <option value="619">桜都字幕组</option>
  <option value="767">天月動漫 &amp;發佈組</option>
  <option value="185">极影字幕社</option>
  <option value="657">LoliHouse</option>
  <option value="151">悠哈C9字幕社</option>
  <option value="749">幻月字幕组</option>
  <option value="390">天使动漫论坛</option>
  <option value="303">动漫国字幕组</option>
  <option value="241">幻樱字幕组</option>
  <option value="47">爱恋字幕社</option>
  <option value="805">DBD制作组</option>
  <option value="604">c.c动漫</option>
  <option value="550">萝莉社活动室</option>
  <option value="283">千夏字幕组</option>
  <option value="772">IET字幕組</option>
  <option value="288">诸神kamigami字幕组</option>
  <option value="804">霜庭云花Sub</option>
  <option value="755">GMTeam</option>
  <option value="454">风车字幕组</option>
  <option value="37">雪飄工作室(FLsnow)</option>
  <option value="764">MCE汉化组</option>
  <option value="488">丸子家族</option>
  <option value="731">星空字幕组</option>
  <option value="574">梦蓝字幕组</option>
  <option value="504">LoveEcho!</option>
  <option value="650">SweetSub</option>
  <option value="630">枫叶字幕组</option>
  <option value="479">Little Subbers!</option>
  <option value="321">轻之国度</option>
  <option value="649">云光字幕组</option>
  <option value="520">豌豆字幕组</option>
  <option value="626">驯兽师联盟</option>
  <option value="666">中肯字幕組</option>
  <option value="781">SW字幕组</option>
  <option value="576">银色子弹字幕组</option>
  <option value="434">风之圣殿</option>
  <option value="665">YWCN字幕组</option>
  <option value="228">KRL字幕组</option>
  <option value="49">华盟字幕社</option>
  <option value="627">波洛咖啡厅</option>
  <option value="88">动音漫影</option>
  <option value="581">VCB-Studio</option>
  <option value="407">DHR動研字幕組</option>
  <option value="719">80v08</option>
  <option value="732">肥猫压制</option>
  <option value="680">Little字幕组</option>
  <option value="613">AI-Raws</option>
  <option value="806">离谱Sub</option>
  <option value="812">虹咲学园烤肉同好会</option>
  <option value="636">ARIA吧汉化组</option>
  <option value="75">柯南事务所</option>
  <option value="821">百冬練習組</option>
  <option value="641">冷番补完字幕组</option>
  <option value="765">爱咕字幕组</option>
  <option value="822">極彩字幕组</option>
  <option value="217">AQUA工作室</option>
  <option value="592">未央阁联盟</option>
  <option value="703">届恋字幕组</option>
  <option value="808">夜莺家族</option>
  <option value="734">TD-RAWS</option>
  <option value="447">夢幻戀櫻</option>
  <option value="790">WBX-SUB</option>
  <option value="807">Liella!の烧烤摊</option>
  <option value="814">Amor字幕组</option>
  <option value="813">MingYSub</option>
  <option value="835">小白GM</option>
  <option value="832">Sakura</option>
  <option value="845">PMFAN字幕组</option>
  <option value="817">EMe</option>
  <option value="818">Alchemist</option>
  <option value="819">黑岩射手吧字幕组</option>
  <option value="816">ANi</option>
  <option value="844">DBFC字幕组</option>
  <option value="836">MSB制作組</option>
*/
