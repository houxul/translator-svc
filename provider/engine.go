package provider

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"translator/model"
)

type provider func([]string, bool) ([]string, error)

var Engine = newEngine()

func newEngine() *engine {
	words := model.ReadWords()
	return &engine{
		providers: []provider{timeWrapper("lingocloud", lingocloudTranslate),
			timeWrapper("baidu", baiduTranslate),
			timeWrapper("tencent", tencentTranslate)},
		records: words,
	}
}

func timeWrapper(label string, p provider) provider {
	return func(srcs []string, en2zh bool) ([]string, error) {
		startTime := time.Now()
		var result []string
		var err error
		defer func(label string) {
			log.Printf("%s cost(%d) en2zh(%t) srcs(%#v) result(%#v)\n", label, time.Since(startTime).Milliseconds(), en2zh, srcs, result)
		}(label)
		result, err = p(srcs, en2zh)
		return result, err
	}
}

type engine struct {
	index     int
	providers []provider
	records   map[string]string
	mutex     sync.Mutex
}

func (e *engine) Query(srcs []string) ([]string, error) {
	if len(srcs) == 0 {
		return []string{}, nil
	}

	if len(srcs) == 1 {
		if isZh(srcs[0]) {
			return e.queryZh(srcs)
		}

		if isStatement(srcs[0]) {
			return e.queryStatement(srcs)
		}
	}

	return e.queryWord(srcs)
}

func (e *engine) queryWord(srcs []string) ([]string, error) {
	dsts := make([]string, len(srcs))
	missingSrcs := make([]string, 0, len(srcs))
	missingSrcIndex := make(map[string]int, len(srcs))
	for i, src := range srcs {
		dst, ok := e.record(src)
		if ok {
			dsts[i] = dst
			continue
		}
		dsts[i] = srcs[i]
		missingSrcs = append(missingSrcs, src)
		missingSrcIndex[src] = i
	}

	if len(missingSrcs) == 0 {
		return dsts, nil
	}

	provider := e.provider()
	missingDsts, err := provider(missingSrcs, true)
	if err != nil {
		return dsts, err
	}
	e.addRecords(missingSrcs, missingDsts)

	for i, src := range missingSrcs {
		dsts[missingSrcIndex[src]] = missingDsts[i]
	}
	return dsts, nil
}

func (e *engine) queryZh(srcs []string) ([]string, error) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	results := make([]string, 0, len(e.providers))
	for i, p := range e.providers {
		wg.Add(1)
		go func(i int, p provider, srcs []string) {
			out, err := p(srcs, false)
			if err != nil {
				fmt.Println("queryZh error", i, srcs[0], err)
				return
			}
			mu.Lock()
			results = append(results, out...)
			mu.Unlock()
			wg.Done()
		}(i, p, srcs)
	}
	wg.Wait()

	return results, nil
}

func (e *engine) queryStatement(srcs []string) ([]string, error) {
	provider := e.provider()
	return provider(srcs, true)
}

func (e *engine) record(src string) (string, bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	dst, ok := e.records[strings.ToLower(src)]
	return dst, ok
}

func (e *engine) addRecords(srcs []string, dsts []string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for i, src := range srcs {
		e.records[strings.ToLower(src)] = dsts[i]
	}
}

func (e *engine) provider() provider {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.index++
	return e.providers[(e.index)%len(e.providers)]
}

func (e *engine) Close() {
	model.WriteWords(e.records)
}
