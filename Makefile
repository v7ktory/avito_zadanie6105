include .env
export

.PHONY: dcup dcdown mup mdown 

# docker-compose
dcup:
	docker-compose up --build -d  && docker-compose logs -f

dcdown:
	docker-compose down --remove-orphans

### migrations
mup: 
	migrate -path migrations -database '$(POSTGRES_LOCALHOST)?sslmode=disable' up
mdown:
	echo "y" | migrate -path migrations -database '$(POSTGRES_LOCALHOST)?sslmode=disable' down
