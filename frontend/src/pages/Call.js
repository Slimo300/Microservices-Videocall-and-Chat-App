import React, {useEffect, useState, useRef, useReducer, useCallback, useMemo} from 'react';
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

    const peerConnection = useRef(new RTCPeerConnection());
    const ws = useRef(null);
    const audioSender = useRef(null);
    const videoSender = useRef(null);

    const [audioState, setAudioState] = useState("");
    const [videoState, setVideoState] = useState("");
    const [videoPrevState, setVideoPrevState] = useState("");

    const [fatal, setFatal] = useState(false);

    const [RTCStreams, dispatch] = useReducer(RTCStreamsReducer, []);

    const [userStream, setUserStream] = useState(null);

    useEffect(() => {
        const startCall = async () => {
            const stream = await await navigator.mediaDevices.getUserMedia({video: initialVideo==="true", audio: initialAudio==="true"});

            setUserStream(stream);

            peerConnection.current.ontrack = (event) => {
                dispatch({type: actionTypes.NEW_STREAM, payload: event.streams[0]});
    
                event.streams[0].onremovetrack = () => {
                    dispatch({type: actionTypes.DELETE_STREAM, payload: event.streams[0].id});
                }
            };

            audioSender.current = peerConnection.current.addTrack(stream.getAudioTracks()[0], stream);
            setAudioState(AUDIO_ACTIVE);

            if (initialVideo === "true") {
                videoSender.current = peerConnection.current.addTrack(stream.getVideoTracks()[0], stream);
                setVideoState(VIDEO_ACTIVE);
            } else {
                setVideoState(VIDEO_INACTIVE);
            }
    
            try {
                ws.current = GetWebRTCWebsocket(id, accessCode, stream.id, initialVideo==="true", initialAudio===false);
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

                        console.log(userData);
                        dispatch({type: actionTypes.SET_USER_INFO, payload: userData});
                        break;
                    case "mute":
                        let data = JSON.parse(msg.data);
                        if (!data) {
                            return console.log("Failed to parse mute message")
                        }
                        dispatch({type: actionTypes.TOGGLE_MUTE, payload: data});
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

    }, [accessCode, id, initialVideo, initialAudio]);

    const ToggleAudio = useCallback(async () => {
        if (audioState !== AUDIO_ACTIVE) {
            const track = (await navigator.mediaDevices.getUserMedia({audio: true})).getAudioTracks()[0];

            setUserStream(stream => {
                stream.addTrack(track);
    
                if (!audioSender.current) {
                    audioSender.current = peerConnection.current.addTrack(track, stream);
                    ws.current.send(JSON.stringify({event: "renegotiate"}));
                } else {
                    audioSender.current.replaceTrack(track);
                }
                ws.current.send(JSON.stringify({event: "mute", data: JSON.stringify({audioEnabled: true})}));
                return stream;
            })

            setAudioState(AUDIO_ACTIVE);
        } else {
            setUserStream(stream => {
                const track = stream.getAudioTracks()[0];
                track.stop();
                stream.removeTrack(track);

                return stream;
            })
            audioSender.current.replaceTrack(null);
            ws.current.send(JSON.stringify({event: "mute", data: JSON.stringify({audioEnabled: false})}));

            setAudioState(AUDIO_INACTIVE);
        }
    }, [audioState]);

    const ToggleVideo = useCallback(async () => {
        if (videoState !== VIDEO_ACTIVE) {
            const track = (await navigator.mediaDevices.getUserMedia({video: true})).getVideoTracks()[0];

            setUserStream(stream => {
                if (videoState === VIDEO_SCREENSHARE) {
                    stream.getVideoTracks()[0].stop();
                    stream.removeTrack(stream.getVideoTracks()[0]);
                }
                
                stream.addTrack(track);
    
                if (!videoSender.current) {
                    videoSender.current = peerConnection.current.addTrack(track, stream);
                    ws.current.send(JSON.stringify({event: "renegotiate"}));
                } else {
                    videoSender.current.replaceTrack(track);
                }
                ws.current.send(JSON.stringify({event: "mute", data: JSON.stringify({videoEnabled: true})}));
                return stream;
            });
    
            setVideoState(VIDEO_ACTIVE);
        } else {
            setUserStream(stream => {
                const track = stream.getVideoTracks()[0];
                track.stop();
                stream.removeTrack(track);

                return stream;
            });
    
            videoSender.current.replaceTrack(null);
            ws.current.send(JSON.stringify({event: "mute", data: JSON.stringify({videoEnabled: false})}));

            setVideoState(VIDEO_INACTIVE);
        }
    }, [videoState]);

    const ToggleScreenShare = useCallback(async () => {
        if (videoState !== VIDEO_SCREENSHARE) {
            const track = (await navigator.mediaDevices.getDisplayMedia({video: true})).getVideoTracks()[0];

            setUserStream(stream => {
                if (videoState === VIDEO_ACTIVE) {
                    stream.getVideoTracks()[0].stop();
                    stream.removeTrack(stream.getVideoTracks()[0]);
                }
                stream.addTrack(track);
    
                if (!videoSender.current) {
                    videoSender.current = peerConnection.current.addTrack(track, stream);
                    ws.current.send(JSON.stringify({event: "renegotiate"}));
                } else {
                    videoSender.current.replaceTrack(track);
                }
                if (videoState === VIDEO_INACTIVE) {
                    ws.current.send(JSON.stringify({event: "mute", data: JSON.stringify({videoEnabled: true})}));
                }
                return stream;
            });
    
            setVideoState(videoState => {
                setVideoPrevState(videoState);

                return VIDEO_SCREENSHARE;
            });

        } else {
            let track = null;
            if (videoPrevState === VIDEO_ACTIVE) {
                track = (await navigator.mediaDevices.getUserMedia({video: true})).getVideoTracks()[0];
            }

            setUserStream(stream => {
                stream.getVideoTracks()[0].stop();
                stream.removeTrack(stream.getVideoTracks()[0]);

                if (videoPrevState === VIDEO_ACTIVE) {
                    stream.addTrack(track);
                }

                return stream;
            });
    
            videoSender.current.replaceTrack(track);
            
            ws.current.send(JSON.stringify({event: "mute", data: JSON.stringify({videoEnabled: true})}));
            setVideoState(videoPrevState);
        }
    }, [videoState, videoPrevState]);

    const EndCall = useCallback(() => {

        peerConnection.current.close();
        ws.current.close();

        dispatch({type: actionTypes.END_SESSION});

        setUserStream(stream => {
            stream.getTracks().forEach((track) => {
                track.stop();
            })

            return null;
        });
    }, []);

    const CallHandler = useMemo(() => {
        return {
            EndCall, ToggleAudio, ToggleVideo, ToggleScreenShare
        }
    }, [EndCall, ToggleAudio, ToggleVideo, ToggleScreenShare])

    if (fatal) return <Navigate to="/not-found" />;

    return (
        <CallScreen CallHandler={CallHandler} userStream={userStream} RTCStreams={RTCStreams} audioState={audioState} videoState={videoState}/>
    )
};

export default VideoConference;