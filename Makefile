NAME = nomad-deploy

BUILD_DIR = build
BUILD_ENV = GOOS=linux GOARCH=amd64 CGO_ENABLED=0

.PHONY: clean dryrun 

clean:
	rm -f $(NAME)
	rm -rf $(BUILD_DIR)

$(BUILD_DIR)/bootstrap:
	$(BUILD_ENV) go build -ldflags="-s -w" -o $(BUILD_DIR)/bootstrap main.go

$(BUILD_DIR)/lambda.zip: $(BUILD_DIR)/bootstrap
	cd $(BUILD_DIR) && zip -j lambda.zip bootstrap

dryrun: $(BUILD_DIR)/lambda.zip
	terraform -chdir=deploy/app init -no-color
	terraform -chdir=deploy/app plan -no-color -var-file=../nonprod.tfvars

all: deploy
