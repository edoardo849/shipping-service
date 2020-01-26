stop:
	docker-compose -f ./deployments/docker-compose.yml down -v

build:
	docker-compose -f ./deployments/docker-compose.yml build

run:
	docker-compose -f ./deployments/docker-compose.yml up