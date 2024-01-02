package search

import (
	"animedown/util"
	"strings"
)

type Extractor func(string) string

var (
	extractors = map[string]Extractor{
		FilterPostDate: func(s string) string {
			// correct "2023/12/27 23:52                        2023/12/27 23:52"
			s = defaultExtractor(s)
			return strings.TrimSpace(s[0:len("xxxx/xx/xx xx:xx")])
		},

		FilterMagnet: func(s string) string {
			// Get href value
			// <a class="download-arrow arrow-magnet" title="磁力下載" href="magnet:?xt=urn:btih:2JEQBJUADCR3GUPEBGCILZGARLAS4BJF&dn=&tr=http%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=udp%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.prq.to%2Fannounce&tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce&tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud&tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=http%3A%2F%2F208.67.16.113%3A8000%2Fannounce">&nbsp;</a>
			// <a class="download-pp" target="_blank" title="保存后可以用手機隨時觀賞" onclick="_hmt.push(['_trackEvent', 'pikpak', 'click', 'list'])" href="https://mypikpak.com/drive/url-checker?url=magnet:?xt=urn:btih:d24900a68018a3b351e4098485e4c08ac12e0525">&nbsp;</a>
			const (
				left  = `href="`
				right = `">`
			)

			s = s[strings.Index(s, left)+len(left):]
			s = s[:strings.Index(s, right)]
			return s
		},

		FilterTitle: func(s string) string {
			// <span class="tag">
			// <a href="/topics/list/team_id/619">桜都字幕组</a>
			// </span>
			// <a href="/topics/view/659528_Hoshikuzu_Telepath_12_1080p.html" target="_blank">[桜都字幕组] 星灵感应 / Hoshikuzu Telepath [12][1080p][简繁内封]</a>

			// ignore span
			const (
				left = `</span>`
			)
			s = s[strings.Index(s, left)+len(left):]
			s = defaultExtractor(s)
			// 去除结尾 "約?條評論"
			const (
				right1 = `約`
				right2 = `條評論`
			)
			s = strings.TrimSpace(s)
			if strings.HasSuffix(s, right2) && strings.Contains(s, right1) {
				s = s[:strings.LastIndex(s, right1)]
			}
			return s
		},
	}
)

func defaultExtractor(s string) string {
	return strings.TrimSpace(util.RemoveHTMLElem(s))
}

func GetExtractor(filterName string) Extractor {
	if extractor, ok := extractors[filterName]; ok {
		return extractor
	}
	return defaultExtractor
}
