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

	// if err := client.GoDoc("./sample-app", "./tmp/godoc"); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// if err := client.GoVulnCheck("./sample-app"); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	resp, err := client.PushByDockerFile("./sample-app", "docker.io", "username", "DOCKERHUB_PASSWORD", "username/sample:1.0")
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
