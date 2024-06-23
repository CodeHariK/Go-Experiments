package spider

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/codeharik/rerun/logger"
	"github.com/gorilla/websocket"
)

func createServer(s *Spider) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlePage)

	mux.HandleFunc("GET /logs/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.handleLog(w, r)
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.handleWebSocket(w, r)
	})
	server := http.Server{
		Addr:    ":7359",
		Handler: mux,
	}
	return &server
}

func (s *Spider) handleLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	idString := r.PathValue("id")

	fmt.Fprint(w, s.stdOutLogs[idString])
}

func (s *Spider) handleExecute(command string) {
	stdOutLogs := logger.CreateStdOutSave(
		make(map[string][]string),
		func(p []byte) (n int, err error) {
			s.BroadcastMessage(fmt.Sprintf("Console Output %s", string(p)), Connection{ID: "SPIDER"})
			return os.Stdout.Write(p)
		},
	)

	stdErrLogs := logger.CreateStdOutSave(
		make(map[string][]string),
		func(p []byte) (n int, err error) {
			s.BroadcastMessage(fmt.Sprintf("Console Error %s", string(p)), Connection{ID: "SPIDER"})
			return os.Stderr.Write(p)
		},
	)

	// Execute the command
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdOutLogs
	cmd.Stderr = stdErrLogs

	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Printf("Exec Error : %v", err)
	// }

	// s.BroadcastMessage(fmt.Sprintf("Console %s", string(output)), Connection{ID: "SPIDER"})

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error Exec : %v", err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error Wait : %v", err)
	}
}

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
				if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
				) {
					fmt.Printf("Connection closed : %s %v\n", connection.ID, err)
					return
				}
			}

			fmt.Printf("%s Received: %s\n", connection.ID, message)

			command := strings.Split(string(message), ":")
			fmt.Println(command)
			if len(command) > 1 {
				s.handleExecute(command[1])
			}

			s.BroadcastMessage(string(message), connection)
		}
	}()
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	fmt.Fprint(w, htmlContent)
}

const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>ReRun</title>
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
		width: 102%;
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
	#rightPane {
		overflow-y: scroll;
	}
	#statusPane {
		overflow-x: scroll;
	}
	.divider {
		width: 3px;
		background-color: #aaaaaa11;
		cursor: ew-resize;
		position: relative;
		z-index: 1;
	}
	.Output {
		padding: 5px;
		font-family: monospace;
		color: black;
	}
	.Error {
		color: red;
	}

	#commandInput {
		width: 80%;
		padding: 10px;
		font-size: 16px;
		border: 2px solid #ccc;
		border-radius: 4px;
		box-shadow: 2px 2px 12px rgba(0, 0, 0, 0.1);
		outline: none;
		transition: border-color 0.3s;
	}

	#commandInput:focus {
		border-color: #007bff;
	}

	#commandInput::placeholder {
		color: #999;
	}
</style>
<body>
	<div class="container">
        <div class="pane" id="leftPane">
			<iframe id="contentFrame" src=""></iframe>
        </div>
        <div class="divider" id="divider"></div>
        <div class="pane">
			<button onclick="sendMessage()">Send Message</button>
			<div id="terminal">
				<input type="text" id="commandInput" placeholder="Type command and press Enter">
			</div>
			<div id="statusPane">
			</div>
			<div id="rightPane"></div>
        </div>
    </div>
</body>
<script>
	const rightPane = document.getElementById('rightPane');
	const statusPane = document.getElementById('statusPane');

	let socket = new WebSocket("ws://localhost:7359/ws");
	const iframeUrl = 'http://localhost:8080/docs';
	const logsUrl = 'http://localhost:7359/logs/';
	const checkInterval = 100;
	
	async function iframeReload() {
		document.getElementById('contentFrame').src = ""
		try {
			const response = await fetch(iframeUrl, { method: 'HEAD' });
			if (response.ok) {
				document.getElementById('contentFrame').src = iframeUrl;
				return
			}
		} catch (error) {}
		setTimeout(iframeReload, checkInterval);
	}

	socket.onopen = function(event) {
		console.log("Connected to WebSocket spider.");
		iframeReload();
	};

	socket.onmessage = function(event) {
		d = event.data
		console.log("-> " + d);
		console.log(d.startsWith("SPIDER:Console"))
		
		if (d.startsWith("SPIDER:Console")){
			console.log("console terminal : " + Math.random())
			return
		}

		if (d.startsWith("SPIDER:ReRun")){
			iframeReload();
			rightPane.innerHTML = ""
			statusPane.innerHTML = ""

			let matches = d.match(/\d+/);
			if (matches) {
				let num = Number(matches[0]);
				console.log(num);

				for (let i = 1; i <= num; i++) {
					const button = document.createElement('button');
					button.innerText = i;
					button.addEventListener('click', async function(event){
						console.log(i)
						try {
							const response = await fetch(logsUrl+i, { method: 'GET' });
							console.log(logsUrl+i)
							let body = await response.text() 
							if (response.ok) {
								rightPane.innerHTML = body;
							}
						} catch (error) {}
					});
					statusPane.appendChild(button);
				}
			}
		}

		const newLog = document.createElement("div");
		newLog.innerHTML = d.replace(/\n/g, "<br>");
		newLog.classList.add("Output");
		if (d.toLowerCase().includes("error")) {
			newLog.classList.add("Error");
		}
		rightPane.appendChild(newLog);
	};

	socket.onclose = function(event) {
		console.log("Disconnected from Spider." + event);
		location.reload()
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
	rightPane.style.width = container.clientWidth * .3 + 'px'

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


	///////


	const input = document.getElementById('commandInput');
	input.addEventListener('keypress', function (event) {
        if (event.key === 'Enter') {
            const command = input.value;
            if (socket.readyState === WebSocket.OPEN) {
				console.log("Console:"+command)
				socket.send("Console:"+command);
                input.value = '';
            } else {
                console.log("WebSocket is not open.");
            }
        }
    });
</script>
</html>
`
