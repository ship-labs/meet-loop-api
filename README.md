# MeetLoop API

A robust Go-based REST API for the MeetLoop platform, featuring JWT authentication, PostgreSQL database integration, and seamless Supabase connectivity.

## Features

- **JWT Authentication** - Secure user authentication and authorization
- **PostgreSQL Database** - Reliable data persistence with migrations
- **Supabase Integration** - Modern backend-as-a-service capabilities
- **CORS Support** - Cross-origin resource sharing for web applications
- **Docker Ready** - Containerized deployment with multi-stage builds
- **Cloud Run Deployment** - Automated CI/CD with GitHub Actions
- **Live Reload** - Hot reloading for development with Air
- **Database Migrations** - Version-controlled schema management
- **Code Generation** - Type-safe database queries with sqlc

## Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24+** - [Download Go](https://golang.org/dl/)
- **PostgreSQL 13+** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Docker** (optional) - [Download Docker](https://www.docker.com/get-started)
- **migrate CLI** - Database migration tool
- **sqlc** (optional) - SQL code generator
- **Air** (optional) - Live reload tool

### Install Required Tools

```bash
# Install migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install sqlc (optional)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install Air for live reload (optional)
go install github.com/air-verse/air@latest
```

## Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/ship-labs/meet-loop-api.git
cd meet-loop-api
```

### 2. Install Dependencies

```bash
make install
```

### 3. Environment Configuration

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Update the `.env` file with your configuration:

```env
DB_PASSWORD=your_database_password
DB_URL=postgres://username:password@localhost:5432/meetloop_db?sslmode=disable
SUPABASE_PROJECT_URL=https://your-project.supabase.co
SUPABASE_API_KEY=your_supabase_anon_key
FRONTEND_URL=http://localhost:3000
JWT_SECRET=your_super_secret_jwt_key
Env=development
PORT=8080
```

### 4. Database Setup

#### Create Database

```bash
# Connect to PostgreSQL and create database
psql -U postgres -c "CREATE DATABASE meetloop_db;"
```

#### Run Migrations

```bash
# Run database migrations
make migrate url="postgres://username:password@localhost:5432/meetloop_db?sslmode=disable"
```

#### Seed Database (Optional)

```bash
# Seed the database with initial data
make seed_db url="postgres://username:password@localhost:5432/meetloop_db?sslmode=disable"
```

## Running the Application

### Development Mode

```bash
# Standard development run
make start

# Or with live reload (requires Air)
make live_reload
```

### Production Mode

```bash
# Build and run
make build
./bin
```

### Using Docker

```bash
# Build the Docker image
docker build -t meet-loop-api .

# Run the container
docker run -p 8080:8080 --env-file .env meet-loop-api
```

## Project Structure

```
meet-loop-api/
├── cmd/                    # Application entrypoints
│   └── main.go            # Main application
├── config/                # Configuration management
├── database/              # Database connection and utilities
├── internal/              # Private application code
│   └── pkg/
│       └── postgres/
│           ├── migration/ # Database migrations
│           └── seed.sql   # Database seed data
├── middleware/            # HTTP middleware
├── .env.example          # Environment variables template
├── .github/              # GitHub Actions workflows
├── Dockerfile           # Container configuration
├── Makefile            # Build and development commands
├── go.mod              # Go module definition
└── README.md           # This file
```

## Development Commands

```bash
# Build the application
make build

# Run the application
make start

# Live reload during development
make live_reload

# Clean dependencies
make clean

# Generate SQL code (if using sqlc)
make sqlc

# Run static analysis
make static_analysis

# Generate new migration
make generate_migration name="create_users_table"

# Generate JWT secret
make generate_32_bit_key
```

## Database Management

### Migrations

```bash
# Create new migration
make generate_migration name="add_user_profile"

# Run migrations
make migrate url="your_database_url"

# Rollback migrations
make migrate_down url="your_database_url"

# Rollback specific number of migrations
make migrate_down_count url="your_database_url" count=2

# Force migration version (use with caution)
make migrate_force url="your_database_url" version=1
```

## Deployment

### Google Cloud Run

The application is configured for automatic deployment to Google Cloud Run using GitHub Actions.

#### Prerequisites

1. **Google Cloud Project** with Cloud Run API enabled
2. **Artifact Registry** repository
3. **Service Account** with appropriate permissions:
   - Artifact Registry Writer
   - Cloud Run Admin
   - Service Account User

#### GitHub Secrets

Configure the following secrets in your GitHub repository:

- `GCP_SERVICE_KEY` - Service account JSON key
- `GCP_PROJECT_ID` - Your Google Cloud project ID
- `DB_PASSWORD` - Database password
- `DB_URL` - Database connection URL
- `SUPABASE_PROJECT_URL` - Supabase project URL
- `SUPABASE_API_KEY` - Supabase API key
- `FRONTEND_URL` - Frontend application URL
- `JWT_SECRET` - JWT signing secret
- `ENV` - Environment (production/staging)

#### Deployment Trigger

Push to the `main` branch to trigger automatic deployment.

## Authentication

The API uses JWT-based authentication. Include the JWT token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### Health Check

```http
GET /
```

Returns the API status and health information.

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Static analysis
make static_analysis
```

## Security

- JWT tokens for authentication
- CORS middleware for cross-origin requests
- Non-root user in Docker container
- Environment-based configuration
- Secure database connections

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation as needed
- Run static analysis before committing
- Use meaningful commit messages

## Technology Stack

- **Language**: Go 1.24+
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Backend Service**: Supabase
- **Containerization**: Docker
- **Deployment**: Google Cloud Run
- **CI/CD**: GitHub Actions
- **Migration Tool**: golang-migrate
- **Code Generation**: sqlc (optional)
- **Live Reload**: Air

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Verify PostgreSQL is running
   - Check database URL in `.env`
   - Ensure database exists

2. **Migration Errors**
   - Check migration files syntax
   - Verify database permissions
   - Use `migrate_force` if needed (carefully)

3. **Build Failures**
   - Ensure Go 1.24+ is installed
   - Run `make install` to refresh dependencies
   - Check for syntax errors

4. **Docker Issues**
   - Verify Docker is running
   - Check Dockerfile syntax
   - Ensure all files are included in build context


## Support

For support and questions:

- Create an issue in the GitHub repository
- Contact the development team
- Check the documentation

---

**Happy coding!**