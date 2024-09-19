protoc:
	powershell ./scripts/protoc.ps1

docker-up:
	docker-compose down && docker-compose up --build -d

docker-down:
	docker-compose down