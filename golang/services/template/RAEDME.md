# Template
以下を用いた一般的なHTTPサーバーのサンプルテンプレート
- gin (http framework)
- gorm (ORM library)

### development(using docker-compose)

- サービス立ち上げ

```shell
docker-compose up -d
```

- DBに入る

```shell
docker exec -it template_service_database mysql -u root -p
```

passwordは`root`

### build for prd
```shell
docker build --target prd -t $IMAGE_NAME:$TAG .
```