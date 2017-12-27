package etcdserverpb

import (
	"fmt"
	"strings"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	_ "github.com/golang/protobuf/proto"
)

const (
	ETCDHost = "172.17.0.1:2379"
)

var conn *grpc.ClientConn
var client *MaintenanceClient

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}
func Init() {
	var err error
	conn, err = grpc.Dial(ETCDHost, grpc.WithInsecure())
	check_error(err)
	// defer conn.Close()
	c := NewMaintenanceClient(conn)
	client = &c
}

func CheckVersion() {
	sr := StatusRequest{}
	res,err := (*client).Status(context.Background(), &sr)
	check_error(err)

	if !strings.HasPrefix(res.Version, "3.") {
		panic(fmt.Sprintf("The version(%+v) has not 3.*", res.Version));
	}

	fmt.Printf("[check_version] res: %+v\n", *res)
}
