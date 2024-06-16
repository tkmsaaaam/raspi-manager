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

	ticker := time.NewTicker(d)
TICK:
	for {
		now := pcommon.NewTimestampFromTime(time.Now())
		i, err := sh.scrape()
		if err != nil {
			log.Panicln("read file error", err)
		} else if i >= 0 {
			sh.mb.RecordClamavInfectedCountDataPoint(now, i, "host")
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

func (ch *clamavHandler) scrape() (int64, error) {
	fp, err := os.Open(ch.config.logFilePath)
	if err != nil {
		return -1, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "Infected files: ") {
			i, err := strconv.ParseInt(strings.TrimSpace(strings.Split(scanner.Text(), ":")[1]), 10, 64)
			if err != nil {
				return 0, err
			}
			return i, nil
		}
	}
	return -1, nil
}
