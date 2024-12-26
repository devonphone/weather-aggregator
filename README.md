# Weather Aggregator Application

## Project Overview
The Weather Aggregator application is a Go-based microservice designed to fetch and cache weather data from multiple weather APIs, including OpenWeather and WeatherAPI. The service uses Redis for caching and includes rate-limiting capabilities.

---

## Features
1. Fetches weather data from multiple providers (OpenWeather, WeatherAPI).
2. Implements caching with Redis to reduce redundant API calls.
3. Supports rate-limiting to prevent abuse.
4. Configurable environment variables for flexible deployment.
5. Dockerized setup for streamlined development and deployment.

---

## Prerequisites
- Docker and Docker Compose
- Go (if running locally)
- Redis (cloud-based or local instance)

---

## Environment Variables
The application requires the following environment variables, provided in a `.env` file:

```env
# Server Configuration
PORT=8080
CACHE_DURATION=30m
RATE_LIMIT_REQUESTS=60
RATE_LIMIT_DURATION=1m

# Redis Configuration
REDIS_ADDR=redis-19314.c252.ap-southeast-1-1.ec2.redns.redis-cloud.com:19314
REDIS_USERNAME=default
REDIS_PASSWORD=jX3UeXtfXq78ST9w4fO9xSQE85WwneCW

# Weather API Keys
OPENWEATHER_API_KEY=2677f7c214be2d4c307f44a4c6422c1e
WEATHERAPI_KEY=56677165c854477d96480002242512
```

---

## Docker Setup

### 1. Build and Run the Application
To build and run the application using Docker:

```bash
docker-compose up --build
```

This will:
- Build the Go application container.
- Start the Redis container and network.
- Expose the application on `localhost:8080`.

---

## Project Structure

```
weather-aggregator/
├── Dockerfile                # Dockerfile for building the Go application
├── docker-compose.yml        # Docker Compose setup
├── go.mod                    # Go module dependencies
├── go.sum                    # Go module checksums
├── .env                      # Environment variables file
├── main.go                   # Application entry point
├── handlers/                 # API handlers
├── providers/                # Weather provider integrations
├── utils/                    # Utility functions (e.g., caching, rate limiting)
├── tests/                    # Unit and integration tests
└── README.md                 # Project documentation
```

---

## How to Use

### Endpoints

1. **Fetch Weather Data**
   - **GET** `/weather?city=<city>`
   - Fetches weather data for the specified location.

2. **Health Check**
   - **GET** `/stats`
   - Checks the status of the application.

### Example Request
```bash
curl http://localhost:8080/weather?city=Jakarta
```

---

## Testing

Run the unit tests using the following command:

```bash
go test tests/unit_test.go
```

Ensure Redis is running before executing the tests. You can use the provided Docker Compose setup to start a Redis instance.

### Test Descriptions

TestOpenWeatherProvider:
Verifies the ability to fetch weather data from the OpenWeather API.
Asserts that the data matches the queried city and contains valid temperature values.

TestWeatherAPIProvider:
Checks the functionality of the WeatherAPI integration.
Ensures the data returned matches the queried city and contains non-zero temperature values.

TestFetchFromProviders:
Simulates concurrent requests to multiple weather providers.
Validates that at least one provider successfully returns valid weather data for the queried city within a timeout period.

---

## Dependencies

The following Go packages are used in the application:

1. `github.com/joho/godotenv`: For loading environment variables from `.env` file.
2. `github.com/go-redis/redis/v9`: For interacting with Redis.
3. `net/http`: For handling HTTP requests and responses.

Install these dependencies using:

```bash
go mod tidy
```

---

## Troubleshooting

1. **Redis Connection Issues**:
   Ensure that the Redis address, username, and password in the `.env` file are correct. Use cloud-based Redis credentials if using a remote instance.

2. **Port Conflicts**:
   Make sure the `PORT` in the `.env` file is not already in use by another service.

3. **Rate-Limiting Errors**:
   If you encounter rate-limit responses, adjust `RATE_LIMIT_REQUESTS` and `RATE_LIMIT_DURATION` in the `.env` file.

---

## Future Enhancements
1. Implement a frontend to visualize weather data.
2. Improve error handling and logging.

---

## Author
**Devon Gasselyno**

If you encounter issues or have suggestions for improvement, feel free to open an issue or contribute to the project.

---

## License
This project is licensed under the MIT License. See the LICENSE file for details.