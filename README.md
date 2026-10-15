# CRUD de productos (Go + GORM + SQLite)

API REST local en un solo archivo (`main.go`) para crear, listar, consultar,
actualizar y eliminar productos. El precio se guarda en pesos enteros
(`80000` = 80.000 pesos). Los datos quedan en `productos.db`, que se crea en la
carpeta desde la que se ejecuta el programa.

| Método | Ruta | Acción y respuesta |
|---|---|---|
| POST | `/productos` | Crear: 201 + producto |
| GET | `/productos` | Listar: 200 + arreglo JSON |
| GET | `/productos/{id}` | Consultar: 200 + producto |
| PUT | `/productos/{id}` | Actualizar: 200 + producto |
| DELETE | `/productos/{id}` | Eliminar: 204 sin cuerpo |

## Requisitos

- Go 1.22 o superior (rutas con método y `{id}`).
- Compilador C y `CGO_ENABLED=1` (el driver SQLite usa CGO).

```bash
go version
gcc --version
go env CGO_ENABLED   # debe mostrar 1; si no: go env -w CGO_ENABLED=1
```

## Ejecutar

```bash
git clone https://github.com/filosofistichada/CRUDgo.git
cd CRUDgo
go mod download
go run .
```

Debe aparecer `API en http://localhost:8080`. Detener con Ctrl+C.

## Probar con curl

```bash
curl -i -X POST http://localhost:8080/productos \
  -H 'Content-Type: application/json' \
  -d '{"nombre":"Teclado","precio":80000}'

curl -i http://localhost:8080/productos
curl -i http://localhost:8080/productos/1

curl -i -X PUT http://localhost:8080/productos/1 \
  -H 'Content-Type: application/json' \
  -d '{"nombre":"Teclado mecanico","precio":95000}'

curl -i -X DELETE http://localhost:8080/productos/1
curl -i http://localhost:8080/productos/1   # 404
```

Errores: 400 si el JSON es inválido, falta `nombre`/`precio`, el precio es
negativo o el ID no es numérico; 404 si el producto no existe.
