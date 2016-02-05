package json

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/influxdata/telegraf"

	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/internal/encoding"
)

type JsonParser struct {
}

func (p *JsonParser) Parse(buf []byte) ([]telegraf.Metric, error) {

	metrics := make([]telegraf.Metric, 0)

	var jsonOut interface{}
	err := json.Unmarshal(buf, &jsonOut)
	if err != nil {
		err = fmt.Errorf("unable to parse out as JSON, %s", err)
		return nil, err
	}

	f := internal.JSONFlattener{}
	err = f.FlattenJSON("", jsonOut)
	if err != nil {
		return nil, err
	}

	metric, err := telegraf.NewMetric("exec", nil, f.Fields, time.Now().UTC())

	if err != nil {
		return nil, err
	}
	return append(metrics, metric), nil
}

func (p *JsonParser) ParseLine(line string) (telegraf.Metric, error) {
	metrics, err := p.Parse([]byte(line + "\n"))

	if err != nil {
		return nil, err
	}

	if len(metrics) < 1 {
		return nil, fmt.Errorf("Can not parse the line: %s, for data format: influx ", line)
	}

	return metrics[0], nil
}

func NewParser() *JsonParser {
	return &JsonParser{}
}

func (p *JsonParser) InitConfig(configs map[string]interface{}) error {
	return nil
}

func init() {
	encoding.Add("json", func() encoding.Parser {
		return NewParser()
	})
}
