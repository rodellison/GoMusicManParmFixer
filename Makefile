.PHONY: build clean deploy gomodgen

buildForAWS:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/main main.go

buildForOSX:
	export GO111MODULE=on
	env GOOS=darwin go build -ldflags="-s -w" -o bin/main main.go

buildForWin:
	set GO111MODULE=on
	set GOOS=windows go build -ldflags="-s -w" -o bin/main main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

cleanWin:
	del /s  /q bin\*.*

test: clean buildForOSX
	go test -covermode count -coverprofile cover.out ./...

testWin: cleanWin buildForWin
	go test -covermode count -coverprofile cover.out ./...

deploy: clean buildForAWS
	sls deploy --verbose --config serverless.yml

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
