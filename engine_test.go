package engine_test

import (
	"context"
	"fmt"
	"github.com/autom8ter/engine"
	"github.com/autom8ter/engine/config"
	"github.com/autom8ter/engine/examples/examplepb/client"
	"github.com/autom8ter/util"
	"github.com/grpc-ecosystem/grpc-gateway/examples/proto/examplepb"
	"github.com/spf13/viper"
	"net/http"
	"testing"
)

func TestGRPC(t *testing.T) {
	var eng = engine.New().With(
		config.WithGRPCListener("tcp", ":3000"),
		config.WithPluginPaths("bin/example.plugin"),
		config.WithPluginSymbol("Plugin"),
		config.WithEnvPrefix("ENGINE"),
	)
	go eng.ServeGRPC()
	var grpcCli = client.ExampleClient(viper.GetString("address"))
	resp, err := grpcCli.EchoBody(context.Background(), &examplepb.SimpleMessage{
		Id:  "yoyoyoyoyo",
		Num: 199,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("GRPC RESPONSE")
	fmt.Println(util.ToPrettyJsonString(resp))
}

func TestHTTP(t *testing.T) {
	var eng = engine.New().With(
		config.WithPluginPaths("bin/example.plugin"),
		config.WithPluginSymbol("Plugin"),
		config.WithEnvPrefix("ENGINE"),
	)
	go http.ListenAndServe(":3001", eng)
	var grpcCli = client.ExampleClient(viper.GetString("address"))
	gresp, err := grpcCli.EchoBody(context.Background(), &examplepb.SimpleMessage{
		Id:  "yoyoyoyoyo",
		Num: 199,
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("GRPC RESPONSE")
	fmt.Println(util.ToPrettyJsonString(gresp))
	resp, err := http.Get("http://0.0.0.0:3001")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("REST RESPONSE")
	fmt.Println(util.ToPrettyJsonString(resp.Status))
}
