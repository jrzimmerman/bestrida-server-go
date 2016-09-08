## Bestrida Go Server

```
PORT=4001 go run main.go
```

```
docker build -t bestrida .
docker run -it --rm -p 4001:4001 --name bestrida-container bestrida
```