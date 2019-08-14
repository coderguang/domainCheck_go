package domainCheckDef

type DbConfig struct {
	Dburl         string `json:"Dburl"`
	Dbport        string `json:"Dbport"`
	Dbuser        string `json:"Dbuser"`
	Dbpwd         string `json:"Dbpwd"`
	Dbname        string `json:"Dbname"`
	DoroutinesNum int    `json:"DoroutinesNum"`
	EmailFrom     string `json:"EmailFrom"`
	EmailPwd      string `json:"EmailPwd"`
	EmailTo       string `json:"EmailTo"`
	Smtp          string `json:"Smtp"`
	SmtpPort      string `json:"SmtpPort"`
}
