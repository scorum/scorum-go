package sign

type Chain struct {
	ID string
}

var ScorumChain = &Chain{
	ID: "0000000000000000000000000000000000000000000000000000000000000000",
}

var TestChain = &Chain{
	ID: "f679096aa28f0019b5ceb1914b2c6a8488913276250a6d50c7037ad60e710257",
}
