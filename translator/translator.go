package translator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"translator/provider"
)

func Translate(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	bs, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll %v\n", err)
		fmt.Fprintf(w, "%v", err)
		return
	}
	var srcs []string
	if err := json.Unmarshal(bs, &srcs); err != nil {
		fmt.Printf("json.Unmarshal %v\n", err)
		fmt.Fprintf(w, "%v", err)
		return
	}

	result, err := provider.Engine.Inquiry(srcs)
	if err != nil {
		fmt.Printf("Engine.Inquiry %v\n", err)
	}

	bs, err = json.Marshal(result)
	fmt.Fprintf(w, "%s", bs)
}
