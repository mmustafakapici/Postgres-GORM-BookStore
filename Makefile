# Makefile

.PHONY: all build up down  

all: build up


build:
	docker-compose build
	
up:
	docker-compose up -d
	
down:
	docker-compose down
	
