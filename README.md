# API For Product Services

## Motivation
My motivation is to create one portofolio project but not complete portofolio/application, just small pieces like just user management, product api, cart, or other small pieces. By doing this project in small pieces, i can move from one to another way, another design system, another database, another architecture.

## Description
Design of this product API is:
- Product
- Variant
- Category
Relation:
- One product have many variant, one variant just have one product
- One product can have many category, on category can have many product

The relation is one to many and many to many

## How To Install This Project
You need some cli tool for golang development:
- [https://github.com/cosmtrek/air](https://github.com/cosmtrek/air) - Hot reloading in golang
- [https://github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) - Migration database
- [https://github.com/go-critic/go-critic](https://github.com/go-critic/go-critic)
- [https://golangci-lint.run/](https://golangci-lint.run/) - Linter
- [https://github.com/securego/gosec](https://github.com/securego/gosec) - Golang security checker

Don' worry, use can install all by type this command:
`make setup`

it will install all to folder `tmp/`
Or, if you want install one by one its ok, but change the `Makefile`.

after that, type `go mod install`
### Migration Database
Makesure you have edit the `.env` for database setup
type `make cmgr name=FILL_THIS_WITH_MIGRATION_NAME`, for making the migration
type `make migup`, for migration up
type `make migdown`, for migration down

## How To Run This Project
type `make run` to run it using air for hot reloading
type `make docker.dev` to run it using docker
type `make docker.stop` to clear the docker that build from previous command

## How To Use This Project
The function is already done but use with your enhancement.

## Thank you
