# Stock Simulator Backend

This is the backend for a stock simulator app, which was originally meant to be released during BITS Pilani's Apogee '24.

It is a **full-fledged real-time stock simulator** that uses RabbitMQ for in-memory trade transactions and pgsql for persistent storage. It is written in Go using the Gin web framework.

The existing engine processes limit buy and limit sell orders. It also has a simple frontend for adding and viewing stocks and news articles (which would eventually be used to simulate stock prices).

I initially used Apache Kafka for the messaging system, but due to resource constraints on the production server (about 1GB of RAM), I had to switch to RabbitMQ.

The authentication is a bit hacky because it was meant to be for a 2-day app, but it works nevertheless. The standard procedure would be to rotate access tokens and refresh tokens. It uses Google OAuth2 for the app-side auth and session cookies for the admin panel.

**Most core components are in the `engine` package.**

# Stack Used

- Go (Gin)
- RabbitMQ
- Postgresql
- SQLC
- Frontend (Admin Panel):
  - Standard Go Templating
  - TailwindCSS
  - HTMX
- Development Tools:
  - Air
  - Make
- Docker

# Setup for Development

- Make sure you have Go, Node, Docker, Air and SQLC installed.
- Clone the repository.
- Copy the `.env.example` file to `.env` and fill in the required details.
- Run `make devdb` to create the development database.
- Run `make migrateup` to run the migrations. (You might need to change the database url in the `migrate` command in the `Makefile`).
- Run `npm i` to install the frontend dependencies.
- Run `npm run dev` to start the tailwind listener.
- Run `air` to start the development server.

# DB Diagram

<div align="center">
<img src="./db_design.png" width="600">
</div>
