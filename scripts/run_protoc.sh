PROTO_PATH="./api/proto/v1"
OUTPUT_PATH="./api/gen/go/v1"
mkdir -p $OUTPUT_PATH

protoc -I $PROTO_PATH \
    --go_out=$OUTPUT_PATH --go_opt=paths=source_relative \
	--go-grpc_out=$OUTPUT_PATH --go-grpc_opt=paths=source_relative \
	$PROTO_PATH/*.proto