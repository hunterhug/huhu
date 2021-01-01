GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -x -o zhihu_linux_amd64-$1 main.go
GOOS=linux GOARCH=386 go build -ldflags "-s -w" -x -o zhihu_linux_386-$1 main.go
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -x -o zhihu_windows_amd64-$1.exe main.go
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -x -o zhihu_windows_386-$1.exe main.go
# Mac
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -x -o zhihu_mac-$1 main.go

mkdir ./binary_zhihu
mkdir ./binary_zhihu/data
mv zhihu_* ./binary_zhihu