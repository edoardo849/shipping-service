stop:
	docker-compose -f ./deployments/docker-compose.yml down -v

build:
	docker-compose -f ./deployments/docker-compose.yml build

run:
	docker-compose -f ./deployments/docker-compose.yml up

integration-tests:
	newman run ./test/integration/req.json -n 40