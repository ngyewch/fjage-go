package fjage

type PerformativeError struct {
	performative Performative
}

func NewPerformativeError(perf Performative) *PerformativeError {
	return &PerformativeError{
		performative: perf,
	}
}

func (err PerformativeError) Performative() Performative {
	return err.performative
}

func (err PerformativeError) Error() string {
	return string(err.performative)
}
