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
	@if [ "$(storage)" = "postgres" ]; then \
		STORAGE=$(storage) docker-compose -f docker-compose.postgres.yml build --no-cache; \
	else \
		STORAGE=$(storage) docker-compose -f docker-compose.cache.yml build --no-cache; \
	fi

test:
	-go test -v ./tests/...
