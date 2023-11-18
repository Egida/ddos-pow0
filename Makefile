# Definitions
ROOT                    := $(PWD)
GOLANG_DOCKER_IMAGE     := golang:1.21
BINARY_NAME             := ddos
COMPOSE_FILE_NAME       := ./compose.yml

.PHONY: migration create up down

#MAKEFLAGS += --silent

build:
	 GOARCH=arm64 GOOS=linux go build -o bin/${BINARY_NAME} -v main.go

composeUp:
	docker-compose -f $(COMPOSE_FILE_NAME) up -d

composeDown:
	docker-compose down

clean:
	rm ${BINARY_NAME}

run:
	./${BINARY_NAME}
