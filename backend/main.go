// Copyright 2021 drillbits
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/drillbits/googlecloud-api-gateway-urldecode/genproto/api"
)

var (
	_ pb.EchoServer = (*echoServer)(nil)
)

type echoServer struct{}

func (s *echoServer) EchoQuery(_ context.Context, req *pb.EchoQueryRequest) (*pb.EchoQueryResponse, error) {
	return &pb.EchoQueryResponse{
		DecodedQuery: req.Query,
	}, nil
}

func main() {
	srv := grpc.NewServer()
	reflection.Register(srv)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	addr := net.JoinHostPort("0.0.0.0", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	pb.RegisterEchoServer(srv, &echoServer{})

	log.Printf("service listening on tcp: %s", listener.Addr().String())
	srv.Serve(listener)
}
