# Database Migrations

This directory contains SQL migration files for the database schema.

## Migration Files

- `001_create_users_table.sql` - Creates the users table with all required fields

## Running Migrations

### Option 1: Using psql

```bash
psql -h 127.0.0.1 -p 5432 -U postgres -d bm_staff -f migrations/001_create_users_table.sql
```

### Option 2: Using a migration tool

You can use tools like:
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [sql-migrate](https://github.com/rubenv/sql-migrate)
- [atlas](https://atlasgo.io/)

### Option 3: Manual execution

Connect to your PostgreSQL database and execute the SQL files in order.

## Migration Order

Migrations should be executed in numerical order:
1. `001_create_users_table.sql`

## Rollback

To rollback a migration, you can create a corresponding rollback file or manually reverse the changes.

