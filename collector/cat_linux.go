//go:build linux
// +build linux

package collector

import (
	"bufio"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

type catCollector struct {
	info   *prometheus.Desc
	logger *slog.Logger
}

func init() {
	registerCollector("cat", defaultEnabled, NewCatCollector)
}

// logger를 매개변수로 받도록 수정
func NewCatCollector(logger *slog.Logger) (Collector, error) {
	return &catCollector{
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "cat", "info"),
			"Key-value info read from /test.txt",
			[]string{"key", "value"},
			nil,
		),
		logger: logger,
	}, nil
}

func (c *catCollector) Update(ch chan<- prometheus.Metric) error {
	file, err := os.Open("/test.txt")
	if err != nil {
		c.logger.Warn("cannot open /test.txt", "error", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			c.logger.Debug("skipping malformed line", "line", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		ch <- prometheus.MustNewConstMetric(
			c.info,
			prometheus.GaugeValue,
			1,
			key, value,
		)
	}
	return sca
