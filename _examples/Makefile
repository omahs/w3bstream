build_log_wasm:
	@cd log && tinygo build -o log.wasm -scheduler=none --no-debug -target=wasi log.go
	@echo built to log.wasm

build_gjson_wasm:
	@cd gjson && tinygo build -o gjson.wasm -scheduler=none --no-debug -target=wasi parse_json.go
	@echo built to gjson.wasm

build_easyjson_wasm:
	@cd easyjson && tinygo build -o easyjson.wasm -scheduler=none --no-debug -target=wasi parse_json.go
	@echo built to easyjson.wasm

build_token_distribute_wasm:
	@cd token_distribute && tinygo build -o token_distribute.wasm -scheduler=none --no-debug -target=wasi token_distribute.go
	@echo built to token_distribute.wasm

build_word_count_wasm:
	@cd word_count && tinygo build -o word_count.wasm -scheduler=none --no-debug -target=wasi word_count.go
	@echo built to word_count.wasm

build_word_count_v2_wasm:
	@cd word_count_v2 && tinygo build -o word_count_v2.wasm -scheduler=none --no-debug -target=wasi word_count_v2.go
	@echo built to word_count_v2.wasm

build_stream_wasm:
	@cd stream && tinygo build -o stream.wasm -scheduler=none --no-debug -target=wasi main.go
	@echo built to stream.wasm

all: build_log_wasm build_gjson_wasm build_easyjson_wasm build_word_count_wasm build_word_count_v2_wasm build_token_distribute_wasm build_stream_wasm
