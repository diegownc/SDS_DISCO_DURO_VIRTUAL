package db

import (
	"database/sql"
	"fmt"

	"github.com/diegownc/SDS_DISCO_DURO_VIRTUAL/argon2sds"
)

//Conexión con base de datos
const (
	host     = "localhost"
	port     = 5432
	user     = "normaluser"
	password = "contra123"
	dbname   = "sds"
)

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//#######################################################################################################
//	FUNCIONES
//#######################################################################################################

func comprobarUsuario(username string, password string, hash_username string, hash_pass string) bool {
	var res bool = true

	match, err := argon2sds.CompareHash(username, hash_username)
	if !match || err != nil {
		res = false
	}

	match, err = argon2sds.CompareHash(password, hash_pass)
	if !match || err != nil {
		res = false
	}

	return res
}

func LoginUsuario(nombreUsuario string, contraUsuario string) bool {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	checkError(err)
	defer db.Close()

	//Comprobamos si ya existe el usuario dentro de la BD:
	var existeUsuario bool = false
	var usernameDB string
	var passwordDB string

	rows, err := db.Query("Select username, password from users")
	checkError(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&usernameDB, &passwordDB)
		checkError(err)

		if comprobarUsuario(nombreUsuario, contraUsuario, usernameDB, passwordDB) {
			existeUsuario = true
		}
	}
	err = rows.Err()
	checkError(err)

	if existeUsuario {
		//ENVIAR TOKEN AL USUARIO
		fmt.Println("Credenciales correctas")
	} else {
		fmt.Println("Crendenciales incorrectas")
	}

	return existeUsuario
}

func RegistroUsuario(nombreUsuario string, contraUsuario string) bool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	checkError(err)
	defer db.Close()

	//Comprobamos si ya existe el usuario dentro de la BD:
	var existeUsuario bool = false
	var usernameDB string
	var passwordDB string

	rows, err := db.Query("Select username, password from users")
	checkError(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&usernameDB, &passwordDB)
		checkError(err)

		if comprobarUsuario(nombreUsuario, contraUsuario, usernameDB, passwordDB) {
			existeUsuario = true
		}
	}
	err = rows.Err()
	checkError(err)

	//¿El usuario existe?
	if !existeUsuario {

		hash_username, err := argon2sds.GenerateHash(nombreUsuario)
		checkError(err)

		hash_pass, err := argon2sds.GenerateHash(contraUsuario)
		checkError(err)

		sqlStatement := `INSERT INTO users (username, password) VALUES ($1, $2)`
		_, err = db.Exec(sqlStatement, hash_username, hash_pass)

		checkError(err)
		fmt.Println("\nSe ha añadido un nuevo usuario a la BD!")
	} else {
		fmt.Println("El usuario ya existe")
	}

	return !existeUsuario
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
