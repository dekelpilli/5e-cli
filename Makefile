build:
	go build -o 5e-cli ./cmd/5e-cli

run:
	go run ./cmd/5e-cli $(ARGS)

session:
	open -a "Google Chrome" https://5e.tools/classes.html
	open -a "Google Chrome" https://5e.tools/bestiary.html
	open -a "Google Chrome" https://5e.tools/spells.html

lint:
	golangci-lint run

.PHONY: companion
companion:
	curl -fsSL https://raw.githubusercontent.com/spies-and-spiders/companion/$$(basename $$(curl -fsSLo /dev/null -w '%{url_effective}' https://github.com/spies-and-spiders/companion/releases/latest))/scripts/install.sh | bash
