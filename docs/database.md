# Database ER Diagram

The diagram below describes the database models discovered in `internal/repository/user_repo.go` and their relationships.

```mermaid
erDiagram
    USERS {
        uint ID PK "primary key"
        string email "unique, not null"
        string member_code "unique, nullable"
        string password_hash
        string first_name
        string last_name
        string phone
        string birthday
        datetime created_at
        int points
    }

    TRANSACTIONS {
        uint ID PK "primary key"
        string from_email
        string to_email
        string to_member_code
        int amount
        string note
        datetime created_at
    }

    %% Relationships:
    %% Transactions reference Users by email for sender (from_email) and recipient (to_email).
    USERS ||--o{ TRANSACTIONS : "sent_transactions (from_email)"
    USERS ||--o{ TRANSACTIONS : "received_transactions (to_email)"

```

Notes
- The repository stores `User` records and `Transaction` records. Transactions link to users by email (fields `from_email` and `to_email`) and optionally include `to_member_code` for member-code based lookups.
- In GORM this is implemented in `internal/repository/user_repo.go` where `Transaction` and `User` are separate tables; foreign-key constraints are not explicitly declared in the current code — relationships are enforced/queried by application logic (email/member_code lookups).

Reference
- Mermaid ER syntax: https://mermaid.js.org/syntax/entityRelationshipDiagram.html
