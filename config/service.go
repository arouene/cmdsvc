package config

type Service struct {
	Name      string
	Route     string
	Command   string
	Exclusive bool
	WorkDir   string
	Environ   []string
	Groups    []string
}

type ServiceList []Service

var Services ServiceList

func (s ServiceList) Get(name string) (Service, bool) {
	for _, svc := range s {
		if svc.Name == name {
			return svc, true
		}
	}
	return Service{}, false
}
