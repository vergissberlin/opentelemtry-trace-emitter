.PHONY: all js-build go-build

all: js-build go-build

js-build:
	@echo "Building JavaScript project..."
	(cd app/js && npm install && npm run build)

go-build:
	@echo "Building Go project..."
	(cd app/go && go build -o bin/app main.go)

# Run both builds in parallel
parallel-build:
	@echo "Building JavaScript and Go projects in parallel..."
	$(MAKE) -j2 all
