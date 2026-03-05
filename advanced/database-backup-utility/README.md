# Database Backup Utility

A powerful command-line tool for backing up and restoring databases with support for multiple database types, cloud storage, compression, and automated scheduling.

## Features

- **Multiple Database Support**: MySQL, PostgreSQL, MongoDB, and SQLite
- **Backup Types**: Full, incremental, and differential backups
- **Storage Options**: Local storage, AWS S3, and Google Cloud Storage (GCS)
- **Compression**: Gzip compression to save storage space
- **Selective Restore**: Restore specific tables or collections instead of entire databases
- **Notifications**: Slack notifications for backup success or failure
- **Automated Scheduling**: Built-in scheduler for automatic backups
- **Retention Policies**: Automatically clean up old backups based on configurable limits

## Prerequisites

- Go 1.21+
- Database client tools installed on your system:
  - `mysqldump` and `mysql` for MySQL
  - `pg_dump` and `psql` for PostgreSQL
  - `mongodump` and `mongorestore` for MongoDB
  - `sqlite3` for SQLite

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/jaygaha/roadmap-go-projects.git
   cd advanced/database-backup-utility
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the utility:
   ```bash
   go build -o dbu ./cmd/database-backup-utility/main.go
   ```

## Quick start

### 1. Configure the utility

Copy the example environment file and configuration:

```bash
cp .env.example .env
```

Edit `config.yaml` to define your databases, storage settings, and notification hooks.

### 2. Test database connectivity

Before running backups, test your database connections:

```bash
./dbu test
```

### 3. Run a backup

```bash
# Backup all databases defined in config
./dbu backup

# Backup specific databases
./dbu backup mydb_mysql mydb_postgres

# Force a specific backup type
./dbu backup --type incremental
```

### 4. Restore from backup

```bash
# Full restore
./dbu restore mydb_mysql --file ./backups/mydb_mysql_full_backup_20260301.sql

# Selective table restore (MySQL/PostgreSQL)
./dbu restore mydb_postgres --file ./backups/backup.sql --tables "users,orders"

# Selective collection restore (MongoDB)
./dbu restore mydb_mongo --file ./backups/backup.archive --collections "users"
```

### 5. Automated scheduling

Run the utility in daemon mode to execute backups on a schedule:

```bash
./dbu schedule --config config.yaml

# Override schedule via flags
./dbu schedule --cron "0 0 * * *" --timezone "UTC"
```

## Configuration example

```yaml
databases:
  - name: "mydb_postgres"
    type: "postgres"
    host: "localhost"
    user: "admin"
    password: "secret"
    port: 5432

  - name: "mydb_mysql"
    type: "mysql"
    host: "localhost"
    user: "root"
    password: "password"
    port: 3306

backup:
  type: "incremental"
  path: "./backups"
  retain: 5
  compress: true

storage:
  type: "s3"
  bucket: "my-backups"
  region: "us-east-1"

notification:
  slack:
    enabled: true
    webhook_url: "https://hooks.slack.com/services/..."
```

## Command reference

### Backup commands

| Command | Description |
|---------|-------------|
| `./dbu backup` | Backup all databases from config |
| `./dbu backup <name>` | Backup a specific database |
| `./dbu backup --type full` | Force full backup type |
| `./dbu backup --type incremental` | Force incremental backup |

### Restore commands

| Command | Description |
|---------|-------------|
| `./dbu restore <name> --file <path>` | Restore database from file |
| `./dbu restore <name> --file <path> --tables "t1,t2"` | Restore specific tables |
| `./dbu restore <name> --file <path> --collections "c1,c2"` | Restore specific collections |

### Other cmmands

| Command | Description |
|---------|-------------|
| `./dbu test` | Test database connectivity |
| `./dbu schedule` | Run in scheduled/daemon mode |
| `./dbu --help` | Show help for all commands |

## Project structure

```
database-backup-utility/
├── cmd/
│   ├── backup.go              # Backup command definitions
│   ├── restore.go             # Restore command definitions
│   ├── schedule.go            # Schedule command definitions
│   ├── test.go                # Test command definitions
│   ├── root.go                # Root command setup
│   └── database-backup-utility/
│       └── main.go            # Application entry point
├── internal/
│   ├── backup/
│   │   ├── backup.go          # Core backup interface
│   │   ├── mysql.go           # MySQL backup implementation
│   │   ├── postgresql.go      # PostgreSQL backup implementation
│   │   ├── mongodb.go         # MongoDB backup implementation
│   │   └── sqlite.go          # SQLite backup implementation
│   ├── config/
│   │   ├── config.go          # Configuration structures
│   │   └── load.go            # Configuration loading
│   ├── notification/
│   │   └── slack.go           # Slack notification service
│   ├── scheduler/
│   │   └── scheduler.go       # Cron-based scheduler
│   ├── storage/
│   │   ├── s3.go              # AWS S3 storage
│   │   └── gcs.go             # Google Cloud Storage
│   └── utils/
│       ├── compress.go        # Compression utilities
│       ├── logger.go          # Logging utilities
│       └── validation.go      # Input validation
├── config.yaml                # Configuration file
├── .env.example               # Example environment variables
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## How it works

### Backup flow

1. **Read Configuration**: Load database credentials and settings from `config.yaml`
2. **Connect to Database**: Establish connection using provided credentials
3. **Execute Backup**: Run database-specific backup command (e.g., `mysqldump`, `pg_dump`)
4. **Compress**: Optionally compress the backup file using gzip
5. **Upload**: Transfer backup to configured storage (local, S3, or GCS)
6. **Notify**: Send Slack notification on success or failure

### Restore flow

1. **Locate Backup File**: Find the backup file to restore
2. **Decompress**: Decompress if the file is gzipped
3. **Connect to Database**: Establish connection to target database
4. **Execute Restore**: Run database-specific restore command
5. **Verify**: Optionally verify the restore was successful

## Troubleshooting

### Common issues

**Database connection failed**

- Verify database credentials in `config.yaml`
- Ensure the database server is running
- Check firewall rules allow connections

**Command not found errors**

- Install required database client tools (see Prerequisites)
- Ensure client tools are in your PATH

**Permission denied for S3/GCS**

- Verify AWS/GCP credentials are set correctly
- Check IAM permissions for the storage bucket

## Contributing

- Challenge: [Database Backup Utility](https://roadmap.sh/projects/database-backup-utility)
- This project is part of the [roadmap.sh](https://roadmap.sh/projects) backend projects series.
- Created by [jaygaha](https://github.com/jaygaha)
