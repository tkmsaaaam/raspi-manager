package clamavreceiver

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receiverhelper"

	"github.com/tkmsaaaam/raspi-manager/receiver/clamavreceiver/internal/metadata"
)

type clamavHandler struct {
	consumer consumer.Metrics
	cancel   context.CancelFunc
	config   *Config
	obsrecv  *receiverhelper.ObsReport
	mb       *metadata.MetricsBuilder
}

func newClamavHandler(consumer consumer.Metrics, cfg *Config, settings receiver.CreateSettings, obsrecv *receiverhelper.ObsReport) (*clamavHandler, error) {
	sh := &clamavHandler{
		consumer: consumer,
		config:   cfg,
		obsrecv:  obsrecv,
		mb:       metadata.NewMetricsBuilder(cfg.MetricsBuilderConfig, settings),
	}
	return sh, nil
}

const (
	dataFormat = "clamav"
)

func (sh *clamavHandler) run(ctx context.Context) error {
	ctx, sh.cancel = context.WithCancel(ctx)
	d, err := time.ParseDuration(sh.config.BufferInterval)
	if err != nil {
		return err
	}

	var hostname string
	var b bool
	hostname, b = os.LookupEnv("HOSTNAME")
	if !b {
		hostname = "localhost"
	}

	ticker := time.NewTicker(d)
TICK:
	for {
		now := pcommon.NewTimestampFromTime(time.Now())
		infectedCount, totalError, time, err := sh.scrape(d, now)
		if err != nil {
			log.Println("read file error", err)
		} else {
			if infectedCount >= 0 {
				sh.mb.RecordClamavInfectedCountDataPoint(now, infectedCount, hostname)
			}
			if totalError >= 0 {
				sh.mb.RecordClamavErrorsCountDataPoint(now, totalError, hostname)
			}
			if time >= 0 {
				sh.mb.RecordClamavScanElapsedTimeDataPoint(now, time, hostname)
			}
		}
		select {
		case <-ticker.C:
			metrics := sh.mb.Emit()
			sh.obsrecv.StartMetricsOp(ctx)
			err := sh.consumer.ConsumeMetrics(ctx, metrics)
			sh.obsrecv.EndMetricsOp(ctx, dataFormat, metrics.DataPointCount(), err)
		case <-ctx.Done():
			break TICK
		}
	}

	return nil
}

func (ch *clamavHandler) scrape(d time.Duration, now pcommon.Timestamp) (int64, int64, float64, error) {
	fp, err := os.Open(ch.config.logFilePath)
	if err != nil {
		return -1, -1, -1, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	var infectedCount int64 = -1
	var totalError int64 = -1
	var elapsedTime float64 = -1
	var date time.Time
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Infected files: ") {
			infectedCount, err = strconv.ParseInt(strings.TrimSpace(strings.Split(line, ":")[1]), 10, 64)
			if err != nil {
				return infectedCount, totalError, elapsedTime, err
			}
		}
		if strings.HasPrefix(line, "Total errors: ") {
			totalError, err = strconv.ParseInt(strings.TrimSpace(strings.Split(line, ":")[1]), 10, 64)
			if err != nil {
				return infectedCount, totalError, elapsedTime, err
			}
		}
		if strings.HasPrefix(line, "Time: ") {
			elapsedTime, err = strconv.ParseFloat(strings.TrimSpace(strings.Split(strings.Split(line, ":")[1], "sec")[0]), 10)
			if err != nil {
				return infectedCount, totalError, elapsedTime, err
			}
		}
		if strings.HasPrefix(line, "End Date:") {
			dateStr := strings.Replace("End Date:   2024:06:23 03:14:12", "End Date:   ", "", 1)
			var e error
			date, e = time.Parse("2006:01:02 15:04:05", dateStr)
			if e != nil {
				return -1, -1, -1, nil
			}
			if now.AsTime().Add(-d).Before(date) {
				return infectedCount, totalError, elapsedTime, err
			}
		}
	}
	return -1, -1, -1, nil
}
