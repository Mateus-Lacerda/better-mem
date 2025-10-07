package config

type workerConfig struct {
	// Max retry for tasks
	MaxRetry int
	// Timeout for task in seconds
	Timeout int
	// Concurrency for task
	Concurrency int
}

func newWorkerConfig() *workerConfig {
	maxRetry := getInt("WORKER_MAX_RETRY", 5)
	timeout := getInt("WORKER_TIMEOUT", 60)
	concurrency := getInt("WORKER_CONCURRENCY", 5)
	return &workerConfig{
		MaxRetry:    maxRetry,
		Timeout:     timeout,
		Concurrency: concurrency,
	}
}

var Worker = newWorkerConfig()
