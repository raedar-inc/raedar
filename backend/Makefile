# note: call scripts from /scripts

BUILDPATH=&{CURDIR}
GOINSTALL=${GO} install
GOCLEAN=${GO} clean
GOGET=${GO} get

default:
	@echo "=============building Local API============="
	docker build -f cmd/app/main.go

up: default
	@echo "=============starting api locally============="
	docker-compose up -d

logs:
	docker-compose logs -f

buildapp:
	@echo "========build the app=========\n"
	go build -o bin/main cmd/app/main.go

runapp:
	@echo "\n"
	@echo "======app running on port:8080 in development......"
	./bin/main

start: buildapp runapp

down:
	docker-compose down

test:
	go test -v -cover ./...

compile:
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 cmd/app/main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 cmd/app/main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 cmd/app/main.go

clean: down
	@echo "=============cleaning up============="
	rm -f api
	docker system prune -f
	docker volume prune -f
