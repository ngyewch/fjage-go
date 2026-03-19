package gateway

type MessageWrapper[T any] struct {
	Clazz string
	Data  T
}
