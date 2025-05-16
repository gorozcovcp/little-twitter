# Ualá Backend Challenge - Microblogging API

## 🌍 Descripción General
Este proyecto implementa una versión simplificada de una plataforma de microblogging similar a Twitter. Permite:

- Publicar mensajes cortos (tweets) de hasta 280 caracteres.
- Seguir a otros usuarios.
- Consultar un timeline que agrupa los tweets de los usuarios seguidos.

La aplicación está construida en **Golang**, utilizando el framework **Gin**, **MongoDB** como base de datos principal, y **Redis** como cache para optimizar lecturas.

---

## ⚙️ Tecnologías utilizadas
- Lenguaje: **Golang**
- Web Framework: **Gin**
- Base de Datos: **MongoDB (NoSQL)**
- Cache: **Redis**
- Contenedores: **Docker / Docker Compose**

---

## 📆 Consideraciones de Arquitectura

### Clean Architecture
- Separación por responsabilidades:
  - Handlers (expuestos por API REST)
  - Lógica de negocio embebida (a modularizar en servicios en futuras versiones)
  - Acceso a base de datos con cliente MongoDB

### Optimizado para Lecturas
- Se utilizan **índices compuestos** en `tweets(user_id, created)` para mejorar el rendimiento del timeline.
- El endpoint de timeline utiliza **paginación con limit y cursor temporal**.
- Se agrega **Redis como capa de cache** para evitar acceso frecuente a MongoDB y responder en milisegundos.
- Cache invalidado automáticamente cuando se publica un nuevo tweet o se sigue a un nuevo usuario.

### Escalabilidad
- El diseño permite escalar horizontalmente mediante:
  - **Sharding por user_id** en MongoDB
  - Separación futura de servicios de escritura y lectura (CQRS)
  - Integración con Redis para cachear timelines frecuentes

---

## 🔧 Endpoints

### POST /tweet
Publica un nuevo tweet
```json
{
  "user_id": "user123",
  "content": "Hola mundo!"
}
```

### POST /follow
Sigue a otro usuario
```json
{
  "user_id": "user123",
  "follow_id": "user456"
}
```

### GET /timeline/:userID
Obtiene el timeline del usuario
```bash
GET /timeline/user123?limit=20&since=2025-05-14T12:00:00Z
```

---

## 🛠️ Instalación y Ejecución

### Requisitos
- Docker
- Docker Compose

### Pasos
```bash
git clone <repo_url>
cd <repo>
docker-compose build
docker-compose up -d
```

La API estará disponible en `http://localhost:8080`

MongoDB: `localhost:27017`
Redis: `localhost:6379`

---
