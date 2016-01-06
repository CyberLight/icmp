package main

import (
    "fmt"
    "net"
    "net/http"
    "golang.org/x/net/websocket"
    "time"
)

var (
    inChannel chan string
    connections int = 0
)


func EchoServerConnection(ws *websocket.Conn) {
    defer func(){
        connections--
        fmt.Println("EchoServerConnection! Connections: ", connections)
    }()
    connections++
    fmt.Println("New Connection! Now all connections: ", connections)
    msg := make([]byte, 2048)
    var expectedBytesCount string
    var n int
    var err error
    if n, err = ws.Read(msg); err != nil {
        panic(err)
    }

    expectedBytesCount = string(msg[:n])
    fmt.Println("Expected bytes count: ", expectedBytesCount)

    select {
    case realBytesCount := <-inChannel:
        fmt.Println("Real bytes count: ", realBytesCount)
        if expectedBytesCount == realBytesCount {
            ws.Write([]byte("OK"));
            break;
        }
    case <-time.After(time.Second * 10):
        fmt.Println("timeout 10 seconds")
    }
}

func main() {
    protocol := "icmp"
    conn, _ := net.ListenPacket("ip4:"+protocol, "0.0.0.0")
    inChannel = make(chan string, 10)
    buf := make([]byte, 1024)
    defer func() {
        conn.Close()
        fmt.Println("icmp Connection Closed!")
    }()

    go func() {
        for {
            n, addr, _ := conn.ReadFrom(buf)
            fmt.Println(addr.String())
            fmt.Println(buf[:n], n)
            fmt.Println("Main loop: Connections: ", connections)
            if connections > 0 {
                inChannel <- fmt.Sprintf("%v", n - 8)
            }else {
                fmt.Println("Main loop: Count bytes received: ", n - 8);
            }
        }
    }()

    http.Handle("/echo", websocket.Handler(EchoServerConnection))
    err := http.ListenAndServe(":12345", nil)
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}
