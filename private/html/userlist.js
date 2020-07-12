const ws = new WebSocket("ws://chesslife.herokuapp.com/public/api/userwebsocket");

ws.addEventListener("open", () => {
    console.log('We are connected');
});
