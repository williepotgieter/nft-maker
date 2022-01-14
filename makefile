dev:
	docker-compose up --remove-orphans

stop:
	docker-compose down

rebuild:
	docker-compose down && docker-compose up --remove-orphans --force-recreate --build

clean:
	docker-compose down && docker volume rm nft-maker_app-db