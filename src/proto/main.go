/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package api

import (
	"log"
	"os"
	"unsafe"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/golang/protobuf/proto"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	hr := &pb.HelloRequest{Name: name}
	data, err := proto.Marshal(hr)
	var pbr  proto.Message
	err = proto.Unmarshal(data, pbr)
	// var ppbr unsafe.Pointer
	ppbr := unsafe.Pointer(&pbr)
	// r, err := c.SayHelloAgain(context.Background(), &pb.HelloRequest{Name: name})
	r, err := c.SayHelloAgain(context.Background(), (*pb.HelloRequest)(ppbr))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
