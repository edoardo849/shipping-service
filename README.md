# Shipping Service

This is a demo on how to implement multiple services in Go using Docker and Docker Compose. The system creates 2 Api services using one Dockerfile configuration for both:

- **API Server**: the service that receives order information and stores them in a MySQL DB. It then sends the order information to a third-party shipping service API
- **Delivery Server**: a mock of the third-party shipping service

## Installation


- Requirements:
    - [Docker](https://www.docker.com/)
    - [Docker Compose](https://docs.docker.com/compose/)
    - [Make](https://www.tutorialspoint.com/unix_commands/make.htm)

```bash

# from the root directory of the project, run:
make build && \
    make run

```

The API Server will run on http://localhost:8080/v1/orders with Basic Auth:
- username: username
- password: password

To run the integration tests install [NewMan](https://www.npmjs.com/package/newman):

```
npm install -g newman

make integration-tests

```