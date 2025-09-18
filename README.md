# Backend Padel Go

Backend API para gestión de canchas de padel desarrollado en Go con Gin, GORM y MySQL.

## Características

- **Autenticación JWT** con refresh tokens
- **Gestión de canchas** con horarios de atención y horarios especiales
- **Sistema de reservas** con verificación de disponibilidad
- **Reseñas y calificaciones** de canchas
- **Integración con MercadoPago** para pagos
- **Búsqueda y filtros** avanzados de canchas
- **API REST** con documentación Swagger
- **Middleware** para CORS, autenticación y validación
- **Base de datos MySQL** con migraciones automáticas

## Tecnologías

- **Go 1.21+**
- **Gin** - Framework web
- **GORM** - ORM para base de datos
- **MySQL** - Base de datos
- **JWT** - Autenticación
- **MercadoPago SDK** - Pagos
- **Swagger** - Documentación API

## Instalación

### Prerrequisitos

- Go 1.21 o superior
- MySQL 8.0 o superior
- Git

### Configuración

1. **Clonar el repositorio**
```bash
git clone <repository-url>
cd backend-padel-go
```

2. **Instalar dependencias**
```bash
go mod tidy
```

3. **Configurar variables de entorno**
```bash
cp .env.example .env
```

Editar el archivo `.env` con tus configuraciones:
```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=padel_db

# JWT
JWT_SECRET=your-super-secret-jwt-key

# MercadoPago
MERCADOPAGO_ACCESS_TOKEN=your_mercadopago_access_token
```

4. **Crear base de datos**
```sql
CREATE DATABASE padel_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

5. **Ejecutar la aplicación**
```bash
go run main.go
```

La API estará disponible en `http://localhost:8080`

## Documentación API

Una vez que la aplicación esté ejecutándose, puedes acceder a la documentación Swagger en:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`

## Estructura del Proyecto

```
backend-padel-go/
├── internal/
│   ├── config/          # Configuración de la aplicación
│   ├── database/        # Configuración de base de datos
│   ├── handlers/        # Controladores HTTP
│   ├── middleware/      # Middleware personalizado
│   ├── models/          # Modelos de datos
│   └── services/        # Lógica de negocio
├── main.go             # Punto de entrada
├── go.mod              # Dependencias Go
├── .env.example        # Variables de entorno de ejemplo
└── README.md           # Este archivo
```

## Endpoints Principales

### Autenticación
- `POST /api/v1/auth/register` - Registro de usuario
- `POST /api/v1/auth/login` - Inicio de sesión
- `POST /api/v1/auth/refresh` - Renovar token

### Canchas
- `GET /api/v1/courts` - Listar todas las canchas
- `GET /api/v1/courts/:id` - Obtener cancha por ID
- `GET /api/v1/courts/nearby` - Canchas cercanas
- `GET /api/v1/courts/search` - Buscar canchas con filtros
- `GET /api/v1/courts/:id/availability` - Disponibilidad de cancha

### Gestión de Canchas (Propietarios)
- `POST /api/v1/owner/courts` - Crear cancha
- `GET /api/v1/owner/courts` - Mis canchas
- `PUT /api/v1/owner/courts/:id` - Actualizar cancha
- `DELETE /api/v1/owner/courts/:id` - Eliminar cancha
- `GET /api/v1/owner/courts/:id/statistics` - Estadísticas de cancha

### Reservas
- `POST /api/v1/bookings` - Crear reserva
- `GET /api/v1/bookings` - Mis reservas
- `GET /api/v1/bookings/:id` - Obtener reserva por ID
- `PUT /api/v1/bookings/:id/cancel` - Cancelar reserva

### Reseñas
- `POST /api/v1/reviews/courts/:id` - Crear reseña
- `PUT /api/v1/reviews/:id` - Actualizar reseña
- `DELETE /api/v1/reviews/:id` - Eliminar reseña
- `GET /api/v1/courts/:id/reviews` - Reseñas de cancha

### Pagos
- `POST /api/v1/payments/preference` - Crear preferencia de pago
- `GET /api/v1/payments/:id/status` - Estado del pago
- `POST /api/v1/payments/webhook` - Webhook de MercadoPago

## Modelos de Datos

### User
- Información básica del usuario
- Roles: user, owner, admin
- Autenticación JWT

### Court
- Información de la cancha
- Horarios de atención
- Horarios especiales
- Estadísticas (rating, reseñas)

### Booking
- Reserva de cancha
- Estado: pending, confirmed, cancelled, completed
- Cálculo automático de precio

### Review
- Reseñas de canchas
- Rating de 1 a 5 estrellas
- Actualización automática de estadísticas

### Payment
- Integración con MercadoPago
- Estados: pending, approved, rejected, cancelled
- Webhook para actualizaciones automáticas

## Desarrollo

### Ejecutar en modo desarrollo
```bash
go run main.go
```

### Ejecutar tests
```bash
go test ./...
```

### Generar documentación Swagger
```bash
swag init
```

### Compilar para producción
```bash
go build -o padel-backend main.go
```

## Variables de Entorno

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `PORT` | Puerto del servidor | 8080 |
| `GIN_MODE` | Modo de Gin (debug/release) | debug |
| `DB_HOST` | Host de la base de datos | localhost |
| `DB_PORT` | Puerto de la base de datos | 3306 |
| `DB_USER` | Usuario de la base de datos | root |
| `DB_PASSWORD` | Contraseña de la base de datos | - |
| `DB_NAME` | Nombre de la base de datos | padel_db |
| `JWT_SECRET` | Clave secreta para JWT | - |
| `MERCADOPAGO_ACCESS_TOKEN` | Token de acceso de MercadoPago | - |

## Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.