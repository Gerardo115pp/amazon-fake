export const TransmisorStates = {
    DISCONNECTED: 0,
    CONNECTED: 1
}

class Transmisor {

    constructor() {
        this.status = 0;
        this.host = null;
        this.socket = null;
        this.disconnection_callback = undefined;
        this.connection_callback = undefined;
        this.message_callback = undefined;
        this.error_callback = undefined;
    }

    isConnected = () => this.status === TransmisorStates.CONNECTED ? true : false;

    connect = (host) => {
        this.host = host;
        this.socket = new WebSocket(host);
        this.socket.onopen = this.onConnect;
        this.socket.onmessage = this.onMessage;
        this.socket.onclose = this.onDisconnect;
    }



    setDisconnectionCallback = callback => {
        if (!this.is_connected) {
            this.disconnection_callback = callback;
        }
    }

    setConnectionCallback = callback => {
        if (!this.is_connected) {
            this.connection_callback = callback;
        }
    }

    onConnect = e => {
        if(this.connection_callback !== undefined) {
            this.connection_callback(e);
        }
    }

    onDisconnect = e => {
        if(this.disconnection_callback !== undefined) {
            this.disconnection_callback(e);
        }
        this.status = TransmisorStates.DISCONNECTED;
    }

    Close = () => this.socket.close(1000, "User logged out")

    onMessage = message_event => {
        if (this.message_callback !== undefined) {
            this.message_callback(message_event)
        } 
    }

    setMessageCallback = callback => this.message_callback = callback;
    
    onError = callback => this.error_callback = callback;

    emit = message => {
        if(message.length > 0) {
            this.socket.send(message);
        } else {
            console.warn(`Message '${message}' is invalid`);
        }
    }
}


export default Transmisor;