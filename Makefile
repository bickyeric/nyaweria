new-migrate:
	podman compose -f .devcontainer/docker-compose.yml -p nyaweria_devcontainer run --rm db-migrate create -dir /db -ext sql new

migrate:
	podman compose -f .devcontainer/docker-compose.yml -p nyaweria_devcontainer run db-migrate

exec-db:
	podman compose -f .devcontainer/docker-compose.yml -p nyaweria_devcontainer exec -it db psql -U nyaweria_rw nyaweria_dev
