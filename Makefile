default: build-whoami-local

clean:
	rm -rfv target

build-whoami-local:
	go build -a -v -o target/whoami

build-whoami:
	GOOS='linux' GOARCH='amd64' CGO_ENABLED=0 go build -a -v -ldflags '-s' -o target/whoami-linux-64

build-docker: build-whoami
	docker build -t jeschu/whoami:latest .

run-docker: build-docker
	docker run -it --rm -p 8888:80 jeschu/whoami:latest

push-docker: build-docker
	docker push jeschu/whoami:latest