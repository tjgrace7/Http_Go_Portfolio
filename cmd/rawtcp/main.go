package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var listener net.Listener
var err error
var serverStartTime time.Time

// Listen starts a TCP listener on the specified network and address, and returns the listener and any error encountered.
func Listen(network, address string) (net.Listener, error) {
	serverStartTime = time.Now()
	listener, err = net.Listen(network, address)

	fmt.Print("Binding to server...")
	return listener, err
}

// Creates a Standard HTTPResponse with the given status code, status text, body, HTTP version, and content type.
func HTTPResponse(statusCode int, statusText string, body string, contentType string) string {
	return fmt.Sprint("HTTP/1.1", " ", statusCode, " ", statusText, "\r\n",
		"Content-Type: ", contentType, "\r\n",
		"Content-Length: ", fmt.Sprint(len(body)), "\r\n",
		"\r\n",
		body)
}

type Request struct {
	Address       string
	Version       string
	Method        string
	Headers       map[string]string
	Body          string
	ContentLength string
	Agent         string
	ContentType   string
	Authorization string
	HeaderParsed  bool
}

// RequestHandler processes the incoming request bytes, extracts the method, address, version, headers, and body, and returns a Request struct and a string indicating whether to continue processing or return an error response.
func RequestHandler(parts []string, r Request) Request {
	requestLine := make(map[string]string)
	headersparsed := false
	for i, line := range parts[1:] {
		if line == "" {
			// The next line is the body
			headersparsed = true
			if i+2 < len(parts) {
				body := strings.Join(parts[i+2:], "\r\n")
				requestLine["Body"] = body
				break
			}
		}
		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		requestLine[key] = value

	}
	s := fmt.Sprint(len(requestLine["Body"]))
	fmt.Println("Content-Type:", requestLine["Content-Type"], "Body Length:", s, "Content-Length:", requestLine["Content-Length"], "Agent:", requestLine["User-Agent"], "Authorization:", requestLine["Authorization"])
	contentlength := requestLine["Content-Length"]
	if r.ContentLength != "" {
		contentlength = r.ContentLength
	}
	agent := requestLine["User-Agent"]
	if r.Agent != "" {
		agent = r.Agent
	}
	contenttype := requestLine["Content-Type"]
	if r.ContentType != "" {
		contenttype = r.ContentType
	}
	authorization := requestLine["Authorization"]
	if r.Authorization != "" {
		authorization = r.Authorization
	}
	return Request{
		Headers:       requestLine,
		Body:          requestLine["Body"],
		ContentLength: contentlength,
		Agent:         agent,
		ContentType:   contenttype,
		Authorization: authorization,
		HeaderParsed:  headersparsed,
	}
}

func GetServerAddress(address string) string {
	return HTTPResponse(200, "OK", "address = "+address, "text/plain")
}

// Get Request Uptime returns the server uptime in a formatted HTTP response.
func GetServerUptime() string {
	uptime := time.Since(serverStartTime)
	return HTTPResponse(200, "OK", "uptime = "+uptime.String(), "text/plain")
}

func PostSubmit(r Request) string {
	s := fmt.Sprintf("%d", len(r.Body))
	return HTTPResponse(200, "OK", "POST body length "+s, "text/plain")
}

// Starts the Server and listens for incoming connections, handling each connection in a separate goroutine.
func runServer(listen net.Listener) {
	defer listen.Close()
	fmt.Println("Server is listening on", listen.Addr())
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go HandleConnection(conn)
	}
}

func Authenticator(r Request) (string, bool) {
	maintain := true
	response := HTTPResponse(200, "Ok", "Keep Going", "text/plain")
	api_key := os.Getenv("API_KEY")
	if api_key == "" {
		fmt.Println("API_KEY not found in .env file")
		response = HTTPResponse(500, "500 Internal Server Error", "Internal API-KEY Reference Not Found", "text/plain")
		maintain = false
		return response, maintain
	}
	if r.Method == "" || r.Address == "" || r.Version == "" {
		fmt.Println("Invalid request: missing method, address, or version")
		response = HTTPResponse(400, "Bad Request", "400 Bad Request", "text/plain")
		maintain = false
		return response, maintain

	}
	if r.Authorization == "" {
		fmt.Println("Missing API Key")
		response = HTTPResponse(400, "Bad Request", "400 Bad Request: Missing API Key", "text/plain")
		maintain = false
		return response, maintain
	}
	token, found := strings.CutPrefix(r.Authorization, "Bearer ")
	if !found {
		response = HTTPResponse(400, "Bad Request", "400 Bad Request: Impropert API Format", "text/plain")
		maintain = false
		return response, maintain
	}
	if token != api_key {
		fmt.Println("Invalid API Key")
		response = HTTPResponse(401, "Unauthorized", "401 Unauthorized", "text/plain")
		maintain = false
		return response, maintain
	}
	return response, maintain
}

