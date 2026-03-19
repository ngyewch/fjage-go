package fjage

type PerformativeError struct {
	perf string
}

func NewPerformativeError(perf string) *PerformativeError {
	return &PerformativeError{
		perf: perf,
	}
}

func (err PerformativeError) Perf() string {
	return err.perf
}

func (err PerformativeError) Error() string {
	return err.perf
}
