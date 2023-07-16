import React, {useEffect, useState, useRef, useCallback, useReducer} from 'react';
import {  useParams, Navigate } from "react-router-dom";

import useQuery from '../hooks/useQuery';
import "../Call.css";
import { GetWebRTCWebsocket } from '../requests/Ws';
import CallScreen from '../components/videocall/CallScreen';
import StartCall from '../components/videocall/StartCall';
import { RTCStreamsReducer, actionTypes } from '../components/videocall/RTCStreams';


const VideoConference = () => {

    const { id } = useParams();
    const accessCode = useQuery().get("accessCode");
    const mocking = useQuery().get("mock");

    const peerConnection = useRef(new RTCPeerConnection());
    const ws = useRef(null);
    const audio = useRef({});
    const video = useRef({});

    const [dataChannel, setDataChannel] = useState(null);

    const [fatal, setFatal] = useState(false);

    const [RTCStreams, dispatch] = useReducer(RTCStreamsReducer, []);

    const [userStream, setUserStream] = useState(null);

    useEffect(() => {
        const startCall = async () => {
            const stream = await StartCall(mocking);

            setUserStream(stream);

            peerConnection.current.ontrack = (event) => {
                dispatch({type: actionTypes.NEW_STREAM, payload: event.streams[0]});
    
                event.streams[0].onremovetrack = () => {
                    dispatch({type: actionTypes.DELETE_STREAM, payload: event.streams[0].id});
                }
            };

            peerConnection.current.ondatachannel = e => {
                e.channel.onopen = evt => {
                    e.channel.send(JSON.stringify({
                        "type": "NewUser",
                        "data": {
                            "username": localStorage.getItem("username"),
                            "streamID": stream.id,
                        },
                    }));
                };

                e.channel.onmessage = evt => {
                    const msgJSON = JSON.parse(evt.data);

                    switch (msgJSON.type) {
                        case "NewUser":
                            if (msgJSON.data.streamID === stream.id) return;
                            dispatch({type: actionTypes.SET_USERNAME, payload: msgJSON.data});
                            break
                        default:
                            console.error("Unsupported message type: ", msgJSON.type);
                    }
                }
        
                setDataChannel(e.channel);
            }
    
            audio.current.track = stream.getAudioTracks()[0];
            audio.current.sender = peerConnection.current.addTrack(audio.current.track, stream);
    
            video.current.track = stream.getVideoTracks()[0];
            video.current.sender = peerConnection.current.addTrack(video.current.track, stream);
            video.current.screenshare = false;
    
            try {
                ws.current = GetWebRTCWebsocket(id, accessCode);
            } catch(err) {
                alert(err);
                setTimeout(() => setFatal(true), 3000);
                return;
            }
    
            peerConnection.current.onicecandidate = (event) => {
                if (!event.candidate) return;
                ws.current.send(JSON.stringify({event: "candidate", data: JSON.stringify(event.candidate)}));
            };
    
            ws.current.onmessage = (event) => {
                let msg = JSON.parse(event.data);
                if (!msg) {
                    return console.log("failed to parse msg");
                }
        
                switch(msg.event) {
                    case "offer":
                        let offer = JSON.parse(msg.data);
                        if (!offer) {
                            return console.log("Failed to parse message");
                        }
                        peerConnection.current.setRemoteDescription(offer);
                        peerConnection.current.createAnswer().then(answer => {
                            peerConnection.current.setLocalDescription(answer);
                            ws.current.send(JSON.stringify({event: "answer", data: JSON.stringify(answer)}));
                        });
                        return;
                    case "candidate":
                        let candidate = JSON.parse(msg.data);
                        if (!candidate) {
                            return console.log("Failed to parse candidate");
                        }
        
                        peerConnection.current.addIceCandidate(candidate);
                        break;
                    default:
                        console.log("Unexpected websocket event: ", msg.event);
                }
            };

        };
        
        if (!window.localStorage.getItem("token")) {
            setFatal(true);
            return;
        }

        startCall();

    }, [accessCode, id, mocking]);

    const EndSession = useCallback(() => {
        peerConnection.current.close();
        ws.current.close();

        dispatch({type: actionTypes.END_SESSION});

    }, [peerConnection, ws]);

    if (fatal) return <Navigate to="/not-found" />;

    return (
        <CallScreen dataChannel={dataChannel} endSession={EndSession} stream={userStream} video={video} audio={audio} RTCStreams={RTCStreams} />
    )
};

export default VideoConference;