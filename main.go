package main

import (
	"CRUD/db"
	"CRUD/models"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

var plantillaIndex, _ = template.ParseFiles("plantillas/index.html")
var crear, _ = template.ParseFiles("plantillas/crear.html")
var editar, _ = template.ParseFiles("plantillas/editar.html")

func main() {
	// manejador de endpoints
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	log.Print("Listening on port 8080...")

	// creacion de un servidor local.
	_ = http.ListenAndServe(":8080", nil)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		dni := r.FormValue("dni")
		nombre := r.FormValue("nombre")
		apellido := r.FormValue("apellido")
		numeroContacto := r.FormValue("numero")
		correo := r.FormValue("email")

		conexion, _ := db.ConexionDB()

		idCliente, _ := strconv.Atoi(id)

		modificarRegistro, err := conexion.Prepare("UPDATE cliente SET cli_DNI=?, cli_nombre=?, cli_apellido=?, cli_numContacto=?, cli_correo=? WHERE cli_id=?")

		if err != nil {
			panic(err.Error())
		}

		_, _ = modificarRegistro.Exec(dni, nombre, apellido, numeroContacto, correo, idCliente)

		log.Println("Registro modificado con exito.")

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conexion, _ := db.ConexionDB()
	registros, e := conexion.Query("SELECT * FROM cliente")

	if e != nil {
		panic(e.Error())
	}

	cliente := models.Cliente{}
	var arreglo []models.Cliente

	for registros.Next() {

		var id int
		var dni, name, lastName, numberContac, email string
		err := registros.Scan(&id, &dni, &name, &lastName, &numberContac, &email)

		if err != nil {
			log.Fatal(err)
		}

		cliente.Id = id
		cliente.Dni = dni
		cliente.Nombre = name
		cliente.Apellido = lastName
		cliente.NumeroContacto = numberContac
		cliente.Correo = email

		arreglo = append(arreglo, cliente)
	}

	_ = plantillaIndex.Execute(w, arreglo)
}

func Crear(w http.ResponseWriter, r *http.Request) {
	_ = crear.Execute(w, nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		dni := r.FormValue("dni")
		nombre := r.FormValue("nombre")
		apellido := r.FormValue("apellido")
		numeroContacto := r.FormValue("numero")
		correo := r.FormValue("email")

		conexion, _ := db.ConexionDB()

		idCliente, _ := strconv.Atoi(id)

		newRegistro, err := conexion.Prepare("INSERT INTO cliente(cli_id, cli_DNI, cli_nombre, cli_apellido, cli_numContacto, cli_correo) VALUES (?,?,?,?,?,?)")

		if err != nil {
			panic(err.Error())
		}
		_, _ = newRegistro.Exec(idCliente, dni, nombre, apellido, numeroContacto, correo)
		log.Println("Registro guardado con exito.")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Borrar(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	conexion, _ := db.ConexionDB()

	borrar, err := conexion.Prepare("DELETE FROM cliente WHERE cli_id = ?")

	if err != nil {
		panic(err.Error())
	}
	_, _ = borrar.Exec(id)
	log.Println("Registro Eliminado con exito.")

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func Editar(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	conexion, _ := db.ConexionDB()

	registro, _ := conexion.Query("SELECT * FROM cliente WHERE cli_id = ? LIMIT 1", id)

	cliente := models.Cliente{}

	for registro.Next() {

		var id int
		var dni, name, lastName, numberContac, email string
		err := registro.Scan(&id, &dni, &name, &lastName, &numberContac, &email)

		if err != nil {
			log.Fatal(err)
		}

		cliente.Id = id
		cliente.Dni = dni
		cliente.Nombre = name
		cliente.Apellido = lastName
		cliente.NumeroContacto = numberContac
		cliente.Correo = email
	}

	_ = editar.Execute(w, cliente)

}
