# Bachelor Backend API

A Go-based REST API using Fiber v3 framework for the Bachelor project. This API provides endpoints for user management, product catalog, cart operations, and analytics.

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 14 or higher
- Make (optional, for using Makefile commands)

### Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd bachelor_new/back
```

2. **Set up environment variables**
Create a `.env` file in the `back` directory:
```env
# Application
APP_PORT=8080
APP_ENV=development  # or production

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db

# JWT
JWT_SECRET=your-secure-secret-key
```

3. **Install dependencies**
```bash
go mod download
```

4. **Run the application**
```bash
go run main.go
```

The server will start at `http://localhost:8080` (or your configured APP_PORT).

## 🔒 Authentication

The API uses JWT (JSON Web Token) for authentication. Protected routes require a valid JWT token in the Authorization header.

### Token Format
```
Authorization: Bearer <your-jwt-token>
```

## 🛡️ CSRF Protection

The API implements CSRF protection for all non-GET requests. Clients must:

1. Get a CSRF token from `/api/csrf-token` (GET request)
2. Include the token in the `X-Csrf-Token` header for all POST/PUT/PATCH/DELETE requests

Example:
```javascript
// Get CSRF token
const response = await fetch('http://localhost:8080/api/csrf-token', {
  credentials: 'include'
});
const { csrfToken } = await response.json();

// Use token in subsequent requests
await fetch('http://localhost:8080/api/products', {
  method: 'POST',
  credentials: 'include',
  headers: {
    'X-Csrf-Token': csrfToken,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(data)
});
```

## 📚 API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### Users
- `GET /api/users` - List all users (Admin only)
- `GET /api/users/:id` - Get user details
- `GET /api/users/:id/activities` - Get user activities
- `GET /api/users/:id/favorites` - Get user favorites
- `GET /api/users/:id/cart` - Get user's cart
- `GET /api/users/:id/orders` - Get user's orders

### Products
- `GET /api/products` - List all products
- `GET /api/products/:id` - Get product details
- `GET /api/products/category/:category` - Get products by category

### Cart
- `GET /api/cart/:id` - Get cart details
- `POST /api/cart/:id/items` - Add item to cart
- `PUT /api/cart/:id/items/:itemId` - Update cart item
- `DELETE /api/cart/:id/items/:itemId` - Remove item from cart

### Orders
- `GET /api/orders` - List all orders (Admin only)
- `GET /api/orders/:id` - Get order details

### Analytics (Admin only)
- `GET /api/analytics/activities` - Get activity analytics
- `GET /api/analytics/requests` - Get request analytics
- `GET /api/analytics/products/popular` - Get popular products
- `GET /api/analytics/users/active` - Get active users

## 📊 Data Models

### User
```json
{
  "id": "uint",
  "email": "string",
  "name": "string",
  "role": "string"
}
```

### Product
```json
{
  "id": "uint",
  "name": "string",
  "description": "string",
  "price": "float64",
  "stock": "int",
  "category": "string",
  "image_url": "string"
}
```

### Cart
```json
{
  "id": "uint",
  "user_id": "uint",
  "items": [
    {
      "id": "uint",
      "product_id": "uint",
      "quantity": "int",
      "product": "Product"
    }
  ]
}
```

### Order
```json
{
  "id": "uint",
  "user_id": "uint",
  "status": "string",
  "total": "float64",
  "items": [
    {
      "id": "uint",
      "product_id": "uint",
      "quantity": "int",
      "price": "float64",
      "product": "Product"
    }
  ]
}
```

## 🔧 Development

### Database Migrations
The application automatically handles database migrations on startup. In development mode (`APP_ENV=development`), it will drop and recreate all tables.

### Error Handling
All errors are returned in a consistent JSON format:
```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes
- 200: Success
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

## 🔐 Security Features

1. **JWT Authentication**
   - Tokens expire after 24 hours
   - Protected routes require valid tokens

2. **CSRF Protection**
   - Double Submit Cookie Pattern
   - Required for all non-GET requests
   - Token expires after 30 minutes of inactivity

3. **CORS**
   - Configurable allowed origins
   - Credentials supported
   - Proper header exposure

4. **Password Security**
   - Passwords are hashed using bcrypt
   - Never stored or transmitted in plain text

## 📝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
