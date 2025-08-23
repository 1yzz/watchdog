# Configuration Guide

The Watchdog gRPC service uses a flexible configuration system that supports multiple sources with clear precedence rules.

## Configuration Sources (Priority Order)

The configuration is loaded in the following order, with higher-numbered items taking precedence:

1. **Default values** (lowest priority) - Built-in fallback values
2. **System environment variables** - Set via `export VAR=value`
3. **`.env` file** - Project-level environment file
4. **`.env.local` file** (highest priority) - Local overrides, git-ignored

## Environment Files

### Supported Files

- **`.env.default`** - Template file (tracked in git)
- **`.env`** - Main environment file (git-ignored)
- **`.env.local`** - Local development overrides (git-ignored)

### Setup Process

1. **Create your environment file:**
   ```bash
   make env-setup
   ```

2. **Edit the created `.env` file:**
   ```bash
   vim .env  # or your preferred editor
   ```

3. **For local development overrides, create `.env.local`:**
   ```bash
   cp .env .env.local
   # Edit .env.local with your specific local settings
   ```

## Configuration Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `50051` | gRPC server listening port |
| `DB_HOST` | `localhost` | MySQL server hostname or IP |
| `DB_PORT` | `3306` | MySQL server port |
| `DB_USERNAME` | `watchdog` | MySQL username |
| `DB_PASSWORD` | `watchdog123` | MySQL password |
| `DB_DATABASE` | `watchdog_db` | MySQL database name |

## Example Configurations

### Local Development
```bash
# .env
PORT=50051
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=watchdog
DB_PASSWORD=watchdog123
DB_DATABASE=watchdog_db
```

### Production with External Database
```bash
# .env
PORT=50051
DB_HOST=prod-mysql.example.com
DB_PORT=3306
DB_USERNAME=prod_watchdog
DB_PASSWORD=secure_production_password
DB_DATABASE=watchdog_production
```

### Remote MySQL Server
```bash
# .env
PORT=50051
DB_HOST=mysql.example.com
DB_PORT=3306
DB_USERNAME=watchdog
DB_PASSWORD=watchdog123
DB_DATABASE=watchdog_db
```

### AWS RDS
```bash
# .env
PORT=50051
DB_HOST=myinstance.cxkpkp6lmgqb.us-west-2.rds.amazonaws.com
DB_PORT=3306
DB_USERNAME=admin
DB_PASSWORD=your_secure_password
DB_DATABASE=watchdog_db
```

### Google Cloud SQL
```bash
# .env
PORT=50051
DB_HOST=10.1.2.3
DB_PORT=3306
DB_USERNAME=watchdog-user
DB_PASSWORD=your_secure_password
DB_DATABASE=watchdog_db
```

## Security Best Practices

1. **Never commit sensitive files:**
   - `.env` and `.env.local` are git-ignored by default
   - Only commit `.env.example` or `.env.default`

2. **Use strong passwords:**
   - Generate complex passwords for production databases
   - Consider using environment-specific credentials

3. **Restrict database access:**
   - Use dedicated database users with minimal required permissions
   - Configure database firewall rules appropriately

4. **Local development isolation:**
   - Use `.env.local` for personal development settings
   - This file has highest priority and is never committed

## Troubleshooting

### Environment File Not Loading
- Check file exists in the current working directory
- Verify file permissions are readable
- Check for syntax errors in .env file
- Review server startup logs for loading messages

### Database Connection Issues
- Test connection with: `make db-test`
- Verify MySQL server is running and accessible
- Check firewall rules and network connectivity
- Validate credentials and database exists

### Configuration Precedence Issues
- Remember the priority order: `.env.local` > `.env` > system env > defaults
- Use `make db-test` to see which configuration is being loaded
- Check server startup logs for configuration source information

## Testing Configuration

Test your database connection:
```bash
make db-test
```

This will:
1. Load your environment configuration
2. Attempt to connect to the database
3. Verify table accessibility
4. Report any issues found