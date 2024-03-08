# Parking API

## Overview
This project is a simple API for Parking APP. Here is endpoints that already exists on this Project.
- Register
- Login
- Create User
- Create a Parking Lot
- Booking Parking Slot
- Finish Parking Slot
- View Parking Slot
- Booking Summary

This project using `clean code architecture`, there are 3 layers in this project (`handler`, `service`, and `repository`)

## Library
In this project, we use these libraries:
- Uber Dig (Dependency Injection)
- Squirrel
- Golang Validator
- Gomock
- Postgresql
- Redis

> We use quirrel rather than ORM because ORM is very weight. So we decide to use squirrel for performance purpose.

## Installation & Run

1. To install this project, you can just using this command

```bash
go mod tidy
```

2. Don't forget to change `.env.example` to `.env`

4. To run `redis` & `postgresql` you can just using this command
```bash
docker compose up -d
```

5. To run `migration` you can just using this command
```bash
./typicalw pg migrate
```

6. To run `DB seeding` you can just using this command
```bash
./typicalw pg seed
```

7. To `run this project`, you can just using this command
```bash
./typicalw r
```

## Swagger

You can check out the swagger documentation in this url
```bash
http://localhost:8089/swagger/index.html
```

## Debug Performance
To check `pprof` you can open this URL
```bash
http://localhost:8089/debug/pprof/
```

## Unit Test
For Unit Test, you can check out to the test file
```bash
{projectFolder}/internal/app/service/parking/parking_test.go
```

To running unit test you can just using this command
```bash
./typicalw test
```
