import React, {useEffect, useState, useMemo, useRef} from 'react';
import {  useParams, useLocation, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faMicrophone, faVideo } from '@fortawesome/free-solid-svg-icons';

import "../Call.css";
import { GetWebRTCWebsocket } from '../requests/Ws';

const ConferenceWsWrapper = () => {

    const [ws, setWs] = useState({});
    const { id } = useParams();

    function useQuery() {
        const { search } = useLocation();
      
        return useMemo(() => new URLSearchParams(search), [search]);
    }
    const accessCode = useQuery().get("accessCode");

    useEffect(() => {
        let ws = GetWebRTCWebsocket(id, accessCode);
        setWs(ws);
    }, []);

    console.log(ws);

    if (ws && ws.readyState && ws.readyState === ws.CLOSED) {
        console.log("NOT FOUND WS")
        return <Navigate to="/not-found" />
    }
        
    if (window.localStorage.getItem("token") === null) {
        console.log("NOT FOUND TOKEN")
        return <Navigate to="/not-found" />
    }
    return (
        <VideoConference ws={ws} />
    )
}

const VideoConference = ( {ws} ) => {

    const [streams, setStreams] = useState({});
    const [userStream, setUserStream] = useState({});

    const mediaVideo = useRef();

    useEffect(() => {
        const startConnection = async () => {
            let stream = await navigator.mediaDevices.getUserMedia({video: true, audio: true});
            mediaVideo.current.srcObject = stream;
            setUserStream(stream);
            let peerConnection = new RTCPeerConnection();

            peerConnection.ontrack = function (event) {
                let streamID = event.streams[0].id;
                if (streams[streamID] === undefined) {
                    let newStreams = streams;
                    newStreams[streamID] = event.streams[0];
                    setStreams(newStreams);
                }

                event.track.onmute = function(event) {
                    // TODO: if video is muted display user.png
                }

                event.streams[0].onremovetrack = ({track}) => {
                    if (!event.streams[0].active) {
                        let newStreams = streams;
                        delete newStreams[event.streams[0].id];
                        setStreams(newStreams);

                        return;
                    }
                    // TODO: if video is removed displayed 
                }
            }

            peerConnection.onicecandidate = e => {
                if (!e.candidate) {
                    return;
                }
                ws.send(JSON.stringify({event: "candidate", data: JSON.stringify(e.candidate)}));
            }

            ws.onmessage = function(evt) {
                let msg = JSON.parse(evt.data);
                if (!msg) {
                    return console.log("failed to parse msg");
                }

                switch(msg.event) {
                    case "offer":
                        let offer = JSON.parse(msg.data);
                        if (!offer) {
                            return console.log("Failed to parse message");
                        }
                        peerConnection.setRemoteDescription(offer);
                        peerConnection.createAnswer().then(answer => {
                            peerConnection.setLocalDescription(answer);
                            ws.send(JSON.stringify({event: "answer", data: JSON.stringify(answer)}));
                        });

                        return;
                    case "candidate":
                        let candidate = JSON.parse(msg.data);
                        if (!candidate) {
                            return console.log("Failed to parse candidate");
                        }

                        peerConnection.addIceCandidate(candidate);                }
            }

            
        }

        startConnection();
    }, []);

    return (
        
        <div>
            <div id="localModal" className="d-flex row justify-content-center"> 
                <video id="localVideo" width="300" height="240" ref={mediaVideo} autoPlay muted></video>
                <div className="buttonWrapper d-flex column">
                    <button className="btn btn-primary" id="microphoneBtn" type="button">
                        <FontAwesomeIcon icon={faMicrophone} />
                    </button>
                    <button className="btn btn-primary" id="cameraBtn" type="button">
                        <FontAwesomeIcon icon={faVideo} />
                    </button>
                </div>
            </div>
            <div id="remoteVideos">
                {Object.keys(streams).forEach((key, index) => {
                    return <PeerVideo stream={streams[key]} />
                })}
            </div>
        </div>
    )
};



const PeerVideo = ({stream}) => {
    return (
        <video width={240} height={200} src={stream} autoPlay />
    )
}

export default ConferenceWsWrapper;