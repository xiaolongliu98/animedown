package downloader

import (
	"animedown/ratecalc"
	"context"
	"fmt"
	alog "github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func TestDownloadMagnet() {
	//magnet69MB := "magnet:?xt=urn:btih:FI6AOOTCXSVVKIYDR4PCUNFZA7W4X6MQ&dn=&tr=http%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=udp%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.prq.to%2Fannounce&tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce&tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud&tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce"
	magnet9MB := `magnet:?xt=urn:btih:AI373HJAXCQ2D24B4AU62UQ4IEKHFZJP&dn=&tr=http%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=udp%3A%2F%2F104.143.10.186%3A8000%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.prq.to%2Fannounce&tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce&tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud&tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce&tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce&tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce`
	//magnetFulilian :=
	magnet := magnet9MB
	err := Download(context.TODO(), "./data", "", magnet)
	if err != nil {
		fmt.Println(err)
	}
}

func Download(ctx context.Context, dir string, fileName string, magnet string, httpProxy ...string) error {
	_ = os.MkdirAll(dir, os.ModePerm)

	config := torrent.NewDefaultClientConfig()
	config.DataDir = dir
	config.Logger = config.Logger.WithFilterLevel(alog.Error)
	config.Logger.SetHandlers()

	if len(httpProxy) != 0 {
		config.HTTPProxy = func(req *http.Request) (*url.URL, error) {
			return url.Parse(httpProxy[0])
		}
	}

	client, err := torrent.NewClient(config)
	if err != nil {
		return err
	}
	defer client.Close()

	// 解析磁力链接
	tor, err := client.AddMagnet(magnet)
	if err != nil {
		return err
	}
	defer tor.Drop()

	tor.AllowDataDownload()

	fmt.Println("Downloading: getting torrent metadata...")
	<-tor.GotInfo()
	totalBytes := tor.Length()
	//totalPieces := tor.NumPieces()
	//fmt.Println("File Name:", tor.Name())
	//fmt.Printf("totalBytes:%d, totalPieces:%d\n", totalBytes, totalPieces)
	if fileName != "" {
		tor.SetDisplayName(fileName)
	}
	tor.DownloadAll()

TagFor:
	for {
		select {
		case <-ctx.Done():
			tor.Drop()
			break TagFor
		case <-tor.Complete.On():
			fmt.Println()
			fmt.Println("Download finished")
			break TagFor
		default:
			stats := tor.Stats()

			readBytes := stats.BytesRead.Int64()
			showProgress(readBytes, totalBytes, stats.ActivePeers, stats.TotalPeers)
			time.Sleep(time.Millisecond * 400)
		}
	}

	return nil
}

var RateCalculator = ratecalc.NewCalculator()

func showProgress(readBytes int64, totalBytes int64, activePeers, totalPeers int) {
	percent := min(100, 100*float64(readBytes)/float64(totalBytes))
	totalMiBytes := float64(totalBytes) / 1024 / 1024
	readMiBytes := min(totalMiBytes, float64(readBytes)/1024/1024)

	// 计算下载速度
	RateCalculator.Add(readBytes)
	rate := RateCalculator.GetAverageAuto()

	progress := fmt.Sprintf("Downloading: %.2f/%.2fMiB - %.2f%% - %s - Peers(%d)",
		readMiBytes, totalMiBytes, percent, rate, activePeers)
	// 动态打印：
	// 1. 清除当前行
	fmt.Printf("\r%s", strings.Repeat(" ", len(progress)+16))
	// 2. 打印新当前行
	fmt.Printf("\r%s", progress)
}
