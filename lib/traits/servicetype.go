package traits

type Computeable struct{}

func (_ Computeable) ServiceType() string {
	return "compute"
}

type Filesable struct{}

func (_ Filesable) ServiceType() string {
	return "files"
}

type Networkingable struct{}

func (_ Networkingable) ServiceType() string {
	return "networking"
}
