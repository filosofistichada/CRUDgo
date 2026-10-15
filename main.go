package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Producto struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre"`
	Precio int64  `json:"precio"` // Pesos enteros, por ejemplo: 25000.
}

type Entrada struct {
	Nombre string `json:"nombre"`
	Precio *int64 `json:"precio"` // Puntero para detectar un campo omitido.
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("productos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err = db.AutoMigrate(&Producto{}); err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /productos", crear)
	mux.HandleFunc("GET /productos", listar)
	mux.HandleFunc("GET /productos/{id}", obtener)
	mux.HandleFunc("PUT /productos/{id}", actualizar)
	mux.HandleFunc("DELETE /productos/{id}", eliminar)
	log.Println("API en http://localhost:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
}

func responder(w http.ResponseWriter, estado int, datos any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(estado)
	if err := json.NewEncoder(w).Encode(datos); err != nil {
		log.Println(err)
	}
}

func fallo(w http.ResponseWriter, estado int, mensaje string) {
	responder(w, estado, map[string]string{"error": mensaje})
}

func leer(w http.ResponseWriter, r *http.Request) (Entrada, bool) {
	var e Entrada
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		fallo(w, 400, "JSON invalido")
		return e, false
	}
	e.Nombre = strings.TrimSpace(e.Nombre)
	if e.Nombre == "" || e.Precio == nil || *e.Precio < 0 {
		fallo(w, 400, "Envia nombre y precio entero mayor o igual a 0")
		return e, false
	}
	return e, true
}

func buscar(w http.ResponseWriter, r *http.Request) (Producto, bool) {
	var p Producto
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		fallo(w, 400, "ID invalido")
		return p, false
	}
	err = db.First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fallo(w, 404, "Producto no encontrado")
		return p, false
	}
	if err != nil {
		fallo(w, 500, "Error de base de datos")
		return p, false
	}
	return p, true
}

func crear(w http.ResponseWriter, r *http.Request) {
	e, ok := leer(w, r)
	if !ok {
		return
	}
	p := Producto{Nombre: e.Nombre, Precio: *e.Precio}
	if err := db.Create(&p).Error; err != nil {
		fallo(w, 500, "No se pudo crear")
		return
	}
	responder(w, 201, p)
}

func listar(w http.ResponseWriter, r *http.Request) {
	productos := []Producto{}
	if err := db.Order("id").Find(&productos).Error; err != nil {
		fallo(w, 500, "No se pudo listar")
		return
	}
	responder(w, 200, productos)
}

func obtener(w http.ResponseWriter, r *http.Request) {
	p, ok := buscar(w, r)
	if ok {
		responder(w, 200, p)
	}
}

func actualizar(w http.ResponseWriter, r *http.Request) {
	p, ok := buscar(w, r)
	if !ok {
		return
	}
	e, ok := leer(w, r)
	if !ok {
		return
	}
	p.Nombre, p.Precio = e.Nombre, *e.Precio
	if err := db.Model(&p).Select("Nombre", "Precio").Updates(p).Error; err != nil {
		fallo(w, 500, "No se pudo actualizar")
		return
	}
	responder(w, 200, p)
}

func eliminar(w http.ResponseWriter, r *http.Request) {
	p, ok := buscar(w, r)
	if !ok {
		return
	}
	if err := db.Delete(&p).Error; err != nil {
		fallo(w, 500, "No se pudo eliminar")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
