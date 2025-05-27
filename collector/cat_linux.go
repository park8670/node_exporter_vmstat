
package collector

import (
	"bufio"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type catCollector struct {
	info *prometheus.Desc
}

func init() {
	registerCollector("cat", defaultEnabled, NewCatCollector)
}

func NewCatCollector() (Collector, error) {
	return &catCollector{
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "cat", "info"),
			"Key-value info read from /test.txt",
			[]string{"key", "value"},
			nil,
		),
	}, nil
}

func (c *catCollector) Update(ch chan<- prometheus.Metric) error {
	file, err := os.Open("/test.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue // skip malformed line
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
	return scanner.Err()
}
