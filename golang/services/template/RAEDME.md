# Template
以下を用いた一般的なgRPCサーバーのサンプルテンプレート
- gorm (ORM library)

## development(using docker-compose)

- サービス立ち上げ

```shell
docker compose up -d
```

- DBに入る

```shell
docker exec -it template_service_database mysql -u root -p
```

passwordは`root`

## build for prd
```shell
docker build --target prd -t $IMAGE_NAME:$TAG .
```


### おまけ
gRPCクライアントのサンプルコード
- `cmd/grpcclient/main.go`

grpcurlのサンプル
```shell
grpcurl -plaintext -d '{"message":"hoge"}' localhost:8080 list ping.PingService.Ping
```
