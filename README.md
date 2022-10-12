## GRPC Contact Manager

This project is a simple project to that creates contact records for registered users.

It is a small project to understand some basics about integrating golang applications to different services and also some technologies.
Such service includes:
* gRPC
* postgres
* cockroachdb
* prometheus for exporting metrics
* grafana for visualization

It is a simple project that does the following:
* creates a user account via gRPC and REST
* creates a contact through a REST

# Setup

* Install [Docker](https://docker.com) and docker-compose
* Clone the Repo
* RUN `cp .envs/app.env.example .envs/.env` and set the environment variables in `.envs/.env`
* To build the docker image, run `make bd`

# Startup

There are different ways to start the applciation.
1. run `make run`. This will read the environment file from `.envs/.env` and use it for the application.
2. run `docker-compose up` This should be after the docker image is built.

# Application Ports:

* Http: 3500
* gRPC: 3501
* grafana: 3000