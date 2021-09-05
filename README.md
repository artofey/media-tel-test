# media-tel-test
Базу данных можно запустить через docker-compose.
После чего потыкать варианты запуска из Makefile.

## Можно уменьшить размер бинарного файла

1. Скомпилировать без отладочной информации
```
go build -ldflags="-s -w" -o main ./cmd
```
2. Сжать при помощи [UPX](https://upx.github.io/)
```
upx --brute -o main.upx main
```
Результат: ***5,7M -> 1,3M***
