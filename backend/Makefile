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

setdevenv:
	@echo "============= setting MODE to DEV_MODE  ============="
	export MODE=DEV_MODE
	@echo MODE is ==> $$MODE

logs:
	docker-compose logs -f

buildapp:
	@echo "========build the app=========\n"
	go build -o bin/main cmd/app/main.go

runapp:
	@echo "\n"
	@echo "======app running on port:8080 in development......"
	./bin/main

start: setdevenv buildapp runapp

down:
	docker-compose down

rundocker:
	docker run -p 8080:8008 -it raedar/backend

runtest:
	go test -v -cover ./...

settestenv:
	@echo "============= setting MODE to TESTING_MODE  ============="
	export MODE=TESTING_MODE

testapp: settestenv runtest

compile:
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 cmd/app/main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 cmd/app/main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 cmd/app/main.go

clean: down
	@echo "=============cleaning up============="
	rm -f api
	docker system prune -f
	docker volume prune -f
