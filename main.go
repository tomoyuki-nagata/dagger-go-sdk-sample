package main

import (
	"context"
	"dagger-go-sdk/internal/dagger"
	"fmt"
	"os"
)

func main() {
	client, err := dagger.NewDaggerClientConnector(context.Background()).DefaultConnect()
	// client, err := dagger.NewDaggerClientConnector(context.Background()).K8sConnect("default", "dagger-engin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()

	// if err := client.GoTest("./sample-app", "./tmp/test_report"); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	if err := client.GoDoc("./sample-app", "./tmp/godoc"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// if err := client.GoVulnCheck("./sample-app"); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// result, err := repository.DockerLogin(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(result)

	// image, err := client.ImageBuild("./sample-app")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// resp, err := client.Push(fmt.Sprintf("ttl.sh/sample-dagger-%.0f", math.Floor(rand.Float64()*10000000)), []dagger.Container{image})
	// fmt.Println(resp)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}
