# iBank-service
An example of a microservices app built with Golang. It is a dead simple internet banking service. It contains
two different services - **account-service** and **transaction-service**.

- **account-service** - It handles creation of accounts and assigning account number. It also has a gRPC service that runs on port **9003**, this gRPC service is used for inter-service communication with **transaction-service** mainly for verifying JWT tokens and retrieving account details.
- **transaction-service** - Handles funds transfer and also keeps history of transactions/fund transfer. It communicates with **account-service** gRPC server on port **9003**

**Nginx** is used as a reverse-proxy/load balancer so as to have a single entry point to all the services


**docker-compose** is used to run all services/containers. The following is the list of all the containers - 

* Nginx - for reverse proxying/load balancing requests
* mysql database - database uses to store account balance and transaction history
* mongoDB - used to store account data
* account-service - main account service
* transaction-service - main transaction service

## Running the app

* clone the repo `git clone github.com/adigunhammedolalekan/ibank-service`
* `> cd ibank-service`
* `> make run`
* `> docker-compose up`

#### Available endpoints

* `localhost:80/api/account` - to create an account `BODY - POST {email, name, password}`
* `localhost:80/api/account/authenticate` - to authenticate an account for JWT token `BODY - POST {account_number, password}`
* `localhost:80/api/txn/fund` - add dummy fund to an account `BODY - POST {amount}, HEADER - Bearer {JWT_TOKEN}`
* `localhost:80/api/txn/transfer` - send money from authenticated account to another account using their account numbers `BODY - POST {to_account, amount}, HEADER - Bearer {JWT_TOKEN}`
* `localhost:80/api/txn/me/history` - fetch authenticated account transaction history `GET, HEADER - {JWT_TOKEN}`

