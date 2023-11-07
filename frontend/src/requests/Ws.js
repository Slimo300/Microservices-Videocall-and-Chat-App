import axiosObject, {webrtcService, webrtcServiceWebsocket, wsService, wsServiceWebsocket} from "./Setup";

export async function GetWebsocket() {

    let response = await axiosObject.get(wsService+"/accessCode");
    
    let socket = new WebSocket(wsServiceWebsocket+'/websocket?accessCode='+response.data.accessCode);
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

export async function GetWebRTCAccessCode(groupID) {

    let response = await axiosObject.get(webrtcService+"/"+groupID+"/accessCode");

    return response.data.accessCode
    
}

export function GetWebRTCWebsocket(groupID, accessCode, streamID) {
    let socket = new WebSocket(webrtcServiceWebsocket+"/"+groupID+"/websocket?accessCode="+accessCode+"&streamID="+streamID);
    
    socket.onopen = () => {
        let date = new Date();
        console.log("Websocket openned\nSocket openned: ", date);
    };

    socket.onclose = (evt) => {
        let date = new Date();
        console.log("WebRTC signaling Websocket closed: ", evt.wasClean, "\ncode: ", evt.code, "\nreason: ", evt.reason, "\ntimestamp: ", date);
    };

    socket.onerror = (evt) => {
        console.log(evt);
    };

    return socket;
}

export async function GetWebRTCWebSocket(groupID, streamID) {

    let {accessCode, username, muting} = (await axiosObject.get(webrtcService+"/"+groupID+"/accessCode")).data;
    console.log(muting)

    let socket = new WebSocket(webrtcServiceWebsocket+"/"+groupID+"/websocket?accessCode="+accessCode+"&streamID="+streamID);
    
    socket.onopen = () => {
        let date = new Date();
        console.log("Websocket openned\nSocket openned: ", date);
    };

    socket.onclose = (evt) => {
        let date = new Date();
        console.log("WebRTC signaling Websocket closed: ", evt.wasClean, "\ncode: ", evt.code, "\nreason: ", evt.reason, "\ntimestamp: ", date);
    };

    socket.onerror = (evt) => {
        console.log(evt);
    };

    return {
        socket, username, muting
    };
}