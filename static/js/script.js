var conn;
var sessionID = localStorage.getItem("sessionID");

function getWebSocketURL() {
  var protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  var host = window.location.host; // This gets the current host name and port
  var path = "/chatbot" + (sessionID ? "?sessionID=" + sessionID : "");
  return protocol + "//" + host + path;
}

function connectWebSocket() {
  conn = new WebSocket(getWebSocketURL());

  conn.onopen = function () {
    console.log("Connected to the chat server.");
  };

  conn.onclose = function () {
    console.log("Connection closed");
    if (localStorage.getItem("sessionID")) {
      localStorage.removeItem("sessionID");
      console.log("Session ID removed from localStorage.");
    }
  };

  conn.onerror = function (error) {
    console.log("WebSocket Error: ", error);
  };

  conn.onmessage = function (e) {
    var data = JSON.parse(e.data);
    console.log("Received data: ", data);
    if (data.type === "session_init") {
      sessionID = data.sessionID;
      localStorage.setItem("sessionID", sessionID);
      console.log("Session initialized with ID:", sessionID);
    } else {
      document.getElementById("chatbox").value +=
        "ChatBot: " + data.message + "\n";
    }
  };
}

function sendMessage() {
  var input = document.getElementById("userInput");
  if (input.value.trim() === "") return;
  conn.send(input.value);
  document.getElementById("chatbox").value += "You: " + input.value + "\n";
  input.value = "";
  input.focus();
}

// Initialize WebSocket connection
connectWebSocket();
