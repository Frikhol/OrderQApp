# OrderQApp

OrderQApp is a queue management application that connects clients who need to stand in line with agents who can stand in line for them.

## Features

- User authentication (login/register)
- Role-based access (client/agent)
- Queue management
- Real-time queue position updates
- Dynamic pricing based on time and conditions

## Prerequisites

- Go 1.21 or later
- PostgreSQL 12 or later
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/orderqapp.git
cd orderqapp
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
```bash
createdb orderqapp
psql orderqapp < migrations/init.sql
```

4. Configure the application:
Create a `conf/app.conf` file with the following content:
```ini
appname = orderqapp
httpport = 8080
runmode = dev

# Database configuration
dbhost = localhost
dbport = 5432
dbuser = postgres
dbpassword = postgres
dbname = orderqapp

# JWT configuration
jwtSecret = your-secret-key
```

## Running the Application

1. Start the application:
```bash
go run main.go
```

2. Access the application at `http://localhost:8080`

## API Endpoints

### Authentication
- POST `/auth/login` - User login
- POST `/auth/register` - User registration

### Orders
- POST `/orders` - Create a new order
- GET `/orders` - Get user's orders
- PUT `/orders/:id/position` - Update queue position
- GET `/orders/available` - Get available orders (agent only)

## Project Structure

```
orderqapp/
├── controllers/     # Request handlers
├── models/         # Data models
├── routers/        # URL routing
├── static/         # Static files (JS, CSS)
├── views/          # HTML templates
├── conf/           # Configuration files
├── migrations/     # Database migrations
├── main.go         # Application entry point
└── README.md       # This file
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 