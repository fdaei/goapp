package outbox

type Config struct {
	IntervalInSeconds int `koanf:"interval_in_seconds"`
	BatchSize         int `koanf:"batch_size"`
	RetryThreshold    int `koanf:"retry_threshold"`
}
