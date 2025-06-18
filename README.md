Real-Time Flight Tracker
This project is a web application that displays real-time flight data on a dashboard. It's built with a Go backend that fetches data from the OpenSky Network API and pushes it to clients over WebSockets. The frontend is a Next.js application that renders the data.

Tech Stack
Backend: Go with gorilla/websocket

Frontend: Next.js (React) with Tailwind CSS

Orchestration: Docker Compose

Project Structure
real-time-dashboard/
├── backend/
│ ├── main.go
│ ├── ws/
│ │ └── hub.go
│ ├── fetcher/
│ │ └── flight.go
│ ├── go.mod
│ └── Dockerfile
├── frontend/
│ ├── pages/
│ │ └── index.tsx
│ ├── package.json
│ └── Dockerfile
├── docker-compose.yml
└── README.md

Getting Started
Prerequisites
Docker installed on your machine.

Running Locally
Clone the repository and create the file structure:
Make sure you have all the files provided placed in the correct directories as shown in the project structure above.

Build and run the containers:
From the root directory (real-time-dashboard/), run the following command:

docker-compose up --build

This command will build the Docker images for both the backend and frontend services and then start them.

Access the application:

The Next.js frontend will be available at http://localhost:3000.

The Go backend WebSocket server is running on port 8080.

Open your browser to http://localhost:3000 to see the live flight tracker dashboard.
