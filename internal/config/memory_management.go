package config

// ShortTermMemoryValidationConfig is the configuration
// that defines how a short term memory is validated
// (for promotion and for discard)
type shortTermMemoryValidationConfig struct {
	// The age limit in hours
	AgeLimitHours int
	// The minimal relevancy to be considered for long term memory
	MinimalRelevancyForPromotion int
	// The minimal relevancy to not be discarded
	MinimalRelevancyForDiscard int
	// The relevancy threshold to be considered for long term memory
	LongTermThreshold float32
}

type memoryManagementConfig struct {
	ManageSTMemoryTaskPeriod    string  // = "@every 10m"
	ManageSTMemoryTaskPeriodInt int     // = 10 * 60
	MemorySimilarityThreshold   float32 // = 0.9
	MaxSimultaneousTasks        int
	STValConfig                 *shortTermMemoryValidationConfig
}

func newMemoryManagementConfig() *memoryManagementConfig {
	manageShortTermMemoryTaskPeriod := getString("MEMORY_MANAGEMENT_MANAGE_SHORT_TERM_MEMORY_TASK_PERIOD", "@every 30s")
	manageShortTermMemoryTaskPeriodInt := getInt("MEMORY_MANAGEMENT_MANAGE_SHORT_TERM_MEMORY_TASK_PERIOD_SECONDS_INT", 10*60)
	memorySimilarityThreshold := getFloat32("MEMORY_MANAGEMENT_MEMORY_SIMILARITY_THRESHOLD", 0.9)
	maxSimultaneousTasks := getInt("MEMORY_MANAGEMENT_MAX_SIMULTANEOUS_TASKS", 10)
	ageLimit := getInt("MEMORY_MANAGEMENT_SHORT_TERM_MEMORY_AGE_LIMIT", 24*7)
	minimalRelevancyForPromotion := getInt("MEMORY_MANAGEMENT_SHORT_TERM_MEMORY_MINIMAL_RELEVANCY_FOR_PROMOTION", 10)
	minimalRelevancyForDiscard := getInt("MEMORY_MANAGEMENT_SHORT_TERM_MEMORY_MINIMAL_RELEVANCY_FOR_DISCARD", 5)
	longTermThreshold := getFloat32("MEMORY_MANAGEMENT_SHORT_TERM_MEMORY_LONG_TERM_THRESHOLD", 0.5)
	return &memoryManagementConfig{
		ManageSTMemoryTaskPeriod:    manageShortTermMemoryTaskPeriod,
		ManageSTMemoryTaskPeriodInt: manageShortTermMemoryTaskPeriodInt,
		MemorySimilarityThreshold:   memorySimilarityThreshold,
		MaxSimultaneousTasks:        maxSimultaneousTasks,
		STValConfig: &shortTermMemoryValidationConfig{
			AgeLimitHours:                ageLimit,
			MinimalRelevancyForPromotion: minimalRelevancyForPromotion,
			MinimalRelevancyForDiscard:   minimalRelevancyForDiscard,
			LongTermThreshold:            longTermThreshold,
		},
	}
}

var MemoryManagement = newMemoryManagementConfig()
