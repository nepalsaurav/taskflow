APP_NAME := postfix_admin
OUT_DIR  := bin
GO       := go

LDFLAGS := -s -w
GCFLAGS := all=-trimpath=$(PWD)

build:
	@mkdir -p $(OUT_DIR)
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	$(GO) build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-gcflags "$(GCFLAGS)" \
		-o $(OUT_DIR)/$(APP_NAME)

strip:
	strip $(OUT_DIR)/$(APP_NAME)

upx:
	upx --best --lzma $(OUT_DIR)/$(APP_NAME)

clean:
	rm -rf $(OUT_DIR)
