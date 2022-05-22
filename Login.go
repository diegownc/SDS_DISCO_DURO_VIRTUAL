package main

import (
	"database/sql"
	"fmt"
)

//Conexión con base de datos
const (
	host     = "localhost"
	port     = 5432
	user     = "normaluser"
	password = "contra123"
	dbname   = "sds"
)

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func main() {
	//Descomentar para logear a un usuario
	//LoginUsuario("usuPrueba123", "usuContra123")

	//Descomentar para registrar un nuevo usuario
	//RegistroUsuario("usuPrueba123", "usuContra123")
}

//#######################################################################################################
//	FUNCIONES
//#######################################################################################################

func comprobarUsuario(username string, password string, hash_username string, hash_pass string) bool {
	var res bool = true

	match, err := CompareHash(username, hash_username)
	if !match || err != nil {
		res = false
	}

	match, err = CompareHash(password, hash_pass)
	if !match || err != nil {
		res = false
	}

	return res
}

func LoginUsuario(nombreUsuario string, contraUsuario string) {
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
	} else {
		fmt.Println("Crendenciales incorrectas")
	}
}

func RegistroUsuario(nombreUsuario string, contraUsuario string) {
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
		config := &PasswordConfig{
			time:    1,
			memory:  64 * 1024,
			threads: 4,
			keyLen:  32,
		}

		hash_username, err := GenerateHash(config, nombreUsuario)
		checkError(err)

		hash_pass, err := GenerateHash(config, contraUsuario)
		checkError(err)

		sqlStatement := `INSERT INTO users (username, password, idfolder) VALUES ($1, $2, $3)`
		_, err = db.Exec(sqlStatement, hash_username, hash_pass, "testing5")

		checkError(err)
		fmt.Println("\nSe ha añadido un nuevo usuario a la BD!")
	} else {
		fmt.Println("El usuario ya existe")
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
