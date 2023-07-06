import React, {useEffect, useState, useRef} from 'react';
import {  useParams, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faMicrophone, faVideo, faShareFromSquare, faPhoneSlash } from '@fortawesome/free-solid-svg-icons';

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

    const [fatal, setFatal] = useState(false);
    const [ended, setEnded] = useState(false);

    const [RTCStreams, setRTCStreams] = useState({});
    const [userStream, setUserStream] = useState(null);

    const audio = useRef({});
    const video = useRef({});

    const toggleAudio = () => {
        if (!audio.current.state) {
            console.log("adding track");
            audio.current.sender = peerConnection.current.addTrack(audio.current.track, userStream);
        } else {
            console.log("removing track");
            peerConnection.current.removeTrack(audio.current.sender);
        }
        audio.current.state = !audio.current.state;
    };
    const toggleVideo = () => {
        if (!video.current.state) {
            video.current.sender = peerConnection.current.addTrack(video.current.track, userStream);
        } else {
            peerConnection.current.removeTrack(video.current.sender);
        }
        video.current.state = !video.current.state;
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
            console.log("New track received: ", event.track);
            setRTCStreams(streams => {
                if (!streams[event.streams[0].id]) {
                    streams[event.streams[0].id] = event.streams[0];
                }
                return {...streams};
            });

            event.track.onmute = (event) => {
                // TODO: if video is muted display user.png
                console.log("Track muted");
            }
            event.streams[0].onremovetrack = ({track}) => {
                console.log("Track removed");
                if (!event.streams[0].active) {
                    setRTCStreams(streams => {
                        delete streams[event.streams[0].id];
                        return {...streams};
                    });
                }
                // TODO: if video is removed display user.png
            }
        };

        peerConnection.current.onnegotiationneeded = (event) => {
            console.log("negotiation needed: ", event);
        };

        audio.current.track = userStream.getAudioTracks()[0];
        audio.current.sender = peerConnection.current.addTrack(audio.current.track, userStream);
        audio.current.state = true;

        video.current.track = userStream.getVideoTracks()[0];
        video.current.sender = peerConnection.current.addTrack(video.current.track, userStream);
        video.current.state = true;

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

    const EndCall = () => {
        peerConnection.current.close()
        ws.current.close();
        userStream.getTracks().forEach((track) => {
            track.stop();
        })
        Object.keys(RTCStreams).forEach((key) => {
            RTCStreams[key].getTracks().forEach((track) => {
                track.stop();
            })
        });

        setEnded(true);
        setUserStream(null);
        setRTCStreams(null);
    };

    const MockPeer = () => {
        const newStream = userStream.clone();
        setRTCStreams(streams => {
            streams[newStream.id] = newStream;
            return {...streams};
        });
    };

    const PrintTracks = () => {
        let res = {};
        Object.keys(RTCStreams).forEach((key) => {
            let tracks = RTCStreams[key].getTracks();
            res[key] = tracks;
        })
        console.log(res);
    };

    if (fatal) return <Navigate to="/not-found" />;
    if (ended) return (
        <div className="container mt-4 pt-4">
        <div className="mt-5 d-flex justify-content-center">
          <div className="mt-5 row">
          <div className="display-1 mb-4 text-center text-primary">Call ended</div>
          </div>
        </div>
        </div>
    );

    return (
        <div>
            <div id="toolbar" className='d-flex justify-content-around rounded p-1'>
                <button className='btn btn-danger' onClick={PrintTracks}>Tracks</button>
                <button className='btn btn-danger' onClick={MockPeer}>Mock Peer</button>
                <button className="btn btn-secondary shadow rounded-circle" id="microphoneBtn" type="button" onClick={toggleAudio}>
                    <FontAwesomeIcon icon={faMicrophone} size='xl'/>
                </button>
                <button className="btn btn-secondary shadow rounded-circle" id="cameraBtn" type="button" onClick={toggleVideo}>
                    <FontAwesomeIcon icon={faVideo} size='xl'/>
                </button>
                <button className="btn btn-secondary shadow rounded-circle" id="shareScreenBtn" type="button" onClick={toggleVideo}>
                    <FontAwesomeIcon icon={faShareFromSquare} size='xl'/>
                </button>
                <button className="btn btn-danger shadow rounded-circle" id="endCallBtn" type="button" onClick={EndCall}>
                    <FontAwesomeIcon icon={faPhoneSlash} size='xl' />
                </button>
            </div>
            <div id="remoteVideos" className='d-flex flex-wrap justify-content-center align-items-center'>
                {userStream?<PeerVideo stream={userStream} muted={true} />:null}
                {Object.keys(RTCStreams).map(streamID => {
                    console.log(RTCStreams);
                    return <PeerVideo stream={RTCStreams[streamID]} />
                })}
            </div> 
        </div>
    )
};

const PeerVideo = ({stream, muted}) => {
    const video = useRef(null);
    useEffect(() => {
        video.current.srcObject = stream;
    });

    return (
        <div id={"stream-"+stream.id} className='peer m-1'>
            <p className='white-text peer-header'>{stream.id}</p>
            <video id={"media-"+stream.id} ref={video} className='peer-video' autoPlay muted={muted} />
        </div>
    )
}


export default VideoConference;