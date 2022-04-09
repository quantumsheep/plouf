# Plouf
Plouf is a simple, fast, and powerful API framework for Go. Its design is based on NestJS with features like dependency injection, modules, controllers, services, and more.

The framework is powered by the most popular Go libraries to let you use already existing middlewares, documentation and support.

Major libraries used by Plouf:
- [echo](https://github.com/labstack/echo)
- [gorm](https://github.com/go-gorm/gorm)
- [logrus](https://github.com/sirupsen/logrus)
- [validator](https://github.com/go-playground/validator)

Some of these libraries are abstracted by Plouf to make it easier to use.

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
