# Personal Vault Server

A robust, production-ready backend API server built with Go, Gin, and PostgreSQL for a personal document management system.

## 🚀 Project Overview

Personal Vault Server is a RESTful API backend that provides secure file storage, user authentication, and document management capabilities. It's designed as a learning project that demonstrates modern Go development practices, clean architecture, and production-ready patterns.

## 🏗️ Architecture & Design Patterns

### Clean Architecture Implementation
```
cmd/server/          # Application entry point
internal/
├── app/
│   ├── entities/     # Domain models
│   ├── handlers/    # HTTP request handlers
│   ├── repositories/ # Data access layer
│   ├── routes/      # Route definitions
│   └── services/    # Business logic layer
├── config/          # Configuration management
├── helper/          # Utility functions
└── middleware/      # HTTP middleware
pkg/database/        # Database connection and setup
```

### Key Design Patterns Demonstrated

1. **Repository Pattern**: Abstracts data access logic
2. **Service Layer Pattern**: Encapsulates business logic
3. **Dependency Injection**: Clean separation of concerns
4. **Middleware Pattern**: Cross-cutting concerns (auth, CORS)
5. **Configuration Pattern**: Environment-based configuration

## 🛠️ Technology Stack

- **Language**: Go 1.25+
- **Web Framework**: Gin (HTTP router and middleware)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Configuration**: Environment variables with godotenv
- **CORS**: Cross-origin resource sharing support

## 📋 Features

### Core Functionality
- **User Authentication**: JWT-based login/register system
- **File Management**: Upload, download, preview, and organize files
- **Folder Operations**: Create, navigate, and manage directory structures
- **Security**: Password hashing, JWT tokens, CORS protection
- **Database Integration**: PostgreSQL with GORM for data persistence

### API Endpoints

#### Authentication
- `POST /auth/login` - User login
- `POST /auth/register` - User registration

#### File Management (Protected Routes)
- `GET /api/drivers/root` - Get root directory contents
- `GET /api/drivers/list` - List files in a directory
- `POST /api/drivers/upload` - Upload files
- `POST /api/drivers/download` - Download files
- `POST /api/drivers/create-folder` - Create new folder
- `GET /api/drivers/preview` - Preview file contents
- `GET /api/drivers/stream` - Stream file content

#### User Management
- `GET /api/users/me` - Get current user details

## 🚀 Getting Started

### Prerequisites
- Go 1.25 or higher
- PostgreSQL database
- Git

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd PersonalVaultServer
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
Create a `.env` file in the root directory:
```env
# Server Configuration
SERVER_PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=personal_vault
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_here
JWT_EXPIRES_IN_HOURS=24

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
CORS_EXPOSE_HEADERS=Content-Length
```

4. **Set up PostgreSQL database**
```sql
CREATE DATABASE personal_vault;
```

5. **Run the application**
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## 🧪 Testing the API

### Using curl

**Register a new user:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

**Access protected route (replace TOKEN with actual JWT):**
```bash
curl -X GET http://localhost:8080/api/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📚 What You Can Learn

### Go Development Best Practices
1. **Project Structure**: Clean, scalable project organization
2. **Error Handling**: Proper error propagation and logging
3. **Configuration Management**: Environment-based configuration
4. **Database Integration**: GORM ORM usage and migrations
5. **Authentication**: JWT implementation and middleware
6. **API Design**: RESTful endpoint design and HTTP status codes

### Software Architecture Concepts
1. **Separation of Concerns**: Clear boundaries between layers
2. **Dependency Injection**: Loose coupling between components
3. **Repository Pattern**: Data access abstraction
4. **Service Layer**: Business logic encapsulation
5. **Middleware Pattern**: Cross-cutting concerns

### Production-Ready Features
1. **Security**: Password hashing, JWT tokens, CORS
2. **Database**: Connection pooling, migrations, ORM
3. **Configuration**: Environment variables, flexible settings
4. **Error Handling**: Graceful error responses
5. **Logging**: Structured logging for debugging

### Advanced Go Concepts
1. **Interfaces**: Go interface usage for dependency injection
2. **Struct Embedding**: Composition over inheritance
3. **Context**: Request context handling
4. **Goroutines**: Concurrent programming (where applicable)
5. **Package Management**: Go modules and dependency management

## 🔧 Development

### Project Structure Explanation

```
cmd/server/main.go          # Application entry point and initialization
internal/
├── app/
│   ├── entities/           # Domain models (User, Request, Response)
│   ├── handlers/           # HTTP handlers (Auth, User, Drive)
│   ├── repositories/       # Data access layer
│   ├── routes/             # Route definitions and middleware setup
│   └── services/           # Business logic layer
├── config/                 # Configuration management
├── helper/                # Utility functions (JWT, password hashing)
└── middleware/             # HTTP middleware (authentication)
pkg/database/               # Database connection and setup
```

### Adding New Features

1. **New Entity**: Add to `internal/app/entities/`
2. **New Repository**: Add to `internal/app/repositories/`
3. **New Service**: Add to `internal/app/services/`
4. **New Handler**: Add to `internal/app/handlers/`
5. **New Route**: Add to `internal/app/routes/`

### Database Migrations

The application uses GORM's AutoMigrate feature. When you add new fields to entities, they will be automatically migrated on startup.

## 🚀 Deployment

### Environment Variables for Production
```env
ENV=production
SERVER_PORT=8080
DB_HOST=your_production_db_host
DB_USER=your_production_db_user
DB_PASSWORD=your_production_db_password
DB_NAME=personal_vault_prod
JWT_SECRET=your_production_jwt_secret
```

### Docker Deployment (Optional)
```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM ORM](https://gorm.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [JWT-Go](https://github.com/golang-jwt/jwt)

---

**Happy Coding! 🚀**

This project serves as an excellent learning resource for Go developers who want to understand modern backend development patterns, clean architecture, and production-ready API design.