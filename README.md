# Pizza Backend API

This is the backend API for the Pizza application built with Go, Gin, and PostgreSQL.

## Prerequisites

- Go 1.21 or later
- PostgreSQL 12 or later
- Git

## Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL:
```sql
CREATE DATABASE pizza_db;
```

4. Configure environment variables:
Create a `.env` file in the root directory with the following variables:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_postgres_user
DB_PASSWORD=your_postgres_password
DB_NAME=pizza_db
PORT=8080
```

5. Run database migrations:
```bash
go run cmd/api/main.go migrate
```

6. Run the application:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080` (or the port specified in your .env file).

## Development

### Hot Reload (Optional)
For development with hot reload, you can use [Air](https://github.com/cosmtrek/air):

1. Install Air:
```bash
go install github.com/cosmtrek/air@latest
```

2. Run the application with Air:
```bash
air
```

### Testing
Run the tests:
```bash
go test ./...
```

### Debugging
1. Using VS Code:
   - Install the Go extension
   - Add this launch configuration to `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch API",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/api/main.go"
        }
    ]
}
```

2. Using Delve directly:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug cmd/api/main.go
```

## API Endpoints

### Items

- `GET /items` - Get all items
- `GET /items/:id` - Get a specific item
- `POST /items` - Create a new item
- `PUT /items/:id` - Update an item
- `DELETE /items/:id` - Delete an item

### Invoices

- `GET /invoices` - Get all invoices
- `GET /invoices/:id` - Get a specific invoice
- `POST /invoices` - Create a new invoice
- `PUT /invoices/:id` - Update an invoice
- `DELETE /invoices/:id` - Delete an invoice

## Project Structure

The project follows the MVC (Model-View-Controller) pattern:

```
backend/
├── cmd/
│   └── api/
│       └── main.go       # Application entry point
├── models/
│   ├── item.go          # Item model
│   └── invoice.go       # Invoice model
├── controllers/
│   ├── item_controller.go    # Item controller
│   └── invoice_controller.go # Invoice controller
├── go.mod              # Go module file
├── go.sum              # Go module checksum
└── README.md          # This file
```

## Models

### Item
```go
type Item struct {
    ID          uint    `json:"id"`
    Name        string  `json:"name"`
    Category    string  `json:"category"`
    Price       float64 `json:"price"`
    Description string  `json:"description"`
}
```

### Invoice
```go
type Invoice struct {
    ID            uint          `json:"id"`
    InvoiceNumber string        `json:"invoice_number"`
    CustomerName  string        `json:"customer_name"`
    CustomerPhone string        `json:"customer_phone"`
    Date         time.Time      `json:"date"`
    Items        []InvoiceItem  `json:"items"`
    Total        float64        `json:"total"`
}
``` 