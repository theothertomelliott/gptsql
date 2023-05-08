# GPT SQL

An experiment to generate working SQL queries for real databases from natural language questions.

It uses GPT-3 (via the OpenAI API) to generate SQL queries based on a user question, and then runs them against the database to get the real results.

Currently, only Postgres databases are supported.

## Quickstart

Set your OpenAI API token in the `OPENAI_API_TOKEN` environment variable. If you have an OpenAI account, you can obtain a token at https://platform.openai.com/account/api-keys.

Set the `POSTGRES_CONN_STRING` environment variable to the URL of your database. For example, if you're using Postgres, you might set it to `postgres://user:password@localhost:5432/dbname`.

You can now run `gptsql`:

```bash
go run .
```

## Using an example database

A Docker Compose file is included to run a Postgres database with some example data.

To provide example data, create a file `data/init.sql` with the SQL commands to create the tables and insert the data.

Then start the database with the command:

```bash
docker compose up
```

You can now use the `example.sh` script to run queries against your example database.