# Petrichor App

## To Start the App Locally

1. Run Docker.
2. Start a local MySQL server.
   a. Use `root` and `pass` as the credentials, and create a database called `test`.
3. In the db package, run `make migrate`.
4. Run `make redis` to start a local Redis server.
5. Run `make test` to run the full test suite.
6. Run `make dev` to run the app locally.

## Static Assets

SVG Loaders provided by [SVG-Loaders](https://github.com/SamHerbert/SVG-Loaders)
