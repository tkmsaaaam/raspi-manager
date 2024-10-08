package clamavreceiver

import (
	"bufio"
	"context"
	"fmt"
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

func newClamavHandler(consumer consumer.Metrics, cfg *Config, settings receiver.Settings, obsrecv *receiverhelper.ObsReport) (*clamavHandler, error) {
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
		n := time.Now()
		now := pcommon.NewTimestampFromTime(n)
		fp, err := os.Open(sh.config.LogFilePath)
		if err != nil {
			log.Println("read file error: ", err)
		}
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		infectedCount, totalError, time, err := scrape(scanner, d, n)
		if err != nil {
			log.Println("scrape failed: ", err)
		} else {
			log.Printf("now: %s, infected: %d, error: %d, time: %f\n", now.AsTime().GoString(), infectedCount, totalError, time)
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

func scrape(scanner *bufio.Scanner, d time.Duration, now time.Time) (int64, int64, float64, error) {
	var b bool = false
	var err error
	var infectedCount int64 = -1
	var totalError int64 = -1
	var elapsedTime float64 = -1
	var date time.Time
	for scanner.Scan() {
		line := scanner.Text()
		if line == "----------- SCAN SUMMARY -----------" {
			b = true
		}
		if !b {
			continue
		}
		if strings.HasPrefix(line, "Infected files: ") {
			infectedCount, err = strconv.ParseInt(strings.TrimSpace(strings.Split(line, ":")[1]), 10, 64)
			if err != nil {
				return -1, -1, -1, fmt.Errorf("infected files value invalid: %s", err)
			}
		}
		if strings.HasPrefix(line, "Total errors: ") {
			totalError, err = strconv.ParseInt(strings.TrimSpace(strings.Split(line, ":")[1]), 10, 64)
			if err != nil {
				return -1, -1, -1, fmt.Errorf("total errors value invalid: %s", err)
			}
		}
		if strings.HasPrefix(line, "Time: ") {
			elapsedTime, err = strconv.ParseFloat(strings.TrimSpace(strings.Split(strings.Split(line, ":")[1], "sec")[0]), 10)
			if err != nil {
				return -1, -1, -1, fmt.Errorf("time value invalid: %s", err)
			}
		}
		if strings.HasPrefix(line, "End Date:") {
			dateStr := strings.Replace(line, "End Date:   ", "", 1)
			date, err = time.Parse("2006:01:02 15:04:05", dateStr)
			if err != nil {
				return -1, -1, -1, fmt.Errorf("time value invalid: %s", err)
			}
			b = false
			if now.Add(-d).Before(date) {
				return infectedCount, totalError, elapsedTime, nil
			} else {
				infectedCount = -1
				totalError = -1
				elapsedTime = -1
				err = nil
			}
		}
	}
	return -1, -1, -1, nil
}
