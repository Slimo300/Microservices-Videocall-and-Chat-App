import axiosObject, {wsService, wsServiceWebsocket} from "./Setup";


export async function GetWebsocket() {

    let response = await axiosObject.get(wsService+"/accessCode");
    
    let socket = new WebSocket(wsServiceWebsocket+'?accessCode='+response.data.accessCode);
    socket.onopen = () => {
        console.log("Websocket openned");
    };
    socket.onclose = () => {
        console.log("Websocket closed");
    };
    return socket;
}
