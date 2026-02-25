start:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

prod:
	docker compose up -build -d

down:
	docker compose down -v

logs:
	docker compose logs -f


