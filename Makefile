BUILD_DIR = build
BUILD_ENV = GOOS=linux GOARCH=amd64 CGO_ENABLED=0

TF_DEFAULT_ARGS = -no-color

.PHONY: clean dryrun deploy

clean:
	rm -rf $(BUILD_DIR)

$(BUILD_DIR)/bootstrap:
	$(BUILD_ENV) go build -ldflags="-s -w" -o $(BUILD_DIR)/bootstrap main.go

$(BUILD_DIR)/lambda.zip: $(BUILD_DIR)/bootstrap
	cd $(BUILD_DIR) && zip -j lambda.zip bootstrap

dryrun: $(BUILD_DIR)/lambda.zip
	terraform -chdir=deploy/app init $(TF_DEFAULT_ARGS)
	terraform -chdir=deploy/app plan $(TF_DEFAULT_ARGS)

deploy: $(BUILD_DIR)/lambda.zip
	terraform -chdir=deploy/app init $(TF_DEFAULT_ARGS)
	terraform -chdir=deploy/app apply $(TF_DEFAULT_ARGS) -auto-approve

all: deploy
