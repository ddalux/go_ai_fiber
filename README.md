Go Fiber + SQLite example

Endpoints:
- POST /register {email,password,firstname,lastname,phone,birthday}
- POST /login {email,password} -> returns {token}
- GET /me (Authorization: Bearer <token>) -> returns user info
- GET /swagger -> API docs (ReDoc)

Run (Go 1.21.6):

Go Fiber + SQLite example

Endpoints:
- POST /register {email,password,firstname,lastname,phone,birthday}
- POST /login {email,password} -> returns {token}
- GET /me (Authorization: Bearer <token>) -> returns user info
- GET /swagger -> API docs (ReDoc)

Run (Go 1.21.6):

```bash
go mod tidy
go run ./cmd/server
```

Set JWT secret:

```bash
export JWT_SECRET="your-secret-here"
```

DB file: `app.db` (SQLite, created automatically)

New features (implemented from screenshot 3 - transfer flow):

- POST /transfer {to, amount, note} (Authorization: Bearer <token>)
	- `to` accepts member code (e.g., LBK001234) or recipient email
	- transfers points from authenticated user to recipient

- GET /contacts/recent
	- returns recent recipients the authenticated user has transferred to

Swagger / API docs

- Open http://localhost:3000/swagger for ReDoc UI
- Open http://localhost:3000/swagger/doc.json for the OpenAPI JSON

Transfer example

```bash
# assume $TOKEN contains a valid Bearer token from /login
curl -X POST http://localhost:3000/transfer \
	-H 'Content-Type: application/json' \
	-H "Authorization: Bearer $TOKEN" \
	-d '{"to":"LBK001234","amount":100,"note":"gift"}'
```

Notes and caveats

- Password hashing uses a simple salt+sha256 for demo compatibility with Go 1.21.6; replace with bcrypt in production.
- Member codes are stored in `member_code` (may be empty for users created without it).
- Transactions are recorded in a `transactions` table.

Next steps

- Add unit tests for transfer and recent-contacts.
- Improve validation and error handling for production readiness.
