package contracts

type IRequest interface {
	IsRequest() bool
}

type Request struct {
}

func NewRequest() IRequest {
	return &Request{}
}

func (r *Request) IsRequest() bool {
	return true
}
