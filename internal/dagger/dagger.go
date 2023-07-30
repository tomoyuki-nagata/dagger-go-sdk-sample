package dagger

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

type DaggerClientConnector struct {
	ctx context.Context
}

func (d DaggerClientConnector) DefaultConnect() (DaggerClient, error) {
	return d.connect(d.ctx)
}

func (d DaggerClientConnector) K8sConnect(kubeNamespace, daggerEnginName string) (DaggerClient, error) {
	err := setupRemoteEngine(d.ctx, kubeNamespace, daggerEnginName)
	if err != nil {
		return DaggerClient{}, err
	}
	return d.connect(d.ctx)
}

func (DaggerClientConnector) connect(ctx context.Context) (DaggerClient, error) {
	var client *dagger.Client
	var err error
	client, err = dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))

	if err != nil {
		return DaggerClient{}, err
	}
	return DaggerClient{
		client: client,
		ctx:    ctx,
	}, nil
}

func NewDaggerClientConnector(ctx context.Context) DaggerClientConnector {
	return DaggerClientConnector{
		ctx: ctx,
	}
}

type DaggerClient struct {
	client *dagger.Client
	ctx    context.Context
}

type Container = *dagger.Container

/*
sourceDirで「gotestsum --junitfile /report/report.xml -- ./...」を実行し、outputDirにレポートを出力する。
*/
func (c DaggerClient) GoTest(sourceDir, outputDir string) error {
	// ローカルのsourceDirのディレクトリを取得
	src := c.client.Host().Directory(sourceDir)

	test := c.client.Container().
		From("golang:1.20.1").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "install", "gotest.tools/gotestsum@v1.9.0"}).
		WithExec([]string{"sh", "-c", "gotestsum --junitfile /report/report.xml -- ./... || touch fail.txt"})

	log, err := test.Stdout(c.ctx)
	fmt.Println(log)
	if err != nil {
		return err
	}

	_, err = test.Directory("/report").Export(c.ctx, outputDir)
	if err != nil {
		return err
	}

	// エラーがない=fail.txtが存在するため、testに失敗したとみなす
	_, err = test.Directory("/src").File("fail.txt").Contents(c.ctx)
	if err == nil {
		return fmt.Errorf("テストが失敗しました")
	}
	return nil
}

/*
godocサーバーを起動し、静的ページとして出力する。
参考URL: https://docs.dagger.io/757394/use-service-containers
*/

func (c DaggerClient) GoDoc(sourceDir, outputDir string) error {
	// ローカルのsourceDirのディレクトリを取得
	src := c.client.Host().Directory(sourceDir)

	// godocコマンドでドキュメントのウェブサーバーを起動する
	httpSrv := c.client.Container().
		From("golang:1.20.1").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "install", "golang.org/x/tools/cmd/godoc@latest"}).
		WithExec([]string{"sh", "-c", "godoc -http=:6060"}).
		WithExposedPort(6060)

	// godocサーバーのページをwgetを利用してダウンロードする
	godoc := c.client.Container().
		From("ubuntu:22.04").
		WithServiceBinding("godoc", httpSrv). // godocというエイリアスでサービスコンテナを登録
		WithExec([]string{"sh", "-c", "apt-get update && apt-get install -y wget"}).
		WithExec([]string{"sh", "-c", "wget -r -np -N -E -k --page-requisites http://godoc:6060/pkg/sample-app/ || echo 'finish'"}) // wgetのstatus codeが8になり後続の処理が失敗するため、エラーを握りつぶす

	log, err := godoc.Stdout(c.ctx)
	fmt.Println(log)
	if err != nil {
		return err
	}

	//godocサーバーからダウンロードした静的ページをエクスポート
	_, err = godoc.Directory("/godoc:6060").Export(c.ctx, outputDir)
	if err != nil {
		return err
	}

	return nil
}

func (c DaggerClient) GoVulnCheck(sourceDir string) error {
	// ローカルのsourceDirのディレクトリを取得
	src := c.client.Host().Directory(sourceDir)

	govulncheck := c.client.Container().
		From("golang:1.20.1").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithExec([]string{"go", "install", "golang.org/x/vuln/cmd/govulncheck@v0.0.0-20230309043308-bbc736fc3bc1"}).
		WithExec([]string{"govulncheck", "./..."})

	log, err := govulncheck.Stdout(c.ctx)
	if err != nil {
		return err
	}
	fmt.Println(log)
	return nil
}

func (c DaggerClient) ImageBuild(dockerfilePath string) (*dagger.Container, error) {
	// ローカルのsourceDirのディレクトリを取得
	src := c.client.Host().Directory(dockerfilePath)

	image := c.client.Container().Build(src)
	return image, nil
}

func (c DaggerClient) Push(repo string, platformVariants []Container) (string, error) {
	response, err := c.client.Container().
		Publish(c.ctx, repo, dagger.ContainerPublishOpts{
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
