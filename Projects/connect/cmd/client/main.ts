import {
    createPromiseClient
} from '@connectrpc/connect'
import {
    createConnectTransport,
} from '@connectrpc/connect-web'
import { GreetService } from '../../gen/greet/v1/greet_connect.ts'
import { GreetRequest } from '../../gen/greet/v1/greet_pb.ts'


// Make the Eliza Service client
const client = createPromiseClient(
    GreetService,
    createConnectTransport({
        baseUrl: 'http://localhost:8080',
    })
)

// Query for the common elements and cache them.
const containerEl = document.getElementById("conversation-container") as HTMLDivElement;
const inputEl = document.getElementById("user-input") as HTMLInputElement;

// Add an event listener to the input so that the user can hit enter and click the Send button
document.getElementById("user-input")?.addEventListener("keyup", (event) => {
    event.preventDefault();
    if (event.key === "Enter") {
        document.getElementById("send-button")?.click();
    }
});

// Adds a node to the DOM representing the conversation with Eliza
function addNode(text: string, sender: string): void {
    const divEl = document.createElement('div');
    const pEl = document.createElement('p');

    const respContainerEl = containerEl.appendChild(divEl);
    respContainerEl.className = `${sender}-resp-container`;

    const respTextEl = respContainerEl.appendChild(pEl);
    respTextEl.className = "resp-text";
    respTextEl.innerText = text;
}

async function send() {
    const sentence = inputEl?.value ?? '';

    addNode(sentence, 'user');

    inputEl.value = '';


    const response = await client.greet({
        name: sentence,
    })

    console.log(response.toJsonString())

    addNode(response.greeting, 'eliza');
}

export function handleSend() {
    send();
}
