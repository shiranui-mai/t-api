package main

import (
	"io"
	"os"
	"strings"
	. "fmt"
	"net/http"
)

var ss = map[string]string{"api":"Test"}

func check_group(g string) bool {
	_,ok := ss[g]
	return ok
}
func test_redirect(w http.ResponseWriter, r *http.Request) {
	Printf("[test_redirect] r: %+v\n", *r)
	io.WriteString(w, Sprintf("r: %+v\n", *r))
}

func api(w http.ResponseWriter, r *http.Request) {
	Printf("[api] r: %+v\n", *r)
	io.WriteString(w, Sprintf("r: %+v\n", *r))
	// io.WriteString(w, Sprintf("Path: %+v\n", r.URL.Path))
	gs := strings.Split(r.URL.Path,"/")
	// io.WriteString(w, Sprintf("gs: %+v\n", gs))

	group_id := gs[1]
	// io.WriteString(w, Sprintf("group_id: %+v\n", group_id))
	if check_group(group_id) {
		Println("Success Group id: ", group_id);
		http.Redirect(w, r, "/Test_redirect", 0)
		io.WriteString(w, Sprintf("Success Group id:  %+v\n", group_id))
	} else {
		Println("Not found Group id: ", group_id);
		io.WriteString(w, Sprintf("Not found Group id:  %+v\n", group_id))
	}
}

func main() {
	Println("Api Success Start", os.Args[0])
	http.HandleFunc("/", api)
	http.HandleFunc("/Test_Redirect", test_redirect)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic(err)
	}
}
