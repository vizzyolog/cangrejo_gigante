package network

type ConnectionServer interface {
	ListenAndServe() error
}
