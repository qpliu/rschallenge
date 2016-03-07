SRC=compute.go notes.go score.go scrape.go

notes: $(SRC)
	go build -o $@ $(SRC)

clean:
	go clean

test:
	go fmt && go vet && go test
