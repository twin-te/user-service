[![test](https://github.com/twin-te/user-service/actions/workflows/test.yml/badge.svg)](https://github.com/twin-te/timetable-service/actions/workflows/test.yml)

# twinte-user-service
時間割アプリ Twin:te - https://app.twinte.net のv3バックエンドの一部です。

ユーザー情報を管理します。

# 利用方法
[ビルド済みDockerImage](https://github.com/orgs/twin-te/packages?repo_name=user-service)が利用できます。

| 環境変数名  | 説明                             | default               |
|------------|----------------------------------|-----------------------|
| PG_HOST     | Postgres接続先のホスト名         | postgres              |
| PG_PORT     | Postgres接続先のポート番号       | 5432                  |
| PG_DATABASE | Postgres接続先のデータベース名   | twinte_user_service |
| PG_USER     | Postgres接続に使用するユーザー名 | postgres              |
| PG_PASSWORD | Postgres接続に使用するパスワード | postgres              |
| LOG_LEVEL   | ログレベル fatal / error / warn / info / debug / trace / off | info              |

# 開発方法
Docker + VSCodeを推奨します。
以下その方法を紹介します。

1. [RemoteDevelopment](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.vscode-remote-extensionpack)拡張機能をインストール
2. このプロジェクトのフォルダを開く
3. 右下に `Folder contains a Dev Container configuration file. Reopen folder to develop in a container` と案内が表示されるので`Reopen in Container`を選択する。（表示されない場合はコマンドパレットを開き`open folder in container`と入力する）
4. node14の開発用コンテナが立ち上がりVSCodeで開かれます。また、別途postgresも立ち上がり利用できるようになります。
5. `yarn install` で依存をインストールします。
6. `yarn proto` でgrpcに必要なファイルを生成します（開発中にprotoを変更した際も実行してください）
7. `yarn dev` で立ち上がります。

また、`yarn test` でテストを実行、`yarn build` でビルドできます。

`yarn client`を実行するとcliでgrpcリクエストを送れる[grpcc](https://github.com/njpatel/grpcc)が利用できます。

# v3バックエンドサービス一覧
 - [API Gateway](https://github.com/twin-te/api-gateway)
 - Auth Callback
 - **User Service**
 - Session Service
 - [Timetable Service](https://github.com/twin-te/timetable-service)
 - [Course Service](https://github.com/twin-te/course-service)
 - Search Service
 - Donation Service
 - School Calendar Service
 - Information Service
 - Task Service