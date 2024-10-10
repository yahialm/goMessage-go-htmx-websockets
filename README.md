# Realtime Chat Application

This is a simple real-time chat application built using **Go**, **WebSockets**, **HTMX**, and **Tailwind CSS**.

## Features
- Real-time messaging using WebSockets
- Dynamic front-end built with HTMX and Tailwind CSS
- A chat history that updates automatically without refreshing the page

## Prerequisites

Before running the project, make sure you have the following installed on your machine:

- [Go](https://golang.org/doc/install)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- An internet connection to fetch Tailwind CSS and HTMX libraries.

## Getting Started

### 1. Clone the Repository
First, clone the repository to your local machine:

```bash
git clone <get_repo_name>
cd <your-repo-folder>
```

### 2. Install Dependencies

This project uses the following Go packages:

- [Gorilla WebSocket](https://github.com/gorilla/websocket) - for handling WebSocket connections.
- [Lucasepe codenames](https://github.com/lucasepe/codename) - for generating unique usernames for connected users.

To install these dependencies, run the following commands:

```bash
go get github.com/gorilla/websocket
go get github.com/lucasepe/codename
```

### 3. Running the Application

To run the project, simply use the Go \`run\` command. This will start the server:

```bash
go run main.go
```

This will launch the server on \`localhost:3000\`

### 4. Accessing the Application

Open your browser and navigate to:

```
http://localhost:3000/
```

Here, you'll be able to chat with other connected clients in real-time.

### 5. File Structure

The basic structure of the project is:

```
|-- message.go         # Message struct
├── main.go            # Entry point of the application
├── hub.go             # Handles WebSocket hub logic
├── client.go          # Manages individual WebSocket connections
├── templates
│   └── message.html   # Template for rendering messages in the chat
|   └── index.html     # home page
└── README.md          # This file
```

### 6. Customization

#### Frontend (HTML and CSS)

The front-end uses **HTMX** for real-time message updates and **Tailwind CSS** for styling. To modify the design:

- HTML is located in the \`main.go\` and \`message.html\` template file.
- Tailwind CSS is included via CDN in the \`<head>\` of the HTML file.

#### WebSocket Logic

The WebSocket server is managed in \`hub.go\` (handles registering clients and broadcasting messages) and \`client.go\` (manages individual connections).

### 7. Troubleshooting

If you encounter any issues, make sure you have installed the dependencies and that the WebSocket connection is properly set up.
