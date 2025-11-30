

docker-up:
	docker compose up -d --build


docker-down:
	docker compose down -v

test-e2e:
	bash e2e-test.sh

