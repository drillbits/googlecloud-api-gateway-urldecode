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
	"crypto/tls"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/drillbits/googlecloud-api-gateway-urldecode/genproto/api"
)

func main() {
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})

	args := os.Args
	if len(args) < 3 {
		log.Fatalf("usage: %s URL QUERY", args[0])
	}

	addr := args[1]
	query := args[2]

	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(cred),
	)
	if err != nil {
		log.Fatalf("failed to connect gRPC server: %s", err)
	}
	defer conn.Close()

	client := pb.NewEchoClient(conn)
	res, err := client.EchoQuery(context.Background(), &pb.EchoQueryRequest{
		Query: query,
	})
	if err != nil {
		log.Fatalf("failed to echo query: %s", err)
	}

	log.Printf("decoded query: %#v", res.GetDecodedQuery())
}
