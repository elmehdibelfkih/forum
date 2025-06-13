up:
	mkdir -p "./requirements/project/logs"
		: > "./requirements/project/logs/internal_errors.log" 
	docker compose -f ./docker-compose.yml up -d

build:
	docker compose -f ./docker-compose.yml build

down:
	docker compose -f ./docker-compose.yml down

status:
	docker compose -f ./docker-compose.yml ps

clean: down
	docker image prune -af
	docker volume prune -f
	docker network prune -f

fclean: clean

re: clean up

