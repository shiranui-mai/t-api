package main

import (
	"io"
	"os"
	"strings"
	. "fmt"
	"net/url"
	"net/http"
	"encoding/json"
	"proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	_ "github.com/golang/protobuf/proto"
)

var ss = map[string]string{"api":"localhost:12121"}

func check_group(g string) bool {
	_,ok := ss[g]
	return ok
}
func check_error(err error) {
	if err != nil {
		panic(err)
	}
}
func test_redirect(w http.ResponseWriter, r *http.Request) {
	Printf("[test_redirect] r: %+v\n", *r)
	io.WriteString(w, Sprintf("r: %+v\n", *r))
}

func get_http_content(r *http.Request) url.Values {
	if r.Method == "GET" {
		return r.URL.Query()
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		check_error(err)
		return r.PostForm
	}
	return nil
}

func api_(w http.ResponseWriter, r *http.Request) {
	Printf("[api] r: %+v\n", *r)
	io.WriteString(w, Sprintf("r: %+v\n", *r))
	// io.WriteString(w, Sprintf("Path: %+v\n", r.URL.Path))
	gs := strings.Split(r.URL.Path,"/")
	// io.WriteString(w, Sprintf("gs: %+v\n", gs))

	group_id := gs[1]
	// io.WriteString(w, Sprintf("group_id: %+v\n", group_id))
	if check_group(group_id) {
		Println("Success Group id: ", group_id);
		// http.Redirect(w, r, "/Test_redirect", 0)
		content := get_http_content(r)
		io.WriteString(w, Sprintf("Success Group id:  %+v\n", group_id))
		io.WriteString(w, Sprintf("content:  %+v\n", content))
		js_content,err := json.Marshal(content)
		check_error(err)
		io.WriteString(w, Sprintf("json content:  %+v\n", string(js_content[:])))

		conn, err := grpc.Dial(ss[group_id], grpc.WithInsecure())
		check_error(err)
		defer conn.Close()
		c := api.NewApiServerClient(conn)

		hr := &api.HeartbeatReq{}
		r, err := c.Heartbeat(context.Background(), hr)
		check_error(err)
		Printf("r: %+v\n", r);
		io.WriteString(w, Sprintf("r:  %+v\n", *r))

	} else {
		Println("Not found Group id: ", group_id);
		io.WriteString(w, Sprintf("Not found Group id:  %+v\n", group_id))
	}
}

func main() {
	Println("Api Success Start", os.Args[0])
	http.HandleFunc("/", api_)
	http.HandleFunc("/Test_Redirect", test_redirect)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic(err)
	}
}
