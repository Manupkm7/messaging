# Seguridad

Este documento describe las medidas de seguridad implementadas en el backend.

## Cifrado

### Cifrado AES-256-GCM

Todos los mensajes (texto, audios e imágenes) se cifran automáticamente antes de almacenarse en la base de datos usando AES-256-GCM.

- **Algoritmo**: AES-256-GCM
- **Clave**: Derivada del `JWT_SECRET` usando SHA-256
- **Nonce**: Generado aleatoriamente para cada mensaje
- **Cifrado automático**: Los mensajes se cifran al crear y se descifran al recuperar

### Contraseñas

Las contraseñas se almacenan usando bcrypt con el costo por defecto (10).

## WSS (WebSocket Secure)

Para habilitar conexiones WebSocket seguras (WSS):

1. Generar certificados TLS:
```bash
# Para desarrollo (self-signed)
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

2. Configurar en `.env`:
```
ENABLE_TLS=true
TLS_CERT_FILE=cert.pem
TLS_KEY_FILE=key.pem
```

3. El servidor iniciará con HTTPS/WSS automáticamente

## Variables de Entorno Importantes

- `JWT_SECRET`: Debe ser una cadena larga y aleatoria. Se usa para:
  - Firmar tokens JWT
  - Derivar la clave AES para cifrado de mensajes
  
**⚠️ IMPORTANTE**: Cambia `JWT_SECRET` en producción y mantenlo seguro.

## Cookies de Sesión

Las cookies de sesión para clientes:
- `HttpOnly`: true (previene acceso desde JavaScript)
- `Secure`: Se establece automáticamente cuando TLS está habilitado
- `SameSite`: Strict (previene CSRF)

## Manejo de Errores

El sistema implementa un manejo de errores centralizado que:
- No expone detalles internos en producción
- Registra errores para debugging
- Devuelve respuestas JSON consistentes

## Recomendaciones de Producción

1. **Cambiar JWT_SECRET**: Usa una cadena aleatoria de al menos 32 caracteres
2. **Habilitar TLS**: Siempre usa certificados válidos en producción
3. **Revisar logs**: Monitorea los logs de errores regularmente
4. **Actualizar dependencias**: Mantén las dependencias actualizadas
5. **Rate limiting**: Considera implementar rate limiting para prevenir abusos

