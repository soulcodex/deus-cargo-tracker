package bus

type Dto interface {
	Type() string
}

type BlockingDto interface {
	Dto
	BlockingKey() string
}
