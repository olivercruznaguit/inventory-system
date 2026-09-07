.PHONY: help up backend frontend install build lint

help:
	@echo "Available commands:"
	@echo "  make up       Start Docker, backend, and frontend"
	@echo "  make backend  Start Docker and backend"
	@echo "  make frontend Start frontend"
	@echo "  make install  Install dependencies"
	@echo "  make build    Build backend and frontend"
	@echo "  make lint     Run frontend lint"

up:
	docker compose up -d
	$(MAKE) backend & $(MAKE) frontend

backend:
	go run ./cmd/api

frontend:
	cd frontend && npm run dev

install:
	go mod download
	cd frontend && npm install

build:
	go build ./...
	cd frontend && npm run build

lint:
	cd frontend && npm run lint