package ioc

import "github.com/golobby/container/v3"

type Service interface {
	Start() error
}

func Start(c container.Container) error {
	return c.Call(func(services []Service) error {
		for _, service := range services {
			err := service.Start()
			if err != nil {
				return err
			}
		}
		return nil
	})
}
