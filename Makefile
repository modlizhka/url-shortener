all:run

run:
	docker compose up

clean:
	-rm -r logs

rebuild:
	docker compose build --no-cache

tests:
	go test -v ./internal/service
