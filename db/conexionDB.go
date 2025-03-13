package db

import (
	"database/sql"
	"os"
)

func ConexionDB() (*sql.DB, error) {
	Driver := "mysql"
	Usuario := "root"
	Password := os.Getenv("PASSWORD_MYSQL")
	Nombre := "aerolinea"

	conexion, err := sql.Open(Driver, Usuario+":"+Password+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}

	return conexion, err
}
