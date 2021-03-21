# 作業ログ

[春の入門祭り 🌸 #7 作って学ぶ GraphQL。gqlgen を用いて鉄道データ検索 API 開発入門](https://future-architect.github.io/articles/20200609/) を参考に GraphQL サーバの構築を行った。

```bash
# リポジトリの初期セットアップ
make clean
make init
make install
```

## PostgreSQL

[駅データ.jp](https://ekidata.jp/dl/?p=1) から駅データをダウンロードし DB サーバを構築する。  
取得した csv ファイルは `docker/postgres/init/{company, join, line, station}.csv` に配置する。

```bash
# postgres サービスを開始
docker-compose up --build -d postgres

# table(company, station_join, line, station) の確認
docker-compose exec postgres psql -U user -c 'select * from company limit 3' db
```

## gqlgen

```bash
# add `graph/schema.graphqls`
go mod init && gqlgen init
```
