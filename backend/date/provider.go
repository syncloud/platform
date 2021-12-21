package date

import "time"

type RealProvider struct {
}

type Provider interface {
	Now() time.Time
}

func New() *RealProvider {
	return &RealProvider{}
}

func (d *RealProvider) Now() time.Time {
	return time.Now()
}
