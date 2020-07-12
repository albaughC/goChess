const ws = new WebSocket("ws://chesslife.herokuapp.com/private/api/userwebsocket");

ws.addEventListener("open", () => {
    console.log('We are connected');
});
