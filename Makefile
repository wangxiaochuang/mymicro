all:
	PORT=1000 go run github.com/wxc/demo

test:
	cd micro && go test -v github.com/wxc/micro/config/...
