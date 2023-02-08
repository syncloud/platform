package ioc

type Service interface {
	Start() error
}

func Start() error {
	Call(func(services []Service) error {
		for _, service := range services {
			err := service.Start()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}
