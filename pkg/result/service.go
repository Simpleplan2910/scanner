package result

type Service interface {
	GetResults()
}

type service struct {
}

func NewService() Service {
	return nil
}
