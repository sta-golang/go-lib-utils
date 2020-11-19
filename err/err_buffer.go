package err

type errBuffer struct {
	isAsync bool
	table map[int]Error
}