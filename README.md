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