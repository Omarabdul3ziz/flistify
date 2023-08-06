package types

type Paths struct {
	ProjectDir string
	RootFS     string
	Boot       string
}

type Command struct {
	Name string
	Args []string
}
