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

4. Run the application:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`.

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