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
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/drillbits/googlecloud-api-gateway-urldecode/genproto/api"
)

const bufSize = 1024 * 1024

func TestEcho(t *testing.T) {
	ctx := context.Background()

	listener := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	pb.RegisterEchoServer(srv, &echoServer{})
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(ctx context.Context, address string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{
				query: "a b",
			},
			want:    "a b",
			wantErr: false,
		},
		{
			args: args{
				query: "a+b",
			},
			want:    "a+b",
			wantErr: false,
		},
		{
			args: args{
				query: "a%20b",
			},
			want:    "a%20b",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := pb.NewEchoClient(conn)
			resp, err := client.EchoQuery(ctx, &pb.EchoQueryRequest{Query: tt.args.query})
			if (err != nil) != tt.wantErr {
				t.Fatal(err)
				t.Errorf("EchoQuery error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := resp.GetDecodedQuery()
			if got != tt.want {
				t.Errorf("EchoQuery = %v, want %v", got, tt.want)
			}
		})
	}
}
