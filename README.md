# StockViz

## Descripción
StockViz es una plataforma para visualizar datos de acciones del mercado bursátil y obtener recomendaciones de inversión basadas en análisis avanzados. La aplicación permite a los usuarios acceder a información actualizada sobre stocks, visualizar su rendimiento de manera intuitiva y analizar tendencias del mercado.

## Características principales
- Visualización de datos actualizados de stocks del mercado
- Análisis inteligente para identificar patrones y tendencias
- Recomendaciones de inversión basadas en análisis avanzados
- Sincronización automática con fuentes de datos externas
- Interfaz de usuario moderna e intuitiva

## Estructura del proyecto
El proyecto está organizado en una arquitectura cliente-servidor:

### Backend (Go)
- **api/**: Implementa los endpoints HTTP y la lógica de manejo de solicitudes
- **analyzer/**: Contiene la lógica para el análisis de datos de stocks
- **config/**: Gestiona la configuración de la aplicación
- **database/**: Proporciona funcionalidades para interactuar con la base de datos
- **models/**: Define las estructuras de datos utilizadas en la aplicación
- **server/**: Implementa el servidor HTTP
- **service/**: Contiene la lógica de negocio

### Frontend (Vue.js)
- **src/**: Código fuente de la aplicación Vue.js
  - **components/**: Componentes reutilizables
  - **views/**: Páginas principales
  - **store/**: Estado global (Pinia)
  - **router/**: Configuración de rutas

## Requisitos previos
- Go 1.16 o superior
- Node.js 14 o superior
- PostgreSQL
- API de datos de stocks (configurada a través de variables de entorno)

## Instalación

### Backend
1. Navega al directorio backend:
   ```bash
   cd backend
   ```

2. Instala las dependencias:
   ```bash
   go mod download
   ```

3. Configura las variables de entorno (crea un archivo .env en el directorio backend):
   ```
   DATABASE_URL=postgresql://usuario:contraseña@localhost:5432/stockviz
   STOCK_API_TOKEN=tu_token_de_api
   STOCK_API_BASE_URL=https://api.ejemplo.com
   SYNC_INTERVAL=60m
   ```

4. Ejecuta la aplicación:
   ```bash
   go run main.go
   ```

### Frontend
1. Navega al directorio frontend:
   ```bash
   cd frontend
   ```

2. Instala las dependencias:
   ```bash
   npm install
   ```
   
3. Para producción, compila los archivos:
   ```bash
   npm run build
   ```

## Uso
Una vez que ambos servicios (backend y frontend) estén en ejecución:

1. Accede a la aplicación a través de tu navegador en `http://localhost:8080`
2. Navega a la sección "Stocks" para ver información actualizada de acciones
3. Consulta la sección "Recomendaciones" para ver análisis y recomendaciones de inversión

## API REST
El backend expone los siguientes endpoints principales:

- `GET /api/stocks`: Obtiene datos de stocks con paginación
- `GET /api/stocks/all`: Obtiene todos los stocks con paginación opcional
- `GET /api/recommendations`: Obtiene recomendaciones de inversión
