# General Go Vue Admin

## Project Description

General Go Vue Admin is an enterprise-grade backend management system designed to provide a robust, scalable, and highly performant architecture. It implements a strict separation of frontend and backend environments. The backend is built on the Go language utilizing the Gin web framework and GORM for database interactions. The frontend leverages Vue 3, Composition API, TypeScript, Vite, and Tailwind CSS to deliver an optimal and highly responsive user interface featuring modern Glassmorphism aesthetics.

This system is engineered for stability and security, incorporating high-concurrency capabilities, dynamic role-based access control (RBAC), and persistent state management, making it an ideal foundation for complex business applications.

## Core Features

- **Dynamic RBAC Authorization**: Granular control over permissions and data access based on hierarchical Roles, Departments, and Posts.
- **Frontend-Backend Decoupling**: API-first design using standard RESTful interfaces with JWT-based stateless authentication.
- **Micro-Frontend Architecture**: Multi-tab (TagsView) keep-alive navigation caching ensures operational state is preserved across route transitions.
- **System Administration Modules**: Comprehensive modules for User, Role, Menu, Department, and Post management.
- **Data Dictionary System**: Flexible configuration layer for mapping business variables internally without hardcoded values.
- **Operation Audit Trails**: Detailed and automated logging of system events, logins, and API mutations.
- **High Concurrency Stability**: Built-in concurrency managers, rate limiters, timeout wrappers, and circuit breakers mapped at the middleware routing layer.

## Technology Stack

### Backend
- **Go**: Version 1.22+
- **Gin**: High-performance HTTP web framework
- **GORM**: Object-Relational Mapping framework
- **MySQL**: Relational database storage
- **Redis**: In-memory data store for caching and queues
- **JWT**: Token-based authentication

### Frontend
- **Vue 3**: Reactive UI framework focusing on the Composition API
- **Vite**: Next-generation frontend tooling and bundler
- **TypeScript**: Static typing for structural soundness
- **Pinia**: Intuitive and type-safe state management
- **Tailwind CSS**: Utility-first CSS framework for rapid UI styling
- **Vue Router**: Official routing with dynamic permission tree integration

## Development Environment Setup

### Prerequisites
- Go 1.22 or higher
- Node.js 18 or higher
- MySQL 8.0 or higher
- Redis 6.0 or higher
- Yarn or npm package manager

### Backend Initialization
1. Navigate to the server directory:
   `cd server`
2. Duplicate the configuration template and establish database credentials:
   `cp config.yaml.example config.yaml`
3. Download Go modules:
   `go mod tidy`
4. Start the backend service using Air (for live-reloading) or Go:
   `air` or `go run main.go`

### Frontend Initialization
1. Navigate to the frontend directory:
   `cd web`
2. Install package dependencies:
   `yarn install` or `npm install`
3. Start the Vite development server:
   `yarn dev` or `npm run dev`

## Deployment

For production deployment, ensure the frontend is compiled into static assets and the Go application is built into a standalone binary.

1. Build the frontend:
   `cd web`
   `yarn build`
   This will generate a `dist` directory. Serve these files via Nginx or equivalent web servers.
2. Build the backend:
   `cd server`
   `GOOS=linux GOARCH=amd64 go build -o server main.go`

## Directory Structure Overview

- `/server`
  - `/api`: Contains Controller, Service, DAO, and Entity definitions.
  - `/common`: Configuration loading, constants, and structured response formatters.
  - `/middleware`: Authentication, CORS, logging, and concurrency middlewares.
  - `/pkg`: Extended utilities and plugins (Database, Redis, JWT encapsulation).
  - `/router`: Centralized route multiplexing and handlers.
- `/web`
  - `/src/api`: Unified API endpoint wrappers mapping to the backend.
  - `/src/components`: Reusable UI components including foundational layouts and widgets.
  - `/src/router`: Frontend routing arrays with navigation guards.
  - `/src/stores`: Pinia states definitions for user profile, authentication, and Multi-Tab tags.
  - `/src/views`: Domain-specific pages and modular templates.

## License

This project operates under the constraints defined within the workspace and is intended for internal commercial deployment architectures.
