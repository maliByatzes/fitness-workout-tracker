# Fitness Workout Tracker

A fitness workout tracker backend application. Where users can
create their own workout exercises and much more.

## Features

- Users can perform basic authentication including registration,
login, logout.
- Users can create as many workouts and add many exercises from
the list of available exercises.
- Users can set up times of when to do the workouts.
- Users can update when they are finished with a certain workout.

## Tech Stack

- **Go (Golang)**: The Go Programming Language.
- **PostgresSQL**: SQL database.
- **Docker**: Open-source software for deploying and running of containerized applications.

## Pre-requisites

- [Go](https://go.dev/) at least version 1.23
- [Docker](https://docker.com)
- [Golang Migrate](github.com/golang-migrate/migrate) must be installed globally on your system

## Installation

1. Clone the repository
   ```sh
   git clone https://github.com/maliByaztes/fitness-workout-tracker
   cd fitness-workout-tracker
   ```

2. Create `.postgres.env` and `.env` in the root directory. Check `postgres.example.env`
and `example.env` files for examples and required variables.

**NOTE**: Remember to change the `DB_URL` in Makefile according to your postgres.env file.

3. Create postgres container
   ```sh
   make up
   ```

4. Run migrations against database
   ```sh
   make run-migrations
   ```

5. Run server
   ```sh
   make server
   ```

The server should be running on `http://localhost:8000`

## API Endpoints

- Gin will log all the available routes when running the server.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

maliByatzes - malib2027@gmail.com

Project Link - https://github.com/maliByatzes/hotel-booking-system
