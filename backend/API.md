# Documentación de la API

## Autenticación

Todas las rutas protegidas requieren un token JWT en el header:
```
Authorization: Bearer <token>
```

## Endpoints

### Autenticación

#### POST /api/auth/register
Registrar un nuevo usuario.

**Body:**
```json
{
  "username": "usuario123",
  "email": "usuario@example.com", // Opcional
  "password": "password123",
  "role": "client" // o "admin" (solo si un admin lo crea)
}
```

**Response:**
```json
{
  "token": "jwt_token_here",
  "user": {
    "id": "uuid",
    "username": "usuario123",
    "email": "usuario@example.com",
    "role": "client",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "session_token": "session_token_for_clients"
}
```

#### POST /api/auth/login
Iniciar sesión.

**Body:**
```json
{
  "username": "usuario123",
  "password": "password123",
  "session_token": "optional_session_token_for_clients"
}
```

**Response:** Igual que register.

#### POST /api/auth/logout
Cerrar sesión (requiere autenticación).

#### GET /api/auth/me
Obtener información del usuario actual (requiere autenticación).

### Usuarios

#### POST /api/users/admin
Crear un nuevo admin (solo admins pueden crear admins).

**Body:**
```json
{
  "username": "admin123",
  "email": "admin@example.com", // Opcional
  "password": "password123"
}
```

### Chats

#### GET /api/chats
Obtener todos los chats del usuario actual.

**Response:**
```json
[
  {
    "id": "uuid",
    "admin_id": "uuid",
    "client_id": "uuid",
    "assigned_to_id": "uuid",
    "is_admin_chat": false,
    "participant_name": "nombre_del_participante",
    "participant_id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "last_message_at": "2024-01-01T00:00:00Z"
  }
]
```

#### POST /api/chats
Crear un nuevo chat.

**Body:**
```json
{
  "client_id": "uuid", // Para chat con cliente
  // O
  "other_admin_id": "uuid" // Para chat entre admins
}
```

#### POST /api/chats/{id}/delegate
Delegar un chat a otro admin.

**Body:**
```json
{
  "to_admin_id": "uuid"
}
```

#### GET /api/chats/{id}/messages
Obtener mensajes de un chat.

**Query Parameters:**
- `limit`: Número de mensajes (default: 50)
- `offset`: Offset para paginación (default: 0)

**Response:**
```json
[
  {
    "id": "uuid",
    "chat_id": "uuid",
    "sender_id": "uuid",
    "content": "Mensaje de texto",
    "has_audio": false,
    "has_image": false,
    "mime_type": null,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

#### POST /api/chats/{id}/messages
Enviar un mensaje.

**Form Data (multipart/form-data):**
- `content`: Texto del mensaje (opcional)
- `audio`: Archivo de audio (opcional)
- `image`: Archivo de imagen (opcional)

**Response:**
```json
{
  "id": "uuid",
  "chat_id": "uuid",
  "sender_id": "uuid",
  "content": "Mensaje de texto",
  "is_audio": false,
  "is_image": false,
  "mime_type": "image/jpeg",
  "created_at": "2024-01-01T00:00:00Z"
}
```

#### GET /api/messages/{id}/file
Obtener archivo (audio o imagen) de un mensaje.

**Response:** Archivo binario con Content-Type apropiado.

### Grupos de Admins

#### POST /api/admin-groups
Crear un grupo de admins.

**Body:**
```json
{
  "name": "Grupo 1",
  "admin_ids": ["uuid1", "uuid2", "uuid3", "uuid4"]
}
```

#### GET /api/admin-groups
Obtener todos los grupos.

#### POST /api/admin-groups/{id}/replace
Reemplazar un admin en un grupo.

**Body:**
```json
{
  "old_admin_id": "uuid",
  "new_admin_id": "uuid"
}
```

### WebSocket

#### WS /ws
Conexión WebSocket para mensajería en tiempo real.

**Autenticación:**
- Header: `Authorization: Bearer <token>`
- O Query parameter: `?token=<token>`

**Mensajes del cliente:**
```json
{
  "action": "subscribe",
  "chat_id": "uuid"
}
```

```json
{
  "action": "unsubscribe",
  "chat_id": "uuid"
}
```

**Mensajes del servidor:**
```json
{
  "id": "uuid",
  "chat_id": "uuid",
  "sender_id": "uuid",
  "content": "Mensaje",
  "is_audio": false,
  "is_image": false,
  "created_at": "2024-01-01T00:00:00Z"
}
```

## Notas Importantes

1. **Cookies de Sesión**: Los clientes reciben una cookie `session_token` que se usa para evitar duplicados de usuarios.

2. **Limpieza Automática**: Todos los chats y mensajes se eliminan automáticamente después de 24 horas. Solo persisten las sesiones y los usuarios.

3. **Almacenamiento de Archivos**: Los audios e imágenes se almacenan directamente en PostgreSQL como BYTEA.

4. **Delegación de Chats**: Solo los admins pueden delegar chats a otros admins.

5. **Grupos de Admins**: Los admins pueden estar agrupados y pueden ser reemplazados por otros admins.

