# protoc-gen-fake-data-ts [WIP]

## Usage

```
$ go build main.go
$ protoc --proto_path=./path/to/proto/files --plugin=protoc-gen-fake-data-ts=./main --fake-data-ts_out=./path/to/output/files ./path/to/proto/files/*.proto
```
