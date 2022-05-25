package main

import "fmt"

func main() {
	//######################################
	//	ES UN MENU PROVISIONAL
	//######################################

	fmt.Println("[ 1 ] - Obtener token")
	fmt.Println("[ 2 ] - Registrar nuevo usuario")
	fmt.Println("[ 3 ] - Salir")
	fmt.Println("Por favor, teclee 1, 2 o 3")
	fmt.Println()
	var eleccion int
	var eleccion2 string

	var continuar bool = true

	for continuar {
		fmt.Scanln(&eleccion)

		switch eleccion {
		case 1:
			//LoginUsuario("usuPrueba", "usuContra123")

			fmt.Println("¿Desea continuar?... (S/N)")
			fmt.Scanln(&eleccion2)
			if eleccion2 == "N" || eleccion2 == "n" {
				continuar = false
			} else if eleccion2 == "S" || eleccion2 == "s" {
				fmt.Println("Por favor, teclee 1, 2 o 3")
			} else {
				fmt.Println("Tecla inesperada, considero que quiere continuar, porfavor teclee 1, 2 o 3")
			}

		case 2:
			//RegistroUsuario("usuPrueba", "usuContra123")

			fmt.Println("¿Desea continuar?... (S/N)")
			fmt.Scanln(&eleccion2)

			if eleccion2 == "N" || eleccion2 == "n" {
				continuar = false
			} else if eleccion2 == "S" || eleccion2 == "s" {
				fmt.Println("Por favor, teclee 1, 2 o 3")
			} else {
				fmt.Println("Tecla inesperada, considero que quiere continuar, porfavor teclee 1, 2 o 3")
			}
		case 3:
			continuar = false
		default:
			fmt.Println("Tecla inesperada, por favor vuelve a intentarlo (1,2,3)")
		}
	}

	fmt.Println("¡Hasta Luego :)!")
}
