frieden = ./frieden.go

all: run


run:
	go run -race ${frieden}


# build for specific OS target
build-%:
	GOOS=$* GOARCH=amd64 go build -o frieden-$* ${frieden}


build:
	go build -o frieden ${frieden}


# clean any generated files
clean:
	rm -rvf frieden frieden-*
