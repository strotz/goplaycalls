package gpc

type Player struct{}
type Report struct{}

func (r *Report) Passed() bool {
	return true
}

func (r *Report) LastError() error {
	return nil
}

func (p *Player) Play() *Report {
	return &Report{}
}

// ParseFile creates a new Player for http request file.
func ParseFile(filePath string) (*Player, error) {
	return &Player{}, nil
}
