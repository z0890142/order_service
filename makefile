generate_doc:
	swag init --parseDependency --parseInternal
build-image:
	docker-compose build  --no-cache
start:
	docker-compose up -d  