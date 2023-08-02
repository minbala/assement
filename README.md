## Assessment Test

## API Documentation (Swagger)

To genereate swagger documentation

```sh
$ cd back_end
$ swag fmt
$ swag init --parseDependency --parseInternal
```

Read documentation at [doc](http://localhost:8080/docs/index.html):

To do database migration

```sh
$ cd back_end/persistence/mysql/migration
$ go run main.go
```
To run the backend application

```sh
$ cd back_end
$ go run test_assessment
```