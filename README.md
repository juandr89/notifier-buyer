# Notifier Buyer

API que alerta a los buyers que tienen programada
entrega el día de mañana sobre posibles retrasos en la entrega de sus paquetes.

## Tabla de Contenidos

- [Características](#características)
- [Tecnologías](#tecnologías)
- [Requisitos](#requisitos)
- [Instalación](#instalación)
- [Uso](#uso)
- [Pruebas](#pruebas)
- [Contribuciones](#contribuciones)
- [Licencia](#licencia)

## Características

- Arquitecturas Limpias
- Clean Code

## Tecnologías

- **Golang**: Lenguaje de programación. (1.22.5)
- **Redis**: Almacenamiento en memoria para gestión de datos.
- **Docker**: Contenerización de la aplicación para facilitar su despliegue.

## Requisitos

- [Docker](https://docs.docker.com/get-docker/) instalado en tu máquina.
- [Docker Compose](https://docs.docker.com/compose/install/) (opcional, pero recomendado para facilitar la gestión de contenedores).
**Golang**: [Instalación de Golang](https://golang.org/doc/install) (opcional)

## Instalación
### Opción 1
1. Clona el repositorio:
   ```bash
   git clone https://github.com/tu_usuario/nombre_de_la_aplicacion.git
   cd nombre_de_la_aplicacion

2. Variables de entorno
   ```bash
    Añadir los respectivos valores de las variables de entorno en el archivo config.yaml

2. Construye la imagen de docker:
    ```bash
    docker build -t nombre_de_la_aplicacion .

3. Inicia Redis usando Docker:
    ```bash
    docker run --name redis -d -p 6379:6379 redis

4. Inicia la aplicación:
    ```bash
    docker run --rm -p 8080:8080 nombre_de_la_aplicacion

### Opción 2
5. Si estás usando Docker Compose, puedes iniciar todos los servicios con:
    ```bash
    docker compose up --build

## Pruebas

1. Ejecutar el set de pruebas
    ```bash
    go test ./...  -coverprofile coverage.out

2. Validar porcentajes de cobertura
    ```bash
    go tool cover -func=coverage.out

