package downloader

import (
	"animedown/util/ratecalc"
	"context"
	"errors"
	alog "github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Stage string

const (
	StagePrepare  Stage = "Prepare"
	StageDownload Stage = "Downloading"
	StageSuccess  Stage = "Success"
	StageCancel   Stage = "Cancel"
	StageError    Stage = "Error"
)

type Observer func(stage Stage, readBytes int64, totalBytes int64, activePeers, totalPeers int,
	downloadRate string, percentage float64)

type DownloadConfig struct {
	Dir       string          // Optional, default: .
	FileName  string          // Optional
	Magnet    string          // Must
	Proxy     string          // Optional
	Ctx       context.Context // Optional, you can cancel the download task by this context
	Obs       Observer        // Optional, you can get the download progress by this observer
	ObsPeriod int             // Optional, default: 400 ms
}

func DownloadBlocked(cfg *DownloadConfig) error {
	if err := reviseConfig(cfg); err != nil {
		return err
	}
	_ = os.MkdirAll(cfg.Dir, os.ModePerm)

	config := torrent.NewDefaultClientConfig()
	config.DataDir = cfg.Dir
	config.Logger = config.Logger.WithFilterLevel(alog.Error)
	config.Logger.SetHandlers()
	if cfg.Proxy != "" {
		config.HTTPProxy = func(req *http.Request) (*url.URL, error) {
			return url.Parse(cfg.Proxy)
		}
	}

	client, err := torrent.NewClient(config)
	if err != nil {
		return err
	}
	defer client.Close()

	// 解析磁力链接
	tor, err := client.AddMagnet(cfg.Magnet)
	if err != nil {
		return err
	}
	cfg.Obs(StagePrepare, 0, 0, 0, 0, "", 0)
	defer dropTorrentSafely(tor)
	tor.AllowDataDownload()
	<-tor.GotInfo()
	totalBytes := tor.Length()
	//totalPieces := tor.NumPieces()
	//fmt.Println("File Name:", tor.Name())
	//fmt.Printf("totalBytes:%d, totalPieces:%d\n", totalBytes, totalPieces)
	if cfg.FileName != "" {
		tor.SetDisplayName(cfg.FileName)
	}
	tor.DownloadAll()
	calc := ratecalc.NewCalculator()

TagFor:
	for {
		select {
		case <-cfg.Ctx.Done():
			cfg.Obs(StageCancel, 0, 0, 0, 0, "", 0)
			break TagFor
		case <-tor.Complete.On():
			cfg.Obs(StageSuccess, totalBytes, totalBytes, 0, 0, "0B/s", 100)
			break TagFor
		default:
			stats := tor.Stats()

			readBytes := stats.BytesRead.Int64()
			percentage := min(100, 100*float64(readBytes)/float64(totalBytes))
			calc.Add(readBytes)
			cfg.Obs(StageDownload, readBytes, totalBytes, stats.ActivePeers, stats.TotalPeers, calc.GetAverageAuto(), percentage)

			time.Sleep(time.Millisecond * 400)
		}
	}

	return nil
}

func reviseConfig(cfg *DownloadConfig) error {
	if cfg == nil {
		return errors.New("config is nil")
	}
	if cfg.Dir == "" {
		cfg.Dir = "."
	}
	if cfg.Magnet == "" {
		return errors.New("magnet is empty")
	}
	if cfg.Ctx == nil {
		cfg.Ctx = context.Background()
	}
	if cfg.Obs == nil {
		cfg.Obs = func(stage Stage, readBytes int64, totalBytes int64, activePeers, totalPeers int,
			downloadRate string, percentage float64) {
		}
	}
	if cfg.ObsPeriod <= 0 {
		cfg.ObsPeriod = 400
	}
	var err error
	cfg.Dir, err = filepath.Abs(cfg.Dir)
	return err
}

func dropTorrentSafely(tor *torrent.Torrent) {
	defer func() {
		recover()
	}()
	if tor != nil {
		tor.Drop()
	}
}
