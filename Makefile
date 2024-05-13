build:
	go build -o bin ./cmd/gohuw

run:
	go run ./cmd/gohuw

compile:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/gohuw-amd64.exe ./cmd/gohuw
	GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o bin/gohuw-386.exe ./cmd/gohuw
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/gohuw-amd64-darwin ./cmd/gohuw
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o bin/gohuw-arm64-darwin ./cmd/gohuw
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/gohuw-amd64-linux ./cmd/gohuw
	GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o bin/gohuw-386-linux ./cmd/gohuw

	upx --ultra-brute --lzma bin/gohuw-amd64.exe bin/gohuw-386.exe bin/gohuw-amd64-darwin bin/gohuw-arm64-darwin bin/gohuw-amd64-linux bin/gohuw-386-linux
