package utils

func MaybePanic(err interface{}) {
	if err != nil {
		panic(err)
	}
}
