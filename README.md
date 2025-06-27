# Fizz &  Buzz

Test for **LEBONCOIN**

# Golang Test

## Introduction

> - This test is for LEBONCOIN interview process. All the information regarding the expectations will be kept secret and never shared.
> - All the below test is meant to be fully realized by Maxime GUITTET and only. 

The original fizz-buzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by "fizz", all multiples of 5 by "buzz", and all multiples of 15 by "fizzbuzz". The output would look like this: "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,...".

## Goal

**Write a simple fizz-buzz REST server.**

Your goal is to implement a web server that will expose a REST API endpoint that:

- Accepts five parameters : three integers int1, int2 and limit, and two strings str1 and str2.
- Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.

The server needs to be:

- Ready for production
- Easy to maintain by other developers

Add a statistics endpoint allowing users to know what the most frequent request has been. This endpoint should:

- Accept no parameter
- Return the parameters corresponding to the most used request, as well as the number of hits for this request

# Fizz'n'Buzz

This project was developed as a test for LEBONCOIN. It was developed by Maxime GUITET (Emixam23).

- Introduction
- Architecture
    - Data Models
      - Database (sql row data)
      - API (json response data)
    - API schemas
- Start the project (two main solutions)
    - Env settings
    - (A) PostgreSQL Docker Start
    - (A) Make Run on Go preinstalled environment
    - (B) Run both the API and the database with docker-compose
- Endpoints at your disposal
- Possible evolutions
- Conclusion

## Introduction

Based on the first interview with LEBONCOIN, I got a test related to the development of an API. The concept is pretty simple as I have to recreate the known algorithm "Fizz And Buzz" within an API.
Also, an additional feature was required regarding the statistics to allow users to know what the most frequent "played" request has been.

As I wanted to add some additional endpoints, the health is available at `/health` but also the history can be fetched at `/history`.
Another endpoint as been added regarding the statistics, it's `/stats`. It returns all the stats grouped by count (hits) and can be sorted as required.

Of course, `/stats/most_used` is the endpoint that fits the requirement of this test which is to allow users to know what's the most frequent "played" request.

I hope this API will please you!

> Note:
> - This project has been realised using Golang 1.17
> - A Makefile is at your disposal to use this API, you have several commands available such as test, lint, build, run, start, vendor etc.
> - A "postgresql database launcher" command/bash file is at your disposal under `/cmds` (The database is run throughout Docker).

## Architecture

Hexagonal Architecture has been applied for the development. PostgreSQL is the database used as it is well known at LEBONCOIN.

### Data Models

In this section you can find/view the different data models used across the layers of this API. 

**Database (sql row data)**

| Property    | Type       | PostgreSQL Column | PostgreSQL type | PostgreSQL constraint |
|-------------|------------|-------------------|-----------------|-----------------------|
| tableName   | struct{}   | fnb_requests      |                 |                       |
| ID          | uint64     | id                | SERIAL          | PRIMARY KEY           |
| RequestDate | *time.Time | request_date      | TIMESTAMP       | NOT NULL              |
| N1          | uint32     | n1                | BIGSERIAL       | NOT NULL              |
| S1          | string     | s1                | TEXT            | NOT NULL              |
| N2          | uint32     | n2                | BIGSERIAL       | NOT NULL              |
| S2          | string     | s2                | TEXT            | NOT NULL              |
| Limit       | uint64     | rlimit            | BIGSERIAL       | NOT NULL              |

**API (json response data)**

*fnbRequest data model example*
```json
{
  "id": 2,
  "request_date": "2021-10-31T14:57:56.307632Z",
  "n1": 3,
  "s1": "fizz",
  "n2": 5,
  "s2": "buzz",
  "limit": 100
}
```

*fnbRequestInputStats data model example*
```json
{
  "n1": 3,
  "s1": "fizz",
  "n2": 5,
  "s2": "buzz",
  "limit": 100,
  "count": 2
}
```

### API schemas

In this section, we will have a look at the defined architecture of this API.

The schemas were made using [Draw IO](draw.io). The below PNGs can be imported into Draw IO for any update(s).

This schema mostly represent the project business/technical representation.

> Note: Each of the Infra/Domains/UI use their own data models. Between each, mappings takes place and the Domain is the source of trust.

<div>
  <img src="/docs/fizz-and-buzz general architecture.drawio.png"/>
</div>

This second schema is representing a possible approach which can be given to a developer that would need a deeper (technical) approach.

> Note: Given this type of schema, it helps the developer to get into the subject with more information. 
> As you can see, the developer knows that configurations will be given to him for some parts, but also he can directly know how to implement the big parts of the project.
> However, even if postgresql data schema is provided, we don't provide any information for the expected uri results, but just a description. This let some freedom to the developer (in this scenario)

