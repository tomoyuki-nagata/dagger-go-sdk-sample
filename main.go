package main

import (
	"context"
	"dagger-go-sdk/internal/dagger"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	// repository.SetupRemoteEngine(ctx)

	client, err := dagger.NewDaggerClientConnector(ctx).K8sConnect("default", "dagger-engin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()

	if err := client.GoTest(ctx, "./sample-app", "./tmp/test_report"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// if err := client.GoDoc(ctx, "./sample-app"); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// if err := client.GoVulnCheck(ctx, "./sample-app"); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// if err := repository.DockerLogin(ctx); err != nil {
	// 	fmt.Println(err)
	// }

	// image, err := client.ImageBuild(ctx, "./sample-app")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// resp, err := client.Push(ctx, fmt.Sprintf("ttl.sh/sample-dagger-%.0f", math.Floor(rand.Float64()*10000000)), []dagger.Container{image})
	// fmt.Println(resp)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}
