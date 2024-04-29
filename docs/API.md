# API Documentation

This document outlines the API endpoints provided by the chat-bot service. It includes details on the request methods, endpoints, expected request bodies, and sample responses.

## Table of Contents

- [WebSocket API](#websocket-api)
  - [Handle WebSocket Connections](#handle-websocket-connections)
- [REST API](#rest-api)
  - [Session Management](#session-management)

---

## WebSocket API

### Handle WebSocket Connections

This endpoint manages the WebSocket connections for real-time chat functionality.

- **URL**: `/ws`
- **Method**: `GET`
- **Query Parameters**:
  - `sessionID` (optional): A unique identifier for the user's session. If not provided, a new session will be initialized.
- **Success Response**:
  - **Code**: 101 Switching Protocols
  - **Content**:
    ```json
    {
      "type": "session_init",
      "sessionID": "generated-session-id"
    }
    ```
- **Error Response**:
  - **Code**: 400 BAD REQUEST
  - **Content**: Error message in plain text or JSON format.
- **Sample Call**:
  ```javascript
  const socket = new WebSocket(
    "wss://yourdomain.com/ws?sessionID=existing-session-id"
  );
  socket.onmessage = function (event) {
    console.log("Message from server ", event.data);
  };
  ```
- **Notes**:
  - Connections are upgraded to WebSocket from a standard HTTP request. Ensure WebSocket support is enabled on the client side.
  - The server checks if a sessionID is valid or not and initializes a session if necessary.

## REST API

### Session Management

Endpoints for managing user sessions, typically used for setting up or clearing sessions prior to or after WebSocket communications.

- **Endpoint**: `/sessions`
- **Method**: `POST`
- **Body**:
  ```json
  {
    "action": "initialize"
  }
  ```
- **Success Response**:
  - **Code**: 200 OK
  - **Content**:
    ```json
    {
      "sessionID": "new-session-id",
      "status": "initialized"
    }
    ```
- **Error Response**:
  - **Code**: 500 INTERNAL SERVER ERROR
  - **Content**: `{ "error": "Failed to initialize session" }`
- **Sample Call**:
  ```bash
  curl -X POST https://yourdomain.com/sessions -d '{"action":"initialize"}'
  ```
- **Notes**:
  - This endpoint is used to initialize a new session or manage existing ones. Additional endpoints may be created to handle other session-related actions such as renewal or termination.