<div>
  <img src="/docs/fizz-and-buzz implementation and definitions example.drawio.png"/>
</div>

## Start the project

A Makefile is at your disposal, as weel as a docker-compose file and scripts. From these you can start by multiple way the given api:
- Running the database locally (with the scripts?) and running the project using your own go environment.
- Run the API from the Makefile and connect to a remote database instance?
- Use docker-compose to run the API and its database with a Dockerized approach.

List of the argument for make:

| Commands & arguments   | Description                                                               | 
|------------------------|---------------------------------------------------------------------------|
| make                   | It's the default command: it runs build, test, and lint below commands    | 
| make test              | Runs both the unit-tests & integration-tests below commands               |
| make lint              | Verify the written Golang code and return possible language missuses      | 
| make unit-tests        | Runs the unit tests and generate a file containing the code test coverage |
| make integration-tests | Runs the integration tests                                                |
| make build             | Builds the app and creates the executable under `.build`                  |
| make vendor            | Vendors the dependencies                                                  |
| make run               | Starts the project using Docker                                           |
| make start             | Starts both the PostgreSQL database and the project using docker-compose  |
| make mock              | Generates the mocks related to the used/consumed layers of the project    |

### Env settings

To make this project working, you will need to configure your environment. The easiest way is to create a `.env` file. Here is a template for you:
```
LOGGER_LEVEL=info
LOGGER_AS_JSON_FORMAT=true
DATABASE_HOST=localhost
DATABASE_NAME=testdb
DATABASE_USER=testapi
DATABASE_PASSWORD=fizznbuzz
DATABASE_PORT=5432
DATABASE_RETRY_AMOUNT_ON_FAIL=3
FNB_SERVICE_ZERO=1
ROUTER_HOST=0.0.0.0
ROUTER_PORT=8080
```

> Note: Some of the values, if not set, are defaulted to predefined value(s) by the API itself.

> Note 2: A `.env.example` file is at your disposal as file example (environment variables definition/setting file).

### PostgreSQL Docker Start

**windows**
```
.\cmds\postgres-start.cmd
```

**linux**
```
./cmds/postgres-start.sh
```
> chmod might be required: chmod +x cmds/postgres-start.sh

Once the database is running and up, starts the project like you want to.

### Make Run on Go preinstalled environment

```
make run
```

### Run both the API and the database with docker-compose

```
make start
```
or
```
docker-compose up
```

## Endpoints at your disposal

| Method | Endpoint         | Uri Expected Param                                              |  Results                                                                                 |
|--------|----------------  |-----------------------------------------------------------------|------------------------------------------------------------------------------------------|
| GET    | /                |                                                                 | Returns the status (health) of the service                                               |
| GET    | /health          |                                                                 | Returns the status (health) of the service                                               |
| GET    | /fizz-and-buzz   | [n1: uint32, s1: string, n2: uint32, s2: string, limit: uint64] | Validate the entry, generate the results and store the requested input into the database |
| GET    | /history         | [limit: uint64]                                                 | Returns the history of the {n} last requested fizz and buzz entries                      |
| GET    | /stats           | [sorted: boolean]                                               | Returns the stats ({sorted} if required as per parameters)                               |
| GET    | /stats/most_used |                                                                 | Returns the most used fizz and buss params                                               |

## Possible evolutions

As this project is very, very small, talking about microservices would result into loosing time. Staying on a monolithic architecture would be wise.
Regarding potential evolution, this hexagonal architecture gives us freedom on many points:
- We can switch our gin rest api to another package, mux gorilla? (easy with hexagonal architecture to switch such things)
- We can switch to gRPC
- We can switch to another database, maybe NoSQL?
- We could also (if dataset gets bigger and bigger) set a cache for the statistics (`/stats`) using redis or a local cache?

Thanks to the chosen architecture, many and many options are given to us. Of course, some evolution might require some side evolution regarding unit/integration testing.

## Conclusion

This project was realised over a couple of hours (app code, documentation in/out code, tests, dockerization).
I used Hexagonal Architecture and implemented everything in Go/PostgresSQL.

The following libraries were used to realise this project:
```
bou.ke/monkey v1.0.2
github.com/DATA-DOG/go-sqlmock v1.5.0
github.com/cucumber/godog v0.12.2
github.com/gin-contrib/cors v1.3.1
github.com/gin-gonic/gin v1.7.4
github.com/golang/mock v1.6.0
github.com/lib/pq v1.10.3
github.com/rdumont/assistdog v0.0.0-20201106100018-168b06230d14
github.com/rs/zerolog v1.25.0
github.com/spf13/viper v1.9.0
github.com/stretchr/testify v1.7.0
```

I hope this project will confirm your thoughts about me.

Max
