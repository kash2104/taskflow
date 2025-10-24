module github.com/kash2104/taskflow/worker

go 1.25.1

require (
	github.com/joho/godotenv v1.5.1
	github.com/kash2104/taskflow/master v0.0.0
	go.mongodb.org/mongo-driver v1.17.4
	google.golang.org/grpc v1.75.1
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace github.com/kash2104/taskflow/master => ../master