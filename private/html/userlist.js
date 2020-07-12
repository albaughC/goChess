const ws = new WebSocket("wss://chesslife.herokuapp.com/private/api/userwebsocket");

ws.addEventListener("open", () => {
    console.log('We are connected');
});
