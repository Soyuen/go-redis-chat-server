
# Go Redis Chat Server
### TOC

- [Introduction](#introduction)
- [Features](#features)
- [Quick Start](#-quick-start)
- [Sequence Diagram](#sequence-diagram)
- [Ideas and Possible Features](#ideas-and-possible-features)



## Introduction

Redis Chat Server is a real-time chat server built on Redis and WebSocket, supporting multi-channel message broadcasting and client connection management. It serves as a backend core for real-time communication systems.

---

## Features

- **Real-time messaging** — Support sending and receiving messages instantly.  
- **Multi-channel support** — Users can chat in different channels at the same time.  
- **Client connection management** — Manage user connections and their lifecycle.  
- **Broadcast messages** — Send messages to all users or to specific channels.  
- **WebSocket based** — Uses WebSocket for stable and low-latency connections.  
- **Redis backend with Pub/Sub** — Uses Redis Pub/Sub for fast and reliable message delivery.  
- **Designed for future scalability** — Built to support easy expansion and distributed deployment.  
- **Modular and extensible codebase** — Clear code structure that is easy to maintain and extend.  

---

## 🚀 Quick Start

### 1. Create `.env`

You can use the following defaults or override them:

```env
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
PORT=8080
APP_DEBUG=true
LOG_TO_FILE=true
LOG_FILE_PATH=logs/app.log
LOG_ERROR_FILE_PATH=logs/app_error.log
```

### 2. Start with Docker

Use the following commands to build and run the server:

```docker
docker compose build --no-cache
docker compose up -d
```
These commands will build the server without using any cache and run it in detached mode.

Server will be available at:
```localhost
http://localhost:8080
```

---
## 🚀 Demo Preview
<img src="https://github.com/Soyuen/picture/blob/main/demo.gif?raw=true" alt="consent_screen" width="600"/>

The frontend and public entry point are currently not open to avoid abuse and excessive traffic.  The demo shown above is for preview purposes only.

## Sequence Diagram
<img src="https://github.com/Soyuen/picture/blob/main/SequenceDiagram.png?raw=true" alt="consent_screen" width="600"/>


## Ideas and Possible Features

This is a list of potential improvements. Some may be implemented in the future.

- [x] Add system messages for user join/leave events
- [ ] Graceful shutdown
- [ ] Display number of users in each channel
- [ ] Set password for joining a channel
- [ ] Name validation mechanism  
- [ ] User login and authentication  
- [ ] User presence tracking and offline message storage  
- [ ] Support for distributed deployment (Redis cluster / multi-instance sync)
