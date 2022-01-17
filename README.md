# Backend Engineer Coding Challenge

## Description
Simple scanner tries to find secret key leaking in public repository

**How to start:**
- Require docker and docker compose installed
- Start with
```sh
docker-compose up
```
- App start in port 6080, API documentation can be found [here](docs/openapi-spec/swagger.json)

**Some more details design:**
- The queueJob separates the scanning job to a pool of workers. The number of workers can be configured through nWorker.
- Important parts are abstracted so can it be easily swap out with service  of your choice.
