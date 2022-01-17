package result

type Service interface {
	GetResults()
}

func NewService() Service {
	return nil
}
