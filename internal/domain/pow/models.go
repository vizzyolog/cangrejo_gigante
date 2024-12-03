package pow

type Challenge struct {
	Nonce      string
	Difficulty int
}

type Solution struct {
	Nonce    string
	Response string
}
