# Backend - Sistema de Mensajería Bidireccional

Backend desarrollado en Go con WebSockets y PostgreSQL para servicio de mensajería bidireccional.

## Características

- Autenticación JWT con roles (admin/client)
- Mensajería en tiempo real con WebSockets
- **Cifrado AES-256-GCM** para todos los mensajes y archivos
- **WSS (WebSocket Secure)** con soporte TLS/HTTPS
- Almacenamiento de archivos (audios/imágenes) en PostgreSQL
- Sistema de grupos de admins con delegación de chats
- Limpieza automática de chats cada 24 horas
- Cookies de sesión para clientes
- **Manejo de errores centralizado** y robusto
- **Bcrypt** para hash de contraseñas

## Requisitos

- Go 1.21 o superior
- PostgreSQL 15 o superior
- Docker y Docker Compose (opcional)

## Instalación

1. Clonar el repositorio
2. Copiar `.env.example` a `.env` y configurar las variables
3. Instalar dependencias:
```bash
go mod download
```

4. Iniciar PostgreSQL con Docker:
```bash
docker-compose up -d
```

5. Ejecutar migraciones:
```bash
go run cmd/migrate/main.go
```

6. Iniciar el servidor:
```bash
go run cmd/server/main.go
```

## Estructura del Proyecto

```
backend/
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── migrate/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── models/
│   ├── handlers/
│   ├── middleware/
│   ├── websocket/
│   ├── services/
│   └── utils/
├── go.mod
├── go.sum
├── docker-compose.yml
└── README.md
```

## API Endpoints

### Autenticación
- `POST /api/auth/register` - Registro de usuarios
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout

### Usuarios
- `GET /api/users/me` - Obtener usuario actual
- `POST /api/users/admin` - Crear admin (solo admins)

### Chats
- `GET /api/chats` - Listar chats
- `POST /api/chats` - Crear chat
- `POST /api/chats/:id/delegate` - Delegar chat a otro admin

### Mensajes
- `GET /api/chats/:id/messages` - Obtener mensajes de un chat
- `POST /api/chats/:id/messages` - Enviar mensaje

### WebSocket
- `WS /ws` - Conexión WebSocket para mensajería en tiempo real

## Variables de Entorno

Ver `.env.example` para todas las variables disponibles.

### Variables Importantes

- `JWT_SECRET`: Clave secreta para JWT y cifrado AES (cambiar en producción)
- `ENABLE_TLS`: Habilitar HTTPS/WSS (true/false)
- `TLS_CERT_FILE`: Ruta al certificado TLS
- `TLS_KEY_FILE`: Ruta a la clave privada TLS

## Seguridad

- Todos los mensajes se cifran automáticamente con AES-256-GCM
- Las contraseñas se hashean con bcrypt
- Soporte para WSS (WebSocket Secure) con TLS
- Cookies de sesión seguras (HttpOnly, Secure cuando TLS está habilitado)

Ver `SECURITY.md` para más detalles sobre seguridad.

