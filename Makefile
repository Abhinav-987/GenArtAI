# run: build
# 	@./bin/genartai

install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod tidy
	@go mod download
	@npm install -D tailwindcss
	@npm install -D daisyui@latest

css:
	@tailwindcss -i view/css/app.css -o public/styles.css

build:
	@templ generate view
	@go build -o genartai main.go

up:
	@go run cmd/migrate/main.go up

drop:
	@go run cmd/drop/main.go up

down:
	@go run cmd/migrate/main.go down

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))


