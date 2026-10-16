package service_open_session

type Request struct{}

func (x Request) Validate() error {
	return nil
}

func NewRequest() Request {
	return Request{}
}
