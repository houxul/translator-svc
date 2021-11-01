package provider

import (
	"fmt"
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

func timeWrapper(laber string, p provider) provider {
	return func(srcs []string, en2zh bool) ([]string, error) {
		startTime := time.Now()
		defer func(laber string) {
			fmt.Printf("%s %d\n", laber, time.Now().Sub(startTime).Milliseconds())
		}(laber)
		return p(srcs, en2zh)
	}
}

type engine struct {
	index     int
	providers []provider
	records   map[string]string
	mutex     sync.Mutex
}

func (e *engine) Inquiry(srcs []string) ([]string, error) {
	if len(srcs) == 0 {
		return []string{}, nil
	}

	if isEnWord(srcs[0]) {
		return e.inquiryEn(srcs)
	}
	return e.inquiryZh(srcs)
}

func (e *engine) inquiryEn(srcs []string) ([]string, error) {
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

func (e *engine) inquiryZh(srcs []string) ([]string, error) {
	provider := e.provider()
	return provider(srcs, false)
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
