Go Fiber + SQLite example

Endpoints:
- POST /register {email,password,firstname,lastname,phone,birthday}
- POST /login {email,password} -> returns {token}
- GET /me (Authorization: Bearer <token>) -> returns user info
- GET /swagger -> API docs (ReDoc)

Run (Go 1.21.6):

```bash
go mod tidy
go run main.go
```

Set JWT secret:

```bash
export JWT_SECRET="your-secret-here"
```

DB file: `app.db` (SQLite, created automatically)
