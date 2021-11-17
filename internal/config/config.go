package config

type Config struct {
	GinMode        string
	ServerPort     string
	Login          string
	Password       string
	AutoStartTasks string
	PathToStorage  string
	PathToWWW      string
	DelayMin       float64
	DelayMax       float64
}
