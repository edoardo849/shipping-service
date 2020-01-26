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

```

curl --location --request POST 'http://localhost:8080/v1/orders' \
--header 'Content-Type: application/javascript' \
--header 'Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
--data-raw '[
    {
      "id": 101,
      "email": "jon@doe.ca",
      "total_price": "254.98",
      "total_weight_grams": 100,
      "order_number": 1234,
      "shipping_lines": [
        {
          "id": 271878346596884015,
          "title": "Generic Shipping",
          "price": "10.00"
        }
      ],
      "shipping_address": {
        "first_name": "Steve",
        "address1": "123 Shipping Street",
        "city": "Shippington",
        "postcode": "se26hg"
      }
    },
    {
      "id": 102,
      "email": "jane@doe.uk",
      "total_price": "30.00",
      "total_weight_grams": 789,
      "order_number": 5678,
      "shipping_lines": [
        {
          "id": 271878346596884016,
          "title": "Next Day",
          "price": "25.00"
        }
      ],
      "shipping_address": {
        "first_name": "Bob",
        "address1": "89 Shipping Lane",
        "city": "Shipville",
        "postcode": "cb227hd"
      }
    }
  ]'

```

To run the integration tests install [NewMan](https://www.npmjs.com/package/newman):

```
npm install -g newman

cd ./test/integration
newman run

```