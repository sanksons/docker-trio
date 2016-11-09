package sqldb

import ()

type Config struct {
	DriverName string
	Username   string
	Password   string
	Host       string
	Port       string
	Dbname     string
	Timezone   string
	MaxOpenCon int
	MaxIdleCon int
}
