run:
	docker build -t breaker-case:test .
	docker-compose up --remove
