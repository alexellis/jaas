from golang:1.8
run mkdir -p /go/src/github.com/alexellis/gotask
copy app.go /go/src/github.com/alexellis/gotask/
workdir /go/src/github.com/alexellis/gotask
run go build

cmd ["./gotask"]