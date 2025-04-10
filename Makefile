BUILD_DIR = build
BUILD_ENV = GOOS=linux GOARCH=amd64 CGO_ENABLED=0

TF_DEFAULT_ARGS = -no-color
TF_MOD_PATH = deploy/app
TF_VARS ?=

.PHONY: clean clean-state clean-mods clean-all dryrun deploy destroy

clean:
	rm -rf $(BUILD_DIR)

clean-state:
	rm -f $(TF_MOD_PATH)/terraform.tfstate
	rm -f $(TF_MOD_PATH)/terraform.tfstate.backup

clean-mods:
	rm -rf $(TF_MOD_PATH)/.terraform
	rm -f $(TF_MOD_PATH)/.terraform.lock.hcl

clean-all: clean clean-state clean-mods

$(BUILD_DIR)/bootstrap:
	$(BUILD_ENV) go build -ldflags="-s -w" -o $(BUILD_DIR)/bootstrap main.go

$(BUILD_DIR)/lambda.zip: $(BUILD_DIR)/bootstrap
	cd $(BUILD_DIR) && zip -j lambda.zip bootstrap

dryrun: $(BUILD_DIR)/lambda.zip
	terraform -chdir=$(TF_MOD_PATH) init $(TF_DEFAULT_ARGS)
	terraform -chdir=$(TF_MOD_PATH) plan $(TF_DEFAULT_ARGS) -var-file=$(TF_VARS)

deploy: $(BUILD_DIR)/lambda.zip
	terraform -chdir=$(TF_MOD_PATH) init $(TF_DEFAULT_ARGS)
	terraform -chdir=$(TF_MOD_PATH) apply $(TF_DEFAULT_ARGS) -auto-approve -var-file=$(TF_VARS)

destroy:
	terraform -chdir=$(TF_MOD_PATH) destroy $(TF_DEFAULT_ARGS) -var-file=$(TF_VARS)