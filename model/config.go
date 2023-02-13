package model

type Config struct {
	Redis struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		Port     string `json:"port"`
		Database int    `json:"db"`
	} `json:"redis"`
	Port string `json:"port"`
}
