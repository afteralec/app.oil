# Petrichor App

## Before Developing Locally

As an assumption, all of these instructions assume that you're developing locally
on MacOS or any of the more widely-used Linux distributions.

Before developing locally, install all [required tooling](#required-tooling) and
familiarize yourself with the commands in the Makefile.

## To Start the App Locally

1. Run Docker.
2. Start a local MySQL server.
   a. Use `root` and `pass` as the credentials, and create a database called `test`.
3. In the db package, run `make migrate`.
4. Run `make redis` to start a local Redis server.
5. Run `make dev` to run the app locally.

## Required Tooling

Core Go version:
[Go 1.22](https://go.dev/dl/)

Node, for running Javascript tests locally:
Node 20.x

I recommend using [nvm](https://github.com/nvm-sh/nvm) to manage Node versions.

For hot reloading a development server:
[Air CLI](https://github.com/cosmtrek/air)

Containerization:
[Docker](https://www.docker.com)

## Static Assets

SVG Loaders provided by [SVG-Loaders](https://github.com/SamHerbert/SVG-Loaders)
