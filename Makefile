all:run

run:
	@if [ "$(storage)" = "postgres" ]; then \
		STORAGE=$(storage) docker-compose -f docker-compose.postgres.yml up; \
	else \
		STORAGE=$(storage) docker-compose -f docker-compose.cache.yml up; \
	fi

clean:
	-rm -r logs

rebuild:
	docker compose build --no-cache

tests:
	go test -v ./tests/...
