### Usage example:
- Below an example is provided showing how one can use this small http server wrapper.
- Keep in mind that this DOES NOT implement any RFC and is made purely as a showcase demonstrating the implementation of a http LIKE server using the Golang net library.
```go
package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	shttp "github.com/WolvenSpirit/go-http-server"
)

func router(r shttp.Request, wr shttp.Response) {
	switch r.Path {
	case "/":
		page := `<html><body>
		<h3>Welcome to Wolven's website</h3>
		</body></html>`
		wr.Write([]byte(page))
		break
	case "/blog":

		break
	default:
		wr.Status = "404 NotFound"
		wr.Write([]byte("Not found."))
	}
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	s := shttp.Server{Addr: "0.0.0.0:8080", ParseReqLength: 128 * 10}
	s.Handler = router
	s.ReqRawHook = func(conn net.Conn) {
		log.Println(conn.RemoteAddr())
	}
	log.Println("Http Server started on", s.Addr)
	go func() {
		if err := s.Listen(); err != nil {
			log.Println(err.Error())
		}
	}()
	<-sig
	log.Println("Shutting down")
	s.Listener.Close()
}


```
