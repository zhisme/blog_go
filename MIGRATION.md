# CSV to SQLite Migration Guide

This guide explains how to migrate your mailing list data from CSV format to SQLite database.

## Overview

The blog backend has been upgraded to use SQLite instead of CSV for data storage. This provides:

- **Better Performance**: No more full-file reads on every save operation
- **Concurrent Safety**: ACID transactions and proper locking
- **Data Integrity**: Unique constraints enforced at database level
- **Reliability**: No file corruption risks
- **Scalability**: Ready for future growth

## Migration Steps

### 1. Backup Your CSV File (Optional)

If you have existing data in `mailing_list.csv`, create a backup first:

```bash
cp mailing_list.csv mailing_list.csv.backup
```

### 2. Run the Migration Script

The migration script will import all data from your CSV file into SQLite:

```bash
go run cmd/migrate/main.go
```

**Options:**
- `-csv`: Path to CSV file (default: `mailing_list.csv`)
- `-db`: Path to SQLite database (default: `blog.db`)

**Example with custom paths:**
```bash
go run cmd/migrate/main.go -csv=/path/to/data.csv -db=/path/to/database.db
```

### 3. Verify Migration

The migration script will display:
- Total records processed
- Successfully imported records
- Skipped duplicates
- Any errors encountered

Example output:
```
=== Migration Complete ===
Total records processed: 100
Successfully imported: 100
Skipped (duplicates): 0
Errors: 0
```

### 4. Update Your Environment (Optional)

You can customize the database path using environment variables:

```bash
export DB_PATH=/path/to/blog.db
```

Default path: `blog.db` in the application root directory

### 5. Test Your Application

Start your server and verify everything works:

```bash
go run cmd/api/main.go
```

The application will now use the SQLite database instead of CSV.

## Database Schema

The SQLite database uses the following schema:

```sql
CREATE TABLE mailing_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_mailing_list_email ON mailing_list(email);
CREATE INDEX idx_mailing_list_created_at ON mailing_list(created_at);
```

## Features

### Duplicate Handling

The migration script automatically handles duplicates:
- Duplicate emails are detected by unique constraint
- Duplicates are logged but don't cause errors
- Migration continues with remaining records

### Concurrent Access

SQLite is configured with:
- **WAL mode** (Write-Ahead Logging) for better concurrent read/write performance
- **Foreign keys** enabled
- Proper ACID transaction support

### Zero-Downtime Migration

You can run the migration script multiple times safely:
- Existing records are preserved
- Only new records are added
- No data loss

## Troubleshooting

### "CSV file does not exist"
Make sure your CSV file exists at the specified path. Default is `mailing_list.csv` in the current directory.

### "Failed to initialize database"
Check that you have write permissions in the directory where the database will be created.

### Migration Errors
The migration script will continue even if some records fail. Check the console output for specific error messages.

## Rollback (If Needed)

If you need to rollback to CSV:

1. Stop your application
2. Restore your CSV backup:
   ```bash
   cp mailing_list.csv.backup mailing_list.csv
   ```
3. Revert code changes using git
4. Restart application

## Configuration

### Environment Variables

- `DB_PATH`: Path to SQLite database file (default: `blog.db`)
- `SERVER_ADDR`: Server listen address (default: `:8080`)

### File Locations

- Database: `blog.db` (configurable via `DB_PATH`)
- Migration script: `cmd/migrate/main.go`
- Schema: `migrations/001_create_mailing_list.sql`

## Testing

Run the test suite to verify everything works:

```bash
go test ./...
```

All tests use in-memory SQLite databases for isolation and speed.

## Notes

- The CSV file is no longer used after migration
- You can safely delete the CSV file after verifying the migration
- The SQLite database file should be added to `.gitignore`
- Consider regular backups of the `.db` file

## Support

For issues or questions, please create an issue in the repository.
