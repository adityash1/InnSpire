# InnSpire API

This is a backend API for a hotel reservation system, written in Go. It provides endpoints for managing users, hotels, rooms, and bookings.

## Getting Started

To run the application, you need to have Go installed on your machine. You can download it from the [official Go website](https://golang.org/dl/).

Once you have Go installed, clone the repository and navigate to the project directory.

You can then run the application using the following command:
```
make run
```
The application will start and listen on the port specified in the HTTP_LISTEN_ADDRESS environment variable.

## Docker

A Dockerfile is provided for containerizing the application. To run the application in a Docker container, use the following command:
```
make docker

OR

docker build -t hotel-reservation-api .
docker run -p 3000:3000 hotel-reservation-api
```
Replace 3000 with the port you want to bind to on your host machine.

### Environment Variables

The application requires several environment variables to run correctly. These are loaded from a `.env` file at startup. The required environment variables include:

- `HTTP_LISTEN_ADDRESS`: The address the server should listen on.
- `MONGO_DB_URL`: The URL of your MongoDB instance.
- `MONGO_DB_NAME`: The name of your MongoDB database.
- `JWT_SECRET`: The secret key used for signing JWT tokens.

### API Endpoints

The API provides several endpoints for managing users, hotels, rooms, and bookings. All endpoints are prefixed with `/api/v1`.

- User Endpoints:
    - `PUT /user/:id`: Update a user.
    - `DELETE /user/:id`: Delete a user.
    - `POST /user`: Create a new user.
    - `GET /user`: Get all users.
    - `GET /user/:id`: Get a specific user.

- Hotel Endpoints:
    - `GET /hotel`: Get all hotels.
    - `GET /hotel/:id`: Get a specific hotel.
    - `GET /hotel/:id/rooms`: Get all rooms in a specific hotel.

- Room Endpoints:
    - `POST /room`: Create a new room.
    - `POST /room/:id/book`: Book a room.

- Booking Endpoints:
    - `GET /booking/:id`: Get a specific booking.
    - `GET /booking/:id/cancel`: Cancel a booking.

- Admin Endpoints:
    - `GET /admin/booking`: Get all bookings.

  
### Database

The application uses MongoDB for data storage. The `db` package contains the interfaces and implementations for interacting with the database.

#### Seeding the Database

The application includes a script for seeding the database with initial data. This script is located in the `scripts/seed.go` file.
The script creates a number of users, hotels, rooms, and bookings. It uses the `fixtures` package, which contains helper functions for creating these entities.

To run the seeding script, use the following command:
```
make seed
```
Please note that the seeding script uses the same environment variables as the main application to connect to the MongoDB database. Make sure these environment variables are correctly set in your `.env` file before running the script.
The script also prints out JWT tokens for the created users. You can use these tokens to authenticate API requests during testing or development.

&nbsp;

**Authentication** - The API uses JWT for authentication. The `api/jwt.go` file contains the middleware for JWT authentication. The `api/auth_handler.go` file contains the handler for the `/auth` endpoint, which is used to authenticate users and issue JWT tokens.

**Models** -  The `types` package contains the data models used in the application. As well as any related types.

**Error Handling** - The application uses a custom error handler for handling HTTP errors. The error handler is defined in the `api` package.
