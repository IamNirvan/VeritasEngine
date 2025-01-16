package models

type Log struct {
	Level   string `mapstructure:"level" json:"level"`
	Methods bool   `mapstructure:"methods" json:"methods"`
}

type WebServer struct {
	Host    string `mapstructure:"host" json:"host"`
	Port    uint16 `mapstructure:"port" json:"port"`
	Timeout uint16 `mapstructure:"timeout" json:"timeout"`
}

type Database struct {
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Dbname   string `mapstructure:"dbname" json:"dbname"`
	Sslmode  string `mapstructure:"sslmode" json:"sslmode"`
}