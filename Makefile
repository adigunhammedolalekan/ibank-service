proto:
	for d in src; do \
			for f in $$d/**/proto/*.proto; do \
				protoc --go_out=plugins=grpc:. $$f; \
				echo compiled: $$f; \
			done \
	done
build:
	./build.sh

run:
	make proto
	make build