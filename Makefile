start:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

prod:
	docker compose up --build -d

down:
	docker compose down -v

logs:
	docker compose logs -f

# Kubernetes / Minikube
k8s-build:
	docker build --platform linux/amd64 -t account:latest -f account/app.dockerfile --target prod .
	docker build --platform linux/amd64 -t catalog:latest -f catalog/app.dockerfile --target prod .
	docker build --platform linux/amd64 -t order:latest -f order/app.dockerfile --target prod .
	docker build --platform linux/amd64 -t graphql:latest -f graphql/app.dockerfile --target prod .
	docker build --platform linux/amd64 -t account-db:latest -f account/db.dockerfile .
	docker build --platform linux/amd64 -t order-db:latest -f order/db.dockerfile .

k8s-load: k8s-build
	minikube image load account:latest
	minikube image load catalog:latest
	minikube image load order:latest
	minikube image load graphql:latest
	minikube image load account-db:latest
	minikube image load order-db:latest

k8s-apply:
	kubectl apply -f k8s/ -R

k8s-deploy: k8s-load k8s-apply

k8s-delete:
	kubectl delete -f k8s/ -R --ignore-not-found --wait=false
