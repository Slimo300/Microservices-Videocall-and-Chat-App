import React, {useEffect, useState, useRef} from 'react';
import {  useParams, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faMicrophone, faVideo } from '@fortawesome/free-solid-svg-icons';

import useQuery from '../hooks/useQuery';
import "../Call.css";
import { GetWebRTCWebsocket } from '../requests/Ws';
import mockVideo from "../videos/mock.webm";


const VideoConference = () => {

    const { id } = useParams();
    const accessCode = useQuery().get("accessCode");
    const mocking = useQuery().get("mock");

    const peerConnection = useRef(new RTCPeerConnection());
    const ws = useRef(null);
    const localVideo = useRef(null);

    const [fatal, setFatal] = useState(false);

    const [RTCStreams, setRTCStreams] = useState({});
    const [userStream, setUserStream] = useState(null);

    const [audioState, setAudioState] = useState(true);
    const [audioTrack, setAudioTrack] = useState(null);
    const [audioSender, setAudioSender] = useState(null);

    const [videoState, setVideoState] = useState(true);
    const [videoTrack, setVideoTrack] = useState(null);
    const [videoSender, setVideoSender] = useState(null);

    const toggleAudio = () => {
        if (!audioState) {
            let audioSender = peerConnection.current.addTrack(audioTrack, userStream);
            setAudioSender(audioSender);
        } else {
            peerConnection.current.removeTrack(audioSender);
        }
        setAudioState(!audioState);
    };

    const toggleVideo = () => {
        if (!videoState) {
            let videoSender = peerConnection.current.addTrack(videoTrack, userStream);
            setVideoSender(videoSender);
        } else {
            peerConnection.current.removeTrack(videoSender);
        }
        setVideoState(!videoState);
    };

    useEffect(() => {
        const startStream = async () => {
            if (mocking) {
                const video = document.createElement("video");
                video.src = mockVideo;
                video.volume = 0.1;
                video.oncanplay = () => {
                    let stream = video.captureStream();
                    video.play();
                    setUserStream(stream);
                }
            } else {
                let stream = await navigator.mediaDevices.getUserMedia({video: true, audio: true});
                setUserStream(stream);
            }
        };
        
        if (!window.localStorage.getItem("token")) {
            setFatal(true);
            return;
        }
        startStream();
    }, [mocking]);

    useEffect(() => {
        if (!userStream) return;

        peerConnection.current.ontrack = (event) => {
            setRTCStreams(streams => {
                if (!streams[event.streams[0].id]) {
                    streams[event.streams[0].id] = event.streams[0];
                }
                return streams;
            });

            event.track.onmute = (event) => {
                // TODO: if video is muted display user.png
            }
            event.streams[0].onremovetrack = ({track}) => {
                if (!event.streams[0].active) {
                    setRTCStreams(streams => {
                        delete streams[event.streams[0].id];
                        return streams;
                    });
                }
                // TODO: if video is removed display user.png
            }
        };

        localVideo.current.srcObject = userStream;

        let newAudioTrack = userStream.getAudioTracks()[0];
        let newAudioSender = peerConnection.current.addTrack(newAudioTrack, userStream);
        console.log("track sent");
        setAudioTrack(newAudioTrack);
        setAudioSender(newAudioSender);

        let newVideoTrack = userStream.getVideoTracks()[0];
        let newVideoSender = peerConnection.current.addTrack(newVideoTrack, userStream);
        console.log("track sent");
        setVideoTrack(newVideoTrack);
        setVideoSender(newVideoSender);

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
                    console.log("received offer");
                    let offer = JSON.parse(msg.data);
                    if (!offer) {
                        return console.log("Failed to parse message");
                    }
                    peerConnection.current.setRemoteDescription(offer);
                    peerConnection.current.createAnswer().then(answer => {
                        peerConnection.current.setLocalDescription(answer);
                        ws.current.send(JSON.stringify({event: "answer", data: JSON.stringify(answer)}));
                        console.log("answer sent");
                    });
    
                    return;
                case "candidate":
                    console.log("received ICE candidate");
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

    }, [userStream, accessCode, id]);


    if (fatal) return <Navigate to="/not-found" />;

    return (
        <div>
            <div id="localModal" className="d-flex row justify-content-center"> 
                <video id="localVideo" width="300" height="240" ref={localVideo} autoPlay muted></video>
                <div className="buttonWrapper d-flex column">
                    <button className="btn btn-primary" id="microphoneBtn" type="button" onClick={toggleAudio}>
                        <FontAwesomeIcon icon={faMicrophone} />
                    </button>
                    <button className="btn btn-primary" id="cameraBtn" type="button" onClick={toggleVideo}>
                        <FontAwesomeIcon icon={faVideo} />
                    </button>
                </div>
            </div>
            <div id="remoteVideos">
                {Object.keys(RTCStreams).map(streamID => {
                    return <PeerVideo stream={RTCStreams[streamID]} />
                })}
            </div>
        </div>
    )
};

const PeerVideo = ({stream}) => {
    return (
        <video width={240} height={200} srcObject={stream} autoPlay />
    )
}


export default VideoConference;