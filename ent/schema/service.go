package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Service holds the schema definition for the Service entity.
type Service struct {
	ent.Schema
}

// Fields of the Service.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		// Auto-incrementing primary key
		field.Int64("id").
			Positive().
			Comment("Service unique identifier (auto-increment)").
			Annotations(entsql.Annotation{
				Options: "AUTO_INCREMENT",
			}),

		field.String("name").
			MaxLen(255).
			NotEmpty().
			Comment("Service name"),

		field.String("endpoint").
			MaxLen(500).
			NotEmpty().
			Comment("Service endpoint URL"),

		field.Enum("type").
			Values(
				"SERVICE_TYPE_UNSPECIFIED",
				"SERVICE_TYPE_HTTP", 
				"SERVICE_TYPE_GRPC",
				"SERVICE_TYPE_DATABASE",
				"SERVICE_TYPE_CACHE",
				"SERVICE_TYPE_QUEUE",
				"SERVICE_TYPE_STORAGE",
				"SERVICE_TYPE_EXTERNAL_API",
				"SERVICE_TYPE_MICROSERVICE", 
				"SERVICE_TYPE_OTHER",
			).
			Default("SERVICE_TYPE_UNSPECIFIED").
			Comment("Type of service being monitored"),

		field.String("status").
			MaxLen(50).
			Default("active").
			Comment("Current service status"),

		field.Time("last_heartbeat").
			Default(time.Now).
			Comment("Timestamp of last heartbeat/update").
			Annotations(entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
				Options: "ON UPDATE CURRENT_TIMESTAMP",
			}),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("When the service was first registered").
			Annotations(entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
			}),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("When the service record was last updated").
			Annotations(entsql.Annotation{
				Default: "CURRENT_TIMESTAMP",
				Options: "ON UPDATE CURRENT_TIMESTAMP",
			}),
	}
}

// Edges of the Service.
func (Service) Edges() []ent.Edge {
	return nil
}

// Indexes of the Service.
func (Service) Indexes() []ent.Index {
	return []ent.Index{
		// Unique constraint on name and endpoint combination
		index.Fields("name", "endpoint").
			Unique(),

		// Index on type for filtering services by type
		index.Fields("type"),

		// Index on status for filtering by status 
		index.Fields("status"),

		// Index on last_heartbeat for time-based queries
		index.Fields("last_heartbeat"),

		// Composite index for common queries (type + status)
		index.Fields("type", "status"),
	}
}

// Annotations of the Service.
func (Service) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "services",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_unicode_ci",
			Options:   "ENGINE=InnoDB",
		},
	}
}