protoc:
	powershell ./scripts/protoc.ps1

docker:
	docker-compose down && docker-compose up --build -d