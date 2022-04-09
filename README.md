# Plouf
Plouf is a simple, fast, and powerful API framework for Go. Its design is based on NestJS with features like dependency injection, modules, controllers, services, and more.

The framework is powered by [echo](https://github.com/labstack/echo) which allow you to use already existing middlewares, documentation and support.

# Example file architecture
```
.
├── modules/
│   └── user/
│       ├── user_controller.go
│       ├── user_module.go
│       └── user_service.go
├── main_module.go
└── main.go
```
