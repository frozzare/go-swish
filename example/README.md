# swish demo

Start a [ngrok](https://ngrok.com/) tunnel:

```
ngrok http 5000
```

Change ngrok url in the application and start the demo server from the root directory:

```
go run examples/main.go
```

Visit [http://localhost:5000/create](http://localhost:5000/create) to create a payment request, copy the payment request id and visit [http://localhost:5000/status?id=](http://localhost:5000/status?id=).