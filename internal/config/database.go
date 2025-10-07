package config

type databaseConfig struct {
	RedisAddress          string
	MongoUri              string
	MongoDatabase         string
	MongoTimeout          int
	QdrantHost            string
	QdrantPort            int
	DefaultVectorSize     uint64
	DefaultCollectionName string
}

func newDatabaseConfig() databaseConfig {
	mongoUri := getString("MONGO_URI", "mongodb://localhost:27017")
	mongoDatabase := getString("MONGO_DATABASE", "better-mem")
	mongoTimeout := getInt("MONGO_TIMEOUT", 10)
	qdrantHost := getString("QDRANT_HOST", "localhost")
	qdrantPort := getInt("QDRANT_PORT", 6334)
	defaultVectorSize := getInt("QDRANT_DEFAULT_VECTOR_SIZE", 384)
	defaultCollectionName := getString("QDRANT_DEFAULT_COLLECTION_NAME", "better-mem-default")
	redisAddress := getString("REDIS_ADDRESS", "localhost:6379")
	return databaseConfig{
		RedisAddress:          redisAddress,
		MongoUri:              mongoUri,
		MongoDatabase:         mongoDatabase,
		MongoTimeout:          mongoTimeout,
		QdrantHost:            qdrantHost,
		QdrantPort:            qdrantPort,
		DefaultVectorSize:     uint64(defaultVectorSize),
		DefaultCollectionName: defaultCollectionName,
	}
}

var Database = newDatabaseConfig()
