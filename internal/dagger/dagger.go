package dagger

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

type Container = *dagger.Container

type DaggerClient struct {
	client *dagger.Client
}

/*
DaggerのClientを作成する。
*/
func NewDaggerClient(ctx context.Context, isOutputLog bool) (DaggerClient, error) {
	var client *dagger.Client
	var err error
	if isOutputLog {
		client, err = dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	} else {
		client, err = dagger.Connect(ctx)
	}

	if err != nil {
		return DaggerClient{}, err
	}
	return DaggerClient{
		client: client,
	}, nil
}

/*
sourceDirで「gotestsum --junitfile /report/report.xml -- ./...」を行い、レポートを出力する。
*/
func (c DaggerClient) GoTest(ctx context.Context, sourceDir, outputDir string) error {
	// ローカルのsourceDirのディレクトリを取得
	src := c.client.Host().Directory(sourceDir)

	test := c.client.Container().
		From("golang:1.20.1").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "install", "gotest.tools/gotestsum@v1.9.0"}).
		WithExec([]string{"sh", "-c", "gotestsum --junitfile /report/report.xml -- ./... || touch fail.txt"})

	log, err := test.Stdout(ctx)
	fmt.Println(log)
	if err != nil {
		return err
	}

	_, err = test.Directory("/report").Export(ctx, outputDir)
	if err != nil {
		return err
	}

	// エラーがない=fail.txtが存在するため、testに失敗したとみなす
	_, err = test.Directory("/src").File("fail.txt").Contents(ctx)
	if err == nil {
		return fmt.Errorf("テストが失敗しました")
	}
	return nil
}

// TODO: 下記を参考にしたが動かない
// https://docs.dagger.io/757394/use-service-containers
func (c DaggerClient) GoDoc(ctx context.Context, sourceDir string) error {
	// ローカルのsourceDirのディレクトリを取得
	src := c.client.Host().Directory(sourceDir)

	// create HTTP service container with exposed port 8080
	httpSrv := c.client.Container().
		From("golang:1.20.1").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		// WithExec([]string{"go", "run", "main.go"}).
		WithExec([]string{"go", "install", "golang.org/x/tools/cmd/godoc@latest"}).
		WithExec([]string{"sh", "-c", "godoc -http localhost:8080"}).
		// WithExec([]string{"sh", "-c", "go run main.go > /dev/null"}).
		WithExposedPort(8080)

	// create client container with service binding
	// access HTTP service and print result
	val, err := c.client.Container().
		From("alpine").
		WithServiceBinding("www", httpSrv).
		WithExec([]string{"wget", "-O-", "http://www:8080"}).
		Stdout(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Println(val)
	return nil
}

func (c DaggerClient) GoVulnCheck(ctx context.Context, sourceDir string) error {
	// ローカルのsourceDirのディレクトリを取得
	src := c.client.Host().Directory(sourceDir)

	govulncheck := c.client.Container().
		From("golang:1.20.1").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@v0.0.0-20230309043308-bbc736fc3bc1"}).
		WithExec([]string{"govulncheck", "./..."})

	log, err := govulncheck.Stdout(ctx)
	if err != nil {
		return err
	}
	fmt.Println(log)
	return nil
}

func (c DaggerClient) ImageBuild(ctx context.Context, dockerfilePath string) (*dagger.Container, error) {
	// ローカルのsourceDirのディレクトリを取得
	src := c.client.Host().Directory(dockerfilePath)

	image := c.client.Container().Build(src)
	return image, nil
}

func (c DaggerClient) Push(ctx context.Context, repo string, platformVariants []Container) (string, error) {
	response, err := c.client.Container().
		Publish(ctx, repo, dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
		})
	if err != nil {
		return "", err
	}
	return response, nil
}

func (c DaggerClient) Close() {
	c.client.Close()
}
