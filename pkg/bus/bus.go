package bus

type Bus interface {
	GetHandler(command Dto) (Handler[any, Dto], error)
}
