# 作業ログ

[春の入門祭り 🌸 #7 作って学ぶ GraphQL。gqlgen を用いて鉄道データ検索 API 開発入門](https://future-architect.github.io/articles/20200609/) を参考に GraphQL サーバの構築を行った。

```bash
# リポジトリの初期セットアップ
make clean
make init
make install
```

## 機能追加作業フロー

**DB 更新**

1. `docker/postgres/init` に csv ファイルを配置する
2. `docker/postgres/init` の sql ファイルにテーブル作成コードを追加する
3. `make db DB_OPTION="--build"` で DB サーバを起動する
4. `docker-compose exec postgres psql -U user -c 'select * from company limit 3' db` でデータが格納されているか動作確認する

**schema 更新**

1. `graph/schema.graphqls` を更新する
2. `make generate` で自動生成コードを更新する
3. `xo` コマンドでクエリの自動生成コードを生成する
4. resolver (データマッピングのグルーコード)を実装する
5. GraphQL Playground で動作確認する

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

# generate codes after updating schema
gqlgen generate
```

## xo

```bash
# setup db server
make db
mkdir -p models

xo 'pgsql://user:password@localhost:5432/db?sslmode=disable' -N -M -B -T StationConn -o models/ << ENDSQL
select li.line_name,
       li.line_name_h,
       li.line_cd,
       st.station_cd,
       st.station_g_cd,
       st.address,
       st.station_name,
       COALESCE(s2l.line_name, '')     as before_line_name,
       COALESCE(st2.station_cd, 0)    as before_station_cd,
       COALESCE(st2.station_name, '') as before_station_name,
       COALESCE(st2.address, '')      as before_address,
       COALESCE(s3l.line_name, '')     as after_line_name,
       COALESCE(st3.station_cd, 0)    as after_station_cd,
       COALESCE(st3.station_name, '') as after_station_name,
       COALESCE(st2.address, '')      as after_address,
       COALESCE(gli.line_name, '')    as transfer_line_name,
       COALESCE(gs.station_cd, 0)     as transfer_station_cd,
       COALESCE(gs.station_name, '')  as transfer_station_name,
       COALESCE(gs.address, '')       as transfer_address
from station st
         inner join line li on st.line_cd = li.line_cd
         left outer join station_join sjb on st.line_cd = sjb.line_cd and st.station_cd = sjb.station_cd2
         left outer join station_join sja on st.line_cd = sja.line_cd and st.station_cd = sja.station_cd1
         left outer join station st2 on sjb.line_cd = st2.line_cd and sjb.station_cd1 = st2.station_cd
         left outer join line s2l on st2.line_cd = s2l.line_cd
         left outer join station st3 on sja.line_cd = st3.line_cd and sja.station_cd2 = st3.station_cd
         left outer join line s3l on st3.line_cd = s3l.line_cd
         left outer join station gs on st.station_g_cd = gs.station_g_cd and st.station_cd <> gs.station_cd
         left outer join line gli on gs.line_cd = gli.line_cd
where st.station_cd = %%stationCD int%%
  and st.e_status = 0
order by st.e_sort
ENDSQL

# 駅CD検索
xo 'pgsql://user:password@localhost:5432/db?sslmode=disable' -N -M -B -T StationByCD -o models/ << ENDSQL
select l.line_cd, l.line_name, s.station_cd, station_g_cd, s.station_name, s.address
from station s
         inner join line l on s.line_cd = l.line_cd
where s.station_cd = %%stationCD int%%
  and s.e_status = 0
ENDSQL

# 乗り換え検索
# 乗換駅検索
xo 'pgsql://user:password@localhost:5432/db?sslmode=disable' -N -M -B -T Transfer -o models/ << ENDSQL
select s.station_cd,
       ls.line_cd,
       ls.line_name,
       s.station_name,
       s.station_g_cd,
       s.address,
       COALESCE(lt.line_cd, 0)     as transfer_line_cd,
       COALESCE(lt.line_name, '')   as transfer_line_name,
       COALESCE(t.station_cd, 0)   as transfer_station_cd,
       COALESCE(t.station_name, '') as transfer_station_name,
       COALESCE(t.address, '')      as transfer_address
from station s
         left outer join station t on s.station_g_cd = t.station_g_cd and s.station_cd <> t.station_cd
         left outer join line ls on s.line_cd = ls.line_cd
         left outer join line lt on t.line_cd = lt.line_cd
where s.station_cd = %%stationCD int%%
ENDSQL

# 駅名検索
xo 'pgsql://user:password@localhost:5432/db?sslmode=disable' -N -M -B -T StationByName -o models/ << ENDSQL
select l.line_cd, l.line_name, s.station_cd, station_g_cd, s.station_name, s.address
from station s
         inner join line l on s.line_cd = l.line_cd
where s.station_name = %%stationName string%%
  and s.e_status = 0
ENDSQL

# 製品ID検索
xo 'pgsql://user:password@localhost:5432/db?sslmode=disable' -N -M -B -T ProductById -o models/ << ENDSQL
select p.product_id,p.company_cd,p.inventory_cd,p.price_jpy,p.product_name,p.description
from product p
where p.product_id = %%productId int%%
ENDSQL

# 製品 CompanyCd 検索
xo 'pgsql://user:password@localhost:5432/db?sslmode=disable' -N -M -B -T ProductByCompanyCd -o models/ << ENDSQL
select p.product_id,p.company_cd,p.inventory_cd,p.price_jpy,p.product_name,p.description
from product p
where p.company_cd = %%companyCd int%%
ENDSQL
```

## GraphQL playground

```bash
# Station
fragment stationF on Station {
  lineName
  stationCD
  stationName
}

query stations {
  stationByCD(stationCD: 1130201) {
    ...stationF
    transferStation {
      ...stationF
    }
  }
  stationByName(stationName: "大崎") {
    ...stationF
  }
}
---
{
  "data": {
    "stationByCD": {
      "lineName": "JR山手線",
      "stationCD": 1130201,
      "stationName": "大崎",
      "transferStation": [
        {
          "lineName": "JR埼京線",
          "stationCD": 1132101,
          "stationName": "大崎"
        },
        {
          "lineName": "JR湘南新宿ライン",
          "stationCD": 1133307,
          "stationName": "大崎"
        },
        {
          "lineName": "りんかい線",
          "stationCD": 9933708,
          "stationName": "大崎"
        }
      ]
    },
    "stationByName": [
      {
        "lineName": "JR山手線",
        "stationCD": 1130201,
        "stationName": "大崎"
      },
      {
        "lineName": "JR埼京線",
        "stationCD": 1132101,
        "stationName": "大崎"
      },
      {
        "lineName": "JR湘南新宿ライン",
        "stationCD": 1133307,
        "stationName": "大崎"
      },
      {
        "lineName": "りんかい線",
        "stationCD": 9933708,
        "stationName": "大崎"
      }
    ]
  }
}

# Product
fragment productF on Product {
  productId
  companyCd
  inventoryCd
  priceJpy
  productName
  description
}

query products {
  productById(productId: 0){
    ...productF
  }
  productByCompanyCd(companyCd: 0) {
    ...productF
  }
}
---
{
  "data": {
    "productById": {
      "productId": 0,
      "companyCd": 0,
      "inventoryCd": 0,
      "priceJpy": 10000,
      "productName": "product-0",
      "description": "Hello world"
    },
    "productByCompanyCd": [
      {
        "productId": 0,
        "companyCd": 0,
        "inventoryCd": 0,
        "priceJpy": 10000,
        "productName": "product-0",
        "description": "Hello world"
      },
      {
        "productId": 3,
        "companyCd": 0,
        "inventoryCd": 0,
        "priceJpy": 40000,
        "productName": "product-3",
        "description": "Bonjour le monde"
      }
    ]
  }
}
```
