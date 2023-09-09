# Workspace

## Dependencies
- [asdf](https://asdf-vm.com/) v0.12.0
- [docker](https://www.docker.com/) 24.0.5

## Development
- setup
```shell
make setup
```

- teardown
```shell
make teardown
```

- test
```shell
make test
```

- code gen(protoc ...etc)
```shell
make gen
```

- how to add plugin
1. `Makefile`の`add-plugin:`に追加する
2. 指定バージョンをインストール `asdf install <plugin> <version>`
3. `.tool-versions`に反映 `asdf local <plugin> <version>`