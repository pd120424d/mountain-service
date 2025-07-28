# Mountain Management System

A web application built with Angular for managing mountain service operations, including employee management, shift scheduling, and emergency response coordination.

## Features

- Multi-language support (English, Serbian Latin, Serbian Cyrillic, Russian)
- Employee management system
- Authentication and authorization
- Responsive design with Bootstrap and Material UI
- Image carousel on home page
- RESTful API integration

## Prerequisites

- Node.js (v18 or later)
- npm (Node Package Manager)
- Angular CLI (v19.0.6)

## Installation

1. Clone the repository
2. Install dependencies:
```bash
npm install
```

## Build Commands

### Development (staging features visible)
```bash
ng serve
# or
ng serve --configuration=development
```

### Staging (staging features visible)
```bash
ng build --configuration=staging
ng serve --configuration=staging
```

### Production (staging features hidden)
```bash
ng build --configuration=production
ng serve --configuration=production
```

## Production Build

1. Build the project:
```bash
ng build --configuration production
```

2. The build artifacts will be stored in the `dist/ui/` directory.

## Docker Deployment

Build and run the application using Docker:

```bash
docker build -t mountain-service-ui .
docker run -p 80:80 mountain-service-ui
```

## Environment Configuration

- Development: `src/environments/environment.ts`
- Production: `src/environments/environment.prod.ts`

Configure the API URL and other environment-specific settings in these files.

## Testing

- Run unit tests:
```bash
ng test
```

- Run end-to-end tests:
```bash
ng e2e
```

## Project Structure

```
ui/
├── src/
│   ├── app/
│   │   ├── components/
│   │   ├── services/
│   │   └── shared/
│   ├── assets/
│   │   └── i18n/          # Translation files
│   └── environments/      # Environment configurations
├── public/               # Public assets
└── dist/                # Production build output
```

## Internationalization

The application supports multiple languages. Translation files are located in `src/assets/i18n/`.
Available languages:
- English (en)
- Serbian Latin (sr-lat)
- Serbian Cyrillic (sr-cyr)
- Russian (ru)

## Dependencies

- Angular v19.0.5
- Angular Material v18.2.10
- Bootstrap v5.3.3
- Font Awesome v6.6.0
- NgX-Translate
- RxJS

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---


