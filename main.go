package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"gopkg.in/yaml.v2"
)

type TestConfig struct {
	Hostname string
	Tests    map[string]int
}

func main() {
	var config TestConfig
	var wg sync.WaitGroup
	source, err := ioutil.ReadFile("tests.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}

	// avoid lookup
	prefix := config.Hostname
	fmt.Printf("Begin status checking of: %v\n\n", prefix)
	wg.Add(len(config.Tests))
	for path, status_code := range config.Tests {
		go func(prefix, path string, status_code int) {
			defer wg.Done()
			resp, _ := http.Get(prefix + path)
			fmt.Printf("Checked path: %v\nValue: %v\nExpectation: %v\n\n", path, resp.StatusCode, status_code)
		}(prefix, path, status_code)
	}

	wg.Wait()
}
