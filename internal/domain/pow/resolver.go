package pow

type Resolver struct {
	Difficulty int
}

func NewPoWResolver(difficulty int) *Resolver {
	return &Resolver{Difficulty: difficulty}
}
