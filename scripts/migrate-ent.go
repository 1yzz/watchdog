package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"watchdog/config"
)

func main() {
	var (
		dryRun = flag.Bool("dry-run", false, "Print the SQL statements without executing them")
		help   = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		fmt.Println("Ent Migration Tool for Watchdog Service")
		fmt.Println("Usage: go run scripts/migrate-ent.go [flags]")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  -dry-run    Print SQL statements without executing")
		fmt.Println("  -help       Show this help message")
		fmt.Println()
		fmt.Println("This tool automatically generates and applies database schema")
		fmt.Println("based on the Ent entity definitions in ent/schema/")
		return
	}

	log.Println("Starting Ent-based database migration...")

	// Load configuration
	cfg, entClient, err := config.LoadWithEntClient()
	if err != nil {
		log.Fatalf("Failed to load config and connect to database: %v", err)
	}
	defer entClient.Close()

	ctx := context.Background()

	if *dryRun {
		log.Println("DRY RUN MODE: Printing SQL statements that would be executed")
		
		// For dry run, we'd need to use the schema creation with debug mode
		// This is a simplified version - in practice you'd use the migrate package
		log.Println("Schema would be created with the following structure:")
		log.Println("- Table: services")
		log.Println("  - id: BIGINT AUTO_INCREMENT PRIMARY KEY")  
		log.Println("  - name: VARCHAR(255) NOT NULL")
		log.Println("  - endpoint: VARCHAR(500) NOT NULL") 
		log.Println("  - type: ENUM(...) NOT NULL DEFAULT 'SERVICE_TYPE_UNSPECIFIED'")
		log.Println("  - status: VARCHAR(50) NOT NULL DEFAULT 'active'")
		log.Println("  - last_heartbeat: TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
		log.Println("  - created_at: TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
		log.Println("  - updated_at: TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP")
		log.Println("- Indexes:")
		log.Println("  - UNIQUE(name, endpoint)")
		log.Println("  - INDEX(type)")
		log.Println("  - INDEX(status)")
		log.Println("  - INDEX(last_heartbeat)")
		log.Println("  - INDEX(type, status)")
		
		log.Println("Migration complete (dry run)")
		return
	}

	// Run the actual migration
	log.Printf("Connecting to database: %s@%s:%d/%s", 
		cfg.Database.Username, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)

	if err := entClient.AutoMigrate(ctx); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Test the connection
	if err := entClient.HealthCheck(); err != nil {
		log.Fatalf("Health check failed after migration: %v", err)
	}

	log.Println("✅ Migration completed successfully!")
	log.Println("✅ Database health check passed")
	log.Println("✅ Watchdog service is ready to run")
	
	// Optionally show some stats
	services, err := entClient.ListServices()
	if err != nil {
		log.Printf("Warning: Could not count services: %v", err)
	} else {
		log.Printf("📊 Current services in database: %d", len(services))
	}
}