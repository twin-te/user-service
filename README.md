# grpc-ts-template

typescript & grpc のテンプレート。

# 推奨開発環境
docker + vscode を使うことで簡単に開発可能。

1. [RemoteDevelopment](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack)拡張機能をインストール
2. このプロジェクトのフォルダを開く
3. 右下に `Folder contains a Dev Container configuration file. Reopen folder to develop in a container` と案内が表示されるので`Reopen in Container`を選択する。（表示されない場合はコマンドパレットを開き`open folder in container`と入力する）
4. node14の開発用コンテナが立ち上がりVSCodeで開かれる。また、`.devcontainer/docker-compose.yml` に任意のサービスを追加するとvscode起動時に一緒に起動できる（データベース等）。

# npmコマンド一覧

|コマンド|説明|
|:--|:--|
|dev| 開発起動|
|proto|protoファイルから型定義を生成(proto-gen.shを実行している)|
|client|grpcリクエストが送れるCLIを起動|
|test|テストを実行|
|build|distにビルド結果を出力|

# とりあえず動かす
```bash
# 準備
yarn && yarn proto

# 開発鯖立ち上げ
yarn dev


## ----以下別窓---- ##

# gRPCリクエストを送るCLIを立ち上げる
yarn client

# CLIが立ち上がったらリクエストを送る
HelloService@localhost:50051> client.greet({name: 'Twin:te'}, pr)

# レスポンスが返ってくれば成功
HelloService@localhost:50051> 
{
  "text": "hello! SIY1121"
}

```

# GitHub Actions
- `.github/workflows/test.yml` pushされるとテストを実行する
- `.github/workflows/release.yml` GitHub上でリリースをPublishするとDockerImageをビルドし、GHCRにプッシュする。

# 必要変更箇所

1. Dockerfile L17
```dockerfile
LABEL org.opencontainers.image.source https://github.com/twin-te/grpc-ts-template
```
後ろのurlを自分のレポジトリのurlに変更する。（DockerImageとレポジトリの紐付けを行う）

2. .github/workflows/release.yml L17
```yml
run: echo "TAG_NAME=ghcr.io/twin-te/grpc-ts-template:${GITHUB_REF#refs/*/}" >> $GITHUB_ENV
```
TAG_NAME=ghcr.io/twin-te/{自分のレポジトリ名} に変更する（GitHubContainerRegistryにプッシュするときに使う）

3. GitHubのSettings > Secrets > New repository secret で以下の環境変数を登録（GitHubAction用）

|名前|説明|
|:---|:---|
|CR_PAT| GitHubContainerRegistry:write の権限を持ったPersonalAccessToken|
|CR_USER|PersonalAccessTokenを作ったユーザー名|