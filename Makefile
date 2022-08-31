

build: ./cmd/main.go
	go build -o ./out/fgv.exe ./cmd/main.go

locps: ./cmd/main.go
	go run ./cmd/main.go -i 127.0.0.1

locsql: ./cmd/main.go
	go run ./cmd/main.go -i 127.0.0.1 -p 3306

locred: ./cmd/main.go
	go run ./cmd/main.go -i 127.0.0.1 -p 6379

locbms: ./cmd/main.go
	go run ./cmd/main.go -i 127.0.0.1 -p 445