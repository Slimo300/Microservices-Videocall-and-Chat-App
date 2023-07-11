import React, {useEffect, useState, useRef} from 'react';
import {  useParams, Navigate } from "react-router-dom";

import useQuery from '../hooks/useQuery';
import "../Call.css";
import { GetWebRTCWebsocket } from '../requests/Ws';
import CallScreen from '../components/videocall/CallScreen';
import StartCall from '../components/videocall/StartCall';


const VideoConference = () => {

    const { id } = useParams();
    const accessCode = useQuery().get("accessCode");
    const mocking = useQuery().get("mock");

    const peerConnection = useRef(new RTCPeerConnection());
    const ws = useRef(null);
    const audio = useRef({});
    const video = useRef({});

    const [fatal, setFatal] = useState(false);

    const [RTCStreams, setRTCStreams] = useState({});
    const [userStream, setUserStream] = useState(null);

    useEffect(() => {
        const startCall = async () => {
            const stream = await StartCall(mocking);

            setUserStream(stream);

            peerConnection.current.ontrack = (event) => {
                console.log("New track received: ", event.track);
                setRTCStreams(streams => {
                    if (!streams[event.streams[0].id]) {
                        streams[event.streams[0].id] = event.streams[0];
                    }
                    return {...streams};
                });
    
                event.streams[0].onremovetrack = () => {
                    console.log("Track removed");
                    if (!event.streams[0].active) {
                        setRTCStreams(streams => {
                            delete streams[event.streams[0].id];
                            return {...streams};
                        });
                    }
                }
            };
    
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

    if (fatal) return <Navigate to="/not-found" />;

    return (
        <CallScreen ws={ws} peerConnection={peerConnection} stream={userStream} video={video} audio={audio} RTCStreams={RTCStreams} setRTCStreams={setRTCStreams}/>
    )
};

export default VideoConference;