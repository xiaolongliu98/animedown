package util

import (
	"bufio"
	"fmt"
	"github.com/nsf/termbox-go"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const (
	ClashProxyURL = "http://127.0.0.1:7890"
)

// RemoveHTMLElem 移除内容中的HTML元素
func RemoveHTMLElem(content string) string {
	// such as: <td nowrap="nowrap" align="center">CONTENT-A<span class="btl_1">CONTENT-B</span></td>CONTENT-C
	// we want: CONTENT-ACONTENT-BCONTENT-C
	sb := strings.Builder{}
	buf := strings.Builder{}
	pair := 0
	for i := 0; i < len(content); i++ {
		switch content[i] {
		case '<':
			pair++
		case '>':
			if pair == 0 {
				sb.WriteByte(content[i])
			} else {
				pair--
				buf.Reset()
			}
		default:
			if pair == 0 {
				sb.WriteByte(content[i])
			} else {
				buf.WriteByte(content[i])
			}
		}
	}

	if buf.Len() > 0 {
		sb.WriteString(buf.String())
	}
	return sb.String()
}

// cache variables
var (
	getProxyClientLock sync.Mutex
	proxyClientCache   map[string]*http.Client
)

func GetHTTPClient(httpProxyURL ...string) *http.Client {
	if len(httpProxyURL) == 0 || (len(httpProxyURL) > 0 && httpProxyURL[0] == "") {
		return http.DefaultClient
	}
	proxyURL := httpProxyURL[0]

	getProxyClientLock.Lock()
	defer getProxyClientLock.Unlock()
	url, _ := url.Parse(proxyURL)

	if proxyClientCache == nil {
		proxyClientCache = make(map[string]*http.Client)
	}

	// lookup cache
	if client, ok := proxyClientCache[proxyURL]; ok {
		return client
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(url),
	}

	client := &http.Client{Transport: transport}
	// set cache
	proxyClientCache[proxyURL] = client
	return client
}

func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	bytes, _, err := reader.ReadLine()
	return string(bytes), err
}

func StandardizeCMDString(str string) string {
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, "\t", " ")
	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, "\n", "")
	for strings.Contains(str, "  ") {
		str = strings.ReplaceAll(str, "  ", " ")
	}
	return str
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func GetTerminalSize() (int, int, error) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	w, h := termbox.Size()
	termbox.Close()
	return w, h, nil
}
