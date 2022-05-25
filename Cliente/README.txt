//Pasos para probar el cliente 

1 - Ejecutar el main.go de la carpeta de fuera para que se ejecute el servidor
go run main.go

2 - Descargar las dependencias de la carpeta cliente
cd Cliente
go mod tidy

3 - Ejecutar el cliente.go en otro terminal (este servidor solo resuelve peticiones del html)
go run cliente.go

4 - Abrir "index.html" en el navegador
Click derecho y abrir en navegador..

5 - Abrir consola del inspeccionador para ver los resultados
Ctrl + Shift + I 
