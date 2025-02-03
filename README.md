Language: [EN](https://github.com/EvansTrein/iqProgers/blob/main/README.md) [RU](https://github.com/EvansTrein/iqProgers/blob/main/readmeRU.md)

# It's a test assignment 
- Make a REST API for financial transactions with three handles
  - Top up the user's balance.
  - Transferring money from one user to another.
  - View 10 recent transactions of a user.
- The user has a balance as well as a list of transactions. 
- Don't forget to use SQL transactions!!!
- Authorization is not necessary!
- Technologies: go, gin, pgx, postgreSQL, docker. 
- Use goose for migrations, do configuration via ENV. 
- Also attach .env.example file to the solution.
- Make the application run through docker-compose using Makefile. When the make run command is called, the containers should be up, migrations should be performed, and the server should start on port 8080.
- For complex queries you can use query-builder. 
- You should also write unit tests for the service layer of the application
- In the application use **clean architecture** approach
(handler -> service -> repository).

You can write with clarifying questions about the task.
After execution, send a link to github with the completed test case.
In case of successful completion we will schedule an interview. The time for completion was given 6 days (I did it in 3).

# How do I start it up?
For Docker:

- Clone or download the repository - type `make run-docker-compose`
- If you don't use make - enter `docker compose --env-file configForDocker.env up --build -d`.

To run locally:
 - Clone or download the repository
 - To run **need DB (PostgreSQL) and migrate**, enter `make migrate`, without make you will need to manually pass the DB path as a flag (see example in .env files).
 - enter `make run`.
 - If you don't use make - enter `go run cmd/main.go -config ./configLocal.env`.

Users were created in the migration. There are 5 of them, id's from 1 to 5. You can test the API.

## About clarifying questions
In the communication, I mentioned that I didn't use `goose` but used `migrate`. I was allowed to use `migrate`.

Also, in the course of the conversation, it became clear that using query-builder is not necessary. Quote: <u>“We welcome writing ‘raw’ queries, so you don't have to use query-builder”</u>.

The conversation, a bit, bogged down, I was clarifying details. Answered all detailed! And the result - Quote: <u>"if there is a desire to make additional conditions, think about idempotency in your API, when working with payment transactions.
Also, what options are there for competitive debiting of balance, if several SQL transactions will work in parallel.
How will you store money, what data type to use"</u>.

First, I misunderstood about **idempotency of API** and **what data type to use** for money. Answers:
 - “It may be that there was a network failure on the client's side and he repeated the request to write off 500 rubles. It should be written off exactly 500, not 1000. You need to use the idempotency key.”
 - “Keep your money in bigint 500 rubles is 50000 in int together with pennies.”

In general, during the conversation, I realized that here they will look with equal importance at the way SQL queries and ACID are written, not only at Golang itself.

## Execution
Users are created through migration.

The balance here is stored in BIGINT. The last two digits are cents.

For example:
- 1000 in the database is stored as BIGINT 1000<a>00</a>, but in the application it is float64 1000.00.
- Conversions of float64 to BIGINT and vice versa take place on the side of the database. It would be more convenient for me to do in the application to send already ready numbers to the database, but I did it on purpose to show SQL skills. 

For idempotency - in the table with transactions there is `idempotency_key`. If a request comes in, but the key is already there, the `http code 200` and the transaction containing this key are returned.