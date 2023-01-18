import axiosObject, {wsService, wsServiceWebsocket} from "./Setup";


export async function GetWebsocket() {

    let response = await axiosObject.get(wsService+"/accessCode");
    
    let socket = new WebSocket(wsServiceWebsocket+'?accessCode='+response.data.accessCode);
    socket.onopen = () => {
        let date = new Date();
        console.log("Websocket openned\nSocket openned: ", date);
    };
    socket.onclose = (ev) => {
        let date = new Date();
        console.log("Websocket closed: ", ev.wasClean, "\ncode: ", ev.code, "\nreason: ", ev.reason, "\ntimestamp: ", date);
    };
    socket.onerror = (ev) => {
        console.log(ev)
    }
    return socket;
}
