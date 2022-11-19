package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// クライアントの用意、ログを出力するように追加で設定
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// リポジトリをクローン(mainブランチの内容)
	src := client.
		Git("https://github.com/jinwatanabe/go_lambda_cicd").
		Branch("main").Tree()

	if err != nil {
		panic(err)
	}

    // CI/CDの実行環境をGoイメージを使うように設定(Goをテスト/ビルドするため)
    // 環境変数を設定する (今回は必要はないが環境変数を設定することもできる)
	golang := client.Container().From("golang:latest").WithEnvVariable("MESSAGE", "CI/CDが成功しました")

    // コンテナの/appにリポジトリの内容をコピー
    // カレントディレクトリを/app/srcに変更(srcにGoのコードがあるため)
	golang = golang.WithMountedDirectory("/app", src).WithWorkdir("/app/src")

    // コンテナの上でコマンドを実行
    // テストコマンド
    // ビルドコマンド
    // ビルド結果を表示(mainがあることを確認)
		// 環境変数を利用してを表示
	golang = golang.Exec(dagger.ContainerExecOpts{
		Args: []string{"sh", "-c", "go test"},
	}).
		Exec(dagger.ContainerExecOpts{
			Args: []string{"sh", "-c", "env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags='-s -w' -o bin/main handler/main.go"},
		}).
		Exec(dagger.ContainerExecOpts{
			Args: []string{"sh", "-c", "ls -la bin"},
		}).
		Exec(dagger.ContainerExecOpts{
			Args: []string{"sh", "-c", "echo $MESSAGE"},
		})


    // CI/CDに失敗したらここでエラーがでて停止する
	if _, err := golang.ExitCode(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Success!!")
}
