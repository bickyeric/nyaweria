new-migrate:
	migrate create -dir ./db -ext sql new

migrate:
	migrate -source file:./db -database postgresql://nyaweria_rw:supersecret123@db:5432/nyaweria_dev?sslmode=disable up

migrate-down:
	migrate -source file:./db -database postgresql://nyaweria_rw:supersecret123@db:5432/nyaweria_dev?sslmode=disable down 1

migrate-version:
	migrate -source file:./db -database postgresql://nyaweria_rw:supersecret123@db:5432/nyaweria_dev?sslmode=disable version

lint:
	golangci-lint run

exec-db:
	podman compose -f .devcontainer/docker-compose.yml -p nyaweria_devcontainer exec -it db psql -U nyaweria_rw nyaweria_dev
