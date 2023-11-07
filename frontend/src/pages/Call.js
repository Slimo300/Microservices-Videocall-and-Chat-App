import React, { useEffect, useState, useRef, useReducer } from 'react';
import {  useParams, Navigate } from "react-router-dom";

import useQuery from '../hooks/useQuery';
import { GetWebRTCWebsocket } from '../requests/Ws';
import CallScreen from '../components/videocall/CallScreen';
import { RTCStreamsReducer, actionTypes } from '../components/videocall/RTCStreams';

export const VIDEO_ACTIVE = "VideoActive";
export const VIDEO_SCREENSHARE = "VideoScreenshare";
export const VIDEO_INACTIVE = "VideoInactive";
export const AUDIO_ACTIVE = "AudioActive";
export const AUDIO_INACTIVE = "AudioInactive";

const VideoConference = () => {

    const { id } = useParams();
    const query = useQuery();
    const accessCode = query.get("accessCode");
    const initialVideo = query.get("initialVideo");
    const initialAudio = query.get("initialAudio");

    const userStream = useRef(null);

    const peerConnection = useRef(null);
    
    const ws = useRef(null);
    const audioSender = useRef(null);
    const videoSender = useRef(null);

    const [, setState] = useState(false);
    const [init, setInit] = useState(false);
    const [fatal, setFatal] = useState(false);

    const [RTCStreams, dispatch] = useReducer(RTCStreamsReducer, []);

    useEffect(() => {
        const startCall = async () => {
            userStream.current = await navigator.mediaDevices.getUserMedia({video: initialVideo==="true", audio: initialAudio==="true"});

            peerConnection.current = new RTCPeerConnection({'iceServers': [
                {
                    urls: `stun:${window._env_.TURN_ADDRESS}:${window._env_.TURN_PORT}`
                },
                {
                    urls: `turn:${window._env_.TURN_ADDRESS}:${window._env_.TURN_PORT}`,
                    username: window._env_.TURN_USER,
                    credential: window._env_.TURN_PASSWORD
                },
                {
                    urls: `turns:${window._env_.TURN_ADDRESS}:${window._env_.TURN_TLS_PORT}`,
                    username: window._env_.TURN_USER,
                    credential: window._env_.TURN_PASSWORD
                },
                {
                    urls: `turns:${window._env_.TURN_ADDRESS}:${window._env_.TURN_TLS_PORT}?transport=tcp`,
                    username: window._env_.TURN_USER,
                    credential: window._env_.TURN_PASSWORD
                }
            ]});

            peerConnection.current.ontrack = (event) => {
                dispatch({type: actionTypes.NEW_STREAM, payload: event.streams[0]});

                event.streams[0].onremovetrack = () => {
                    setState(state => { return !state });
                }
            };

            if (initialAudio === "true") {
                audioSender.current = peerConnection.current.addTrack(userStream.current.getAudioTracks()[0], userStream.current);
                // setAudioState(AUDIO_ACTIVE);
            } else {
                // setAudioState(AUDIO_INACTIVE);
            }

            if (initialVideo === "true") {
                videoSender.current = peerConnection.current.addTrack(userStream.current.getVideoTracks()[0], userStream.current);
                // setVideoState(VIDEO_ACTIVE);
            } else {
                // setVideoState(VIDEO_INACTIVE);
            }
    
            try {
                ws.current = GetWebRTCWebsocket(id, accessCode, userStream.current.id);
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
                    case "user_info":
                        let userData = JSON.parse(msg.data);
                        if (!userData) {
                            return console.log("Failed to parse newUser message");
                        }

                        dispatch({type: actionTypes.SET_USER_INFO, payload: userData});
                        break;
                    case "disconnected":
                        dispatch({type: actionTypes.USER_DISCONNECTED, payload: msg.data});
                        break;
                    default:
                        console.log("Unexpected websocket event: ", msg.event);
                }
            };

            setInit(true);

        };
        
        if (!window.localStorage.getItem("token")) {
            setFatal(true);
            return;
        }

        startCall();

    }, [accessCode, id, initialVideo, initialAudio]);

    if (fatal) return <Navigate to="/not-found" />;

    if (init) return (
        <CallScreen peerConnection={peerConnection} ws={ws} audioSender={audioSender} videoSender={videoSender} dispatch={dispatch} userStream={userStream} RTCStreams={RTCStreams} />
    )
    else return null;
};

export default VideoConference;