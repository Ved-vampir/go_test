
package main

import (
    "fmt"
    "io/ioutil"
	"bufio"
  	"log"
  	"os"
 )

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



func filterMap(dump string, filters []string) (taggedMetricMap, error) {
	data, err := parseDump(dump)
	if err != nil {
		return nil, err
	}
	
	fmap := make(map[string]interface{})
	for filter := range filters {
		fmap[filter] := data[filter]
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