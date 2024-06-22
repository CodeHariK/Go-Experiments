package spider

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func createServer(s *Spider) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerPage)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.handleWebSocket(w, r)
	})
	server := http.Server{
		Addr:    ":7359",
		Handler: mux,
	}
	return &server
}

// handleWebSocket handles WebSocket connections.
func (s *Spider) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	connection := Connection{
		conn: conn,
		ID:   r.RemoteAddr, // Use the remote address as a simple identifier
	}

	fmt.Printf("Adding %s\n", connection.ID)
	s.addConn <- connection

	go func() {
		defer func() {
			fmt.Printf("Removing %s\n", connection.ID)
			s.removeConn <- connection
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("Connection closed : %s %v\n", connection.ID, err)
					return
				}
			}

			fmt.Printf("%s Received: %s\n", connection.ID, message)

			s.BroadcastMessage(string(message), connection)
		}
	}()
}

func handlerPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	fmt.Fprint(w, htmlContent)
}

const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>My Web Page</title>
</head>
<style>
	html, body {
		height: 100%;
		margin: 0;
		padding: 0;
		overflow: hidden;
	}
	iframe {
		position: absolute;
		bottom: 0;
		left: 0;
		width: 100%;
		height: 100%;
		border: none;
	}
	.container {
		display: flex;
		height: 100%;
	}
	.pane {
		overflow: hidden;
		position: relative;
	}
	.divider {
		width: 3px;
		background-color: #aaaaaa11;
		cursor: ew-resize;
		position: relative;
		z-index: 1;
	}
</style>
<body>
	<div class="container">
        <div class="pane" id="leftPane">
			<iframe id="contentFrame" src=""></iframe>
        </div>
        <div class="divider" id="divider"></div>
        <div class="pane" id="rightPane">
			<button onclick="sendMessage()">Send Message</button>
        </div>
    </div>
</body>
<script>
	const rightPane = document.getElementById('rightPane');

	let socket = new WebSocket("ws://localhost:7359/ws");
	const iframeUrl = 'http://localhost:8080/docs';
	const checkInterval = 100; 
	
	async function iframeReload() {
		document.getElementById('contentFrame').src = ""
		try {
			const response = await fetch(iframeUrl, { method: 'GET' });
			console.log(response)
			if (response.ok) {
				document.getElementById('contentFrame').src = iframeUrl;
				return
			}
		} catch (error) {
			console.log(error)
		}
		setTimeout(iframeReload, checkInterval);
	}

	socket.onopen = function(event) {
		console.log("Connected to WebSocket spider.");
		iframeReload();
	};

	socket.onmessage = function(event) {
		d = event.data
		console.log("-> " + d);
		
		if (d.startsWith("ReRun")){
			iframeReload();
			// location.reload()
		}

		rightPane.innerHTML += event.data + "<br>"
	};

	socket.onclose = function(event) {
		console.log("Disconnected from Spider." + event);
	};

	socket.onerror = function(event) {
		console.error('Spider error:', event);
	};

	function sendMessage() {
		if (socket.readyState === WebSocket.OPEN) {
			socket.send("Hello, spider!");
			console.log("Hello, spider!");
		} else {
			console.log("WebSocket is not open.");
		}
	}


	////////


	const divider = document.getElementById('divider');
	const leftPane = document.getElementById('leftPane');
	const container = document.querySelector('.container')

	leftPane.style.width = container.clientWidth * .7 + 'px'

	let isDragging = false;

	divider.addEventListener('mousedown', function(e) {
		isDragging = true;
		leftPane.style.pointerEvents = 'none';
		document.addEventListener('mousemove', onMouseMove);
		document.addEventListener('mouseup', onMouseUp);
	});

	function onMouseUp() {
		isDragging = false;
		leftPane.style.pointerEvents = 'auto';
		document.removeEventListener('mousemove', onMouseMove);
		document.removeEventListener('mouseup', onMouseUp);
	}

	function onMouseMove(e) {
		console.log(e.clientX)

		if (!isDragging) return;
		const containerOffsetLeft = container.offsetLeft;
		const pointerRelativeXpos = e.clientX - containerOffsetLeft;
		const containerWidth = container.clientWidth;
		const dividerWidth = divider.offsetWidth;
		const minLeftPaneWidth = 100;
		const minRightPaneWidth = 100;

		if (pointerRelativeXpos < minLeftPaneWidth || pointerRelativeXpos > containerWidth - minRightPaneWidth - dividerWidth) {
			return;
		}

		const leftPaneWidth = pointerRelativeXpos;
		leftPane.style.width = e.clientX - containerOffsetLeft + 'px';
		rightPane.style.width = containerWidth - e.clientX - dividerWidth + 'px';
	}
</script>
</html>
`