func bufferloop(conn net.Conn, r Request) (Request, error) {
	for !r.HeaderParsed {
		fmt.Println("Headers Loop")
		headerbuf := make([]byte, 1024)
		n, err := conn.Read(headerbuf)
		if err != nil {
			break
		}
		parts := strings.Split(string(headerbuf[:n]), "\r\n")
		r = RequestHandler(parts, r)
	}
	if r.ContentLength == "" {
		return r, fmt.Errorf("Missing Content Length")
	}
	maxlength, lenerr := strconv.Atoi(r.ContentLength)
	if lenerr != nil {
		return r, fmt.Errorf("Content-Length Not Integer")
	}
	for maxlength > len(r.Body) {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		//Checks headers to see if there is something missing as well as body to see if it's missing to run it through RequestHandler. Otherwise appends bytes to r.Body
		if !r.HeaderParsed {
			parts := strings.Split(string(buf[:n]), "\r\n")
			r = RequestHandler(parts, r)
		} else {
			r.Body = r.Body + string(buf[:n])
			fmt.Println("The difference between max length and r.Body is ", maxlength-len(r.Body))
		}

	}
	if maxlength != len(r.Body) {
		return r, fmt.Errorf("Full Body Not Received")
	}
	fmt.Println("Loop Finished")
	return r, nil
}

// Handle Connection processes an incoming connection, breaks the connection into parts, determines the HTTP method, address, and version. Responds accurately according to the request.
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connection from", conn.RemoteAddr())
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	var r Request
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	parts := strings.Split(string(buf[:n]), "\r\n")
	if len(strings.Split(parts[0], " ")) < 3 {
		conn.Write([]byte(HTTPResponse(404, "Not Found", "404 Not Found", "text/plain")))
		return
	}
	//Parses submitted request into usable data
	r = RequestHandler(parts, r)
	r, err = bufferloop(conn, r)
	if err != nil {
		conn.Write([]byte(HTTPResponse(400, "400 Inccorect Request", err.Error(), "text/plain")))
		return
	}

	r.Method = strings.TrimSpace(strings.Split(parts[0], " ")[0])
	r.Address = strings.TrimSpace(strings.Split(parts[0], " ")[1])
	r.Version = strings.TrimSpace(strings.Split(parts[0], " ")[2])
	//Authenticates request
	auth, check := Authenticator(r)
	//kills connection if request is denied
	if !check {
		conn.Write([]byte(auth))
		return
	}
	if r.Method == "GET" {
		fmt.Print("Address", r.Address)
		var response string
		switch r.Address {
		case "/address":
			response = GetServerAddress(conn.LocalAddr().String())

		case "/uptime":
			response = GetServerUptime()

		default:
			response = HTTPResponse(404, "Not Found", "404 Not Found", "text/plain")
		}
		fmt.Print("\nResponse:", response)
		conn.Write([]byte(response))
	} else if r.Method == "POST" {
		fmt.Print("Address", r.Address)
		var response string
		switch r.Address {
		case "/address":
			response = GetServerAddress(conn.LocalAddr().String())

		case "/submit":
			response = PostSubmit(r)

		default:
			response = HTTPResponse(404, "Not Found", "404 Not Found", "text/plain")
		}
		fmt.Print("\nResponse:", response)
		conn.Write([]byte(response))
	} else {
		conn.Write([]byte(HTTPResponse(405, "Method Not Allowed", r.Method+"Method Not Allowed", "text/plain")))
	}
}

// Main function initializes the server, sets the default port, starts listening for incoming connections, and runs the server.
func main() {
	godotenv.Load()
	port := "8080" // Default port if not specified in .env

	listener, err = Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Listen Error:", err)
		return
	}
	defer listener.Close()

	runServer(listener)

}
