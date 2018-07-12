
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"bufio"
  	"log"
  	"os"
)

type metric struct {
	pathStack []string // lifo stack of name components
	value     float64
}

// Pops names of pathStack to build the flattened name for a metric
func (m *metric) name() string {
	buf := bytes.Buffer{}
	for i := len(m.pathStack) - 1; i >= 0; i-- {
		if buf.Len() > 0 {
			buf.WriteString(".")
		}
		buf.WriteString(m.pathStack[i])
	}
	return buf.String()
}


type taggedMetricMap map[string]metricMap
type metricMap map[string]interface{}

// Parses a raw JSON string into a taggedMetricMap
// Delegates the actual parsing to newTaggedMetricMap(..)
func parseDump(dump string) (taggedMetricMap, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(dump), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json: '%s': %v", dump, err)
	}

	return newTaggedMetricMap(data), nil
}

// Builds a TaggedMetricMap out of a generic string map.
// The top-level key is used as a tag and all sub-keys are flattened into metrics
func newTaggedMetricMap(data map[string]interface{}) taggedMetricMap {
	tmm := make(taggedMetricMap)
	for tag, datapoints := range data {
		mm := make(metricMap)
		for _, m := range flatten(datapoints) {
			mm[m.name()] = m.value
		}
		tmm[tag] = mm
	}
	return tmm
}

// Recursively flattens any k-v hierarchy present in data.
// Nested keys are flattened into ordered slices associated with a metric value.
// The key slices are treated as stacks, and are expected to be reversed and concatenated
// when passed as metrics to the accumulator. (see (*metric).name())
func flatten(data interface{}) []*metric {
	var metrics []*metric

	switch val := data.(type) {
	case float64:
		metrics = []*metric{&metric{make([]string, 0, 1), val}}
	case map[string]interface{}:
		metrics = make([]*metric, 0, len(val))
		for k, v := range val {
			for _, m := range flatten(v) {
				m.pathStack = append(m.pathStack, k)
				metrics = append(metrics, m)
			}
		}
	default:
		log.Printf("I! Ignoring unexpected type '%T' for value %v", val, val)
	}

	return metrics
}

func filterMap(dump string, filters []string) (taggedMetricMap, error) {
	data, err := parseDump(dump)
	if err != nil {
		return nil, err
	}
	
	fmap := make(taggedMetricMap)
	for tag, metrics := range data {
		for _, filter := range filters {
			if val, ok := metrics[filter]; ok {
				fmap[tag][filter] = val
			}
		}
	}

	return fmap, nil
}

func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func main() {
    b, err := ioutil.ReadFile("file.txt") 
    if err != nil {
        fmt.Print(err)
    }

    str := string(b) 

    filters, err := readLines("filters.txt")

    filterMap(str, filters)
}