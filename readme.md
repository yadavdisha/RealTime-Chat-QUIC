## Student Collaboration Chat Protocol Implementation using QUIC in GO

## Overview
The Student Collaboration Chat Protocol Application is a network-based chat application designed to facilitate communication among students. It uses the QUIC protocol for secure, reliable communication and supports both server and client modes. This project is implemented in Go.



## QUIC Features in Chat Protocol

#### Real-Time Messaging

- **Basic Messaging:** Implement basic functionality for sending and receiving messages between clients.
- **QUIC Protocol:** Utilizes the `quic-go` library to handle QUIC sockets for efficient and reliable communication.
- **Message Encoding/Decoding:** Custom PDU (Protocol Data Unit) format is used for message encoding and decoding to ensure efficient data transmission.
- **Emoji Support:** Supports sending and receiving messages with emojis as part of the Unicode standard.

#### User Management

- **Username Handling:** Users can join the chat with a specified username.
- **User Join/Leave Notifications:** Notifies all users when someone joins or leaves the chat.

#### Broadcast Messaging

- **Message Broadcasting:** Sends messages to all connected clients, ensuring everyone in the chat room receives the messages.

#### Active Users Listing

- **List Active Users:** Allows users to request a list of active users in the chat room.

#### Authentication

- **TLS for Secure Connections:** Incorporated TLS to secure connections. The client uses a certificate file for authentication, ensuring that only authorized clients can connect to the server. If a certificate file is not provided, the client defaults to using a basic TLS configuration.

#### Multi-client Support

- **Handling Multiple Clients:** The server can handle multiple clients simultaneously. Each client initiates a connection and opens a stream to communicate with the server. The server can manage multiple streams concurrently, allowing multiple clients to communicate at the same time.

#### Error Handling

- **Error Logging:** Error handling is critical in our protocol. Both the client and server log any errors encountered during communication. If an error occurs, such as an issue with reading from a stream or decoding a PDU, it is logged.

#### Chat Messages from Client to Server

- **Initial Handshake and Message Exchange:** Clients can send chat messages to the server, and the server broadcasts these messages to all connected clients. This ensures real-time communication among all participants in the chat room


## Protocol Requirements

**STATEFUL:** Both the client and server implement a stateful protocol, ensuring that the protocol adheres to our deterministic finite automaton (DFA). This ensures that both ends of the communication can validate the state of the connection at any given time.

## States of the DFA:
- **Initial State:** The client sends a "hello" message to the server.
- **Waiting for Server Response:** The server responds with a welcome message and user join notification.
- **Waiting for Client Message:** The client sends their chat message.
- **Broadcasting Message:** The server broadcasts the message to all connected clients.
- **Handling Controls:** The client sends control commands such as joining and leaving the chat.
- **Error State:** If an error occurs, the connection logs the error and attempts to recover.
- **Waiting for Acknowledgments:** The server waits for acknowledgments from the client to ensure message delivery.

**SERVICE:** The server binds to a default port number (4242), with the client defaulting to this port, configurable via command line arguments for flexibility.

**CLIENT:** The client allows specifying the server's hostname or IP address through command line arguments, ensuring adaptability to different network configurations.

**SERVER:** The server configuration, including certificate file, key file, address, and port, is provided via command line arguments, avoiding hard-coded values and enhancing configurability.

**UI:** The client uses a command line interface, encapsulating the protocol logic and providing a user-friendly experience without exposing protocol commands to the user.

## Extra Credit Options Implemented
**Concurrent or Asynchronous Server**: Our server implementation is capable of handling multiple clients concurrently using the quic-go library, which allows managing multiple streams over a single connection efficiently.

**Using a Systems Programming Language**: This project is implemented in Go, a systems programming language, which adds complexity and robustness compared to using a high-level language.

**Implementation Robustness**: We have implemented several features including real-time messaging, user management, TLS encryption, and error handling, demonstrating a robust and comprehensive solution beyond a basic prototype.

**Working with a Cloud-Based Git System**: Our project is managed using GitHub, with regular commits and proper version control. 

**Design Excellence**: We strive for high code quality with clear documentation, maintainability, and a clean project structure. The code is modular, well-documented, and follows best practices, making it easy to understand and extend

**A Short Video Demo Presentation Link**: https://shorturl.at/wkyqj 
https://1513041.mediaspace.kaltura.com/media/Student+Collaboration+Chat+Protocol+Implementation+in+QUIC+using+GO/1_ioeyoux9



## Configuration & Implementation

### Server Configuration
The server binds to a configurable port number, which can be specified via command line arguments. It uses the following settings:

- **GenTLS:** Flag to generate TLS config.
- **CertFile:** Path to the certificate file.
- **KeyFile:** Path to the key file.
- **Address:** Server address.
- **Port:** Server port.

### Client Configuration
The client can specify the server's hostname or IP address via command line arguments. The configuration includes:

- **ServerAddr:** Server address.
- **PortNumber:** Server port.
- **CertFile:** Path to the certificate file.
- **Username:** Username for the chat.

## Initial Configuration in Echo

### Server Parameters
- **PORT_NUMBER:** 4242
- **SERVER_IP:** "0.0.0.0"

### Client Parameters
- **SERVER_ADDR:** "localhost"

## Command Line Output
With the command line, the user enters their messages, and the application processes and displays them.

## Running the Application
There is a single binary that is used to run both the client and the server:

- **Server:** `go run cmd/echo/echo.go -server`
- **Client:** `go run cmd/echo/echo.go -client -username "username" (you can specify any username)`






  






