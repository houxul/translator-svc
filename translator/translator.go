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

// func JSON(w http.ResponseWriter, v interface{}) {
// 	buf := &bytes.Buffer{}
// 	enc := json.NewEncoder(buf)
// 	enc.SetEscapeHTML(true)
// 	if err := enc.Encode(v); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	w.Write(buf.Bytes())
// }
