package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func getTunnelHost() string {
	return tunnelHost
}

func getTargetHost() string {
	return targetHost
}

func logging(level int, a ...any) {
	if verbal && level >= logLevel {
		fmt.Println(a...)
	}
}

func tunneling(port string) {
	conn, err := net.Dial("tcp", getTunnelHost())

	if nil != err {
		logging(10, "Failed to connect server")
		return
	}

	handshakeWriter := newBufferWriter()
	handshakeWriter.writeString("manager")
	handshakeWriter.writeString("manager-sess-" + createUuid())
	handshakeWriter.writeString(alias)
	handshakeWriter.writeString(port)
	write(conn, handshakeWriter.getBytes())

	recv(conn)
}

func write(conn net.Conn, packet []byte) {
	packetWriter := newBufferWriter()
	packetWriter.writeFixBuffer(packet)

	conn.Write(packetWriter.getBytes())
}

func packetProc(buffer []byte) {
	br := newBufferReader(buffer)

	msgType := br.readString()

	if msgType == "newConnection" {

		sessionId := br.readString()

		targetConnection, targetErr := net.Dial("tcp", getTargetHost())
		tunnelConnection, tunnelErr := net.Dial("tcp", getTunnelHost())

		if nil != targetErr {
			logging(10, "Failed to connect target.")
			return
		}

		if nil != tunnelErr {
			logging(10, "Failed to make a new tunnel session.")
			return
		}

		handshakeBuffer := newBufferWriter()
		handshakeBuffer.writeString("agent")
		handshakeBuffer.writeString(sessionId)
		handshakeBuffer.writeString(alias)

		write(tunnelConnection, handshakeBuffer.getBytes())

		go func(from net.Conn, to net.Conn) {
			for {
				data := make([]byte, 1024)
				n, err := from.Read(data)

				if nil != err {
					logging(5, "[ERROR] Target -> Tunnel", err)
					logging(5, "[CLOSE] Close Tunnel Connection")
					to.Close()
					return
				}

				if n > 0 {
					_, err := to.Write(data[:n])
					if nil != err {
						logging(5, "[ERROR] An unexpected error occured when write data to tunnel")
						logging(5, "[ERROR]", err)
					}
				} else {
					return
				}

			}
		}(targetConnection, tunnelConnection)

		go func(from net.Conn, to net.Conn) {
			for {
				data := make([]byte, 1024)
				n, err := from.Read(data)

				if nil != err {
					logging(5, "[ERROR] Tunnel -> Target", err)
					logging(5, "[CLOSE] Close Tunnel Connection")
					from.Close()
					return
				}

				if n > 0 {
					_, err := to.Write(data[:n])
					if nil != err {
						logging(5, "[ERROR] An unexpected error occured when write data to target")
						logging(5, "[ERROR]", err)
					}
				} else {
					return
				}

			}
		}(tunnelConnection, targetConnection)
	} else if msgType == "listenOk" {
		host := br.readString()

		fmt.Println("Tunnel open at", host)
	} else if msgType == "listenFailed" {
		host := br.readString()

		fmt.Println("Failed to open tunnel at", host)
		os.Exit(0)
	}
}

func recv(conn net.Conn) {
	buffer := make([]byte, 1024)
	bw := newBufferWriter()

	for {
		n, err := conn.Read(buffer)

		if nil != err {
			return
		}

		if n > 0 {
			bw.writeBuffer(buffer[:n])
			loops := 0

			for {
				if len(bw.getBytes()) < 4 {
					break
				}

				br := newBufferReader(bw.getBytes())
				packetSize := br.readInt()
				loops += 1

				if packetSize+4 <= int32(br.reader.Size()) {
					packet := make([]byte, packetSize)

					br.reader.Read(packet)
					packetProc(packet)

					cropped, err := io.ReadAll(br.reader)
					if err != nil {
						logging(5, "[ERROR] An unexpected error occured when crop remain datagram.", err)
					}

					bw = newBufferWriter()
					bw.writeBuffer(cropped)
				} else {
					break
				}
			}
		}
	}
}
