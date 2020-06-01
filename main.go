package wshttpserver

import (
	"bytes"
	"log"
	"net"
	"strings"
)

// Server represents a HTTP server
type Server struct {
	Listener       net.Listener
	Addr           string
	ReqRawHook     func(net.Conn)
	ParseReqLength int
	Handler        func(Request, Response)
}

// Request - HTTP Client request
type Request struct {
	Path     string
	Headers  map[string]string
	Method   string
	Protocol string
}

// Response - HTTP Server response
type Response struct {
	Headers  map[string]string
	Protocol string
	Status   string
	Conn     net.Conn
}

func (wr *Response) Write(b []byte) {
	wr.Conn.Write([]byte(strings.Join([]string{wr.Protocol, " ", wr.Status, "\n\n", string(b)}, "")))
	wr.Conn.Close()
}

func (s *Server) rawHandle(conn net.Conn) {
	s.ReqRawHook(conn)
	var r Request = Request{Headers: make(map[string]string, 30)}
	if s.ParseReqLength == 0 {
		s.ParseReqLength = 128 * 10
	}
	b := make([]byte, s.ParseReqLength)
	if n, err := conn.Read(b); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Read", n, "bytes")
	}
	lines := bytes.Split(b, []byte("\n"))
	for k, ln := range lines {
		//log.Println(string(ln))
		if k == 0 {
			line := strings.Split(string(ln), " ")
			r.Method = line[0]
			r.Path = line[1]
			r.Protocol = line[2]
		} else {
			line := strings.Split(string(ln), ":")
			if len(line) > 1 {
				r.Headers[line[0]] = strings.TrimSpace(line[1])
			} else {
				r.Headers[line[0]] = ""
			}
		}
	}
	wr := Response{Conn: conn}
	wr.Status = "200 OK"
	wr.Protocol = "HTTP/1.1"
	s.Handler(r, wr)
}

func (s *Server) bind() error {
	var err error
	s.Listener, err = net.Listen("tcp", s.Addr)
	return err
}

// Listen binds the tcp listener and starts an accept loop
func (s *Server) Listen() error {
	if err := s.bind(); err != nil {
		log.Println(err.Error())
		return err
	}
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		s.rawHandle(conn)
	}
}
