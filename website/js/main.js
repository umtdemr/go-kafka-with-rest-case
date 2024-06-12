window.addEventListener('load', () => {
    main()
})

function main() {
    const wsUrl = "ws://127.0.0.1:8081/ws"
    const ws = new WebSocket(wsUrl)
    ws.onerror = err => {
        console.log('an error is occurred while doing that', err)
    }
    ws.onopen = () => {
        ws.send('getData')
    }
    ws.onmessage = (event) => {
        handleWebsocketMessage(event.data)
    }
}

function handleWebsocketMessage(data) {
    const obj = JSON.parse(data)
    console.log(obj)
}