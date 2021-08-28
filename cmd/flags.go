package main

import "flag"

func init() {
	flag.StringVar(&inputFile, "input-file", "file.csv", "name of CSV-file")
	flag.StringVar(&dbIP, "db-ip", "127.0.0.1", "database IP-address")
	flag.StringVar(&dbPort, "db-port", "5432", "database port")
	flag.StringVar(&dbName, "db-name", "csv", "database name")
	flag.StringVar(&dbLogin, "db-login", "postgres", "database user name")
	flag.StringVar(&dbPasswd, "db-passwd", "postgres", "database user password")
	flag.Parse()
}
