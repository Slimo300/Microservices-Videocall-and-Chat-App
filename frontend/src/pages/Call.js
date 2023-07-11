import React, {useEffect, useState, useRef} from 'react';
import {  useParams, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faMicrophone, faMicrophoneSlash, faVideo, faVideoSlash, faShareFromSquare, faPhoneSlash } from '@fortawesome/free-solid-svg-icons';

import useQuery from '../hooks/useQuery';
import "../Call.css";
import { GetWebRTCWebsocket } from '../requests/Ws';
import mockVideo from "../statics/videos/mock.webm";
import mutedUser from "../statics/images/video-user.png";
import MediaButton from '../components/videocall/MediaButton';


const VideoConference = () => {

    const { id } = useParams();
    const accessCode = useQuery().get("accessCode");
    const mocking = useQuery().get("mock");

    const peerConnection = useRef(new RTCPeerConnection());
    const ws = useRef(null);
    const audio = useRef({});
    const video = useRef({});

    const [fatal, setFatal] = useState(false);
    const [ended, setEnded] = useState(false);

    const [RTCStreams, setRTCStreams] = useState({});
    const [userStream, setUserStream] = useState(null);

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

        audio.current.track = userStream.getAudioTracks()[0];
        audio.current.sender = peerConnection.current.addTrack(audio.current.track, userStream);

        video.current.track = userStream.getVideoTracks()[0];
        video.current.sender = peerConnection.current.addTrack(video.current.track, userStream);
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

    const ToggleScreenShare = async () => {
        let track;
        if (video.current.screenshare) {
            track = userStream.getVideoTracks()[0];
        } else {
            let stream = await navigator.mediaDevices.getDisplayMedia({video: true, audio: false})
            track = stream.getTracks()[0];
            track.onended = () => {
                console.log("screen share track ended");
                let track = userStream.getVideoTracks()[0];
                video.current.sender.replaceTrack(track);
                video.current.screenshare = false;
            }
        }
        video.current.sender.replaceTrack(track);
        video.current.screenshare = !video.current.screenshare;
    }

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

    const ShowTracks = ()=> {
        Object.keys(RTCStreams).forEach(key => {
            RTCStreams[key].getTracks().forEach(track => {
                console.log(track, track.muted);
            })
        })
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
                <button className='btn btn-danger' onClick={ShowTracks}>Show Tracks</button>
                <MediaButton isActive={true} mediaRef={audio} ws={ws} activeIcon={faMicrophone} inactiveIcon={faMicrophoneSlash} />
                <MediaButton isActive={true} mediaRef={video} ws={ws} activeIcon={faVideo} inactiveIcon={faVideoSlash} />
                <button className="btn btn-secondary shadow rounded-circle" type="button" onClick={ToggleScreenShare}>
                    <FontAwesomeIcon icon={faShareFromSquare} size='xl'/>
                </button>
                <button className="btn btn-danger shadow rounded-circle" type="button" onClick={EndCall}>
                    <FontAwesomeIcon icon={faPhoneSlash} size='xl' />
                </button>
            </div>
            <div id="remoteVideos" className='d-flex flex-wrap justify-content-center align-items-center'>
                {userStream?<PeerVideo stream={userStream} isUser={true} />:null}
                {Object.keys(RTCStreams).map(streamID => {
                    return <PeerVideo key={streamID} stream={RTCStreams[streamID]} />
                })}
            </div> 
        </div>
    )
};

const PeerVideo = ({stream, isUser}) => {
    const video = useRef(null);
    const [muted, setMuted] = useState(false);

    useEffect(() => {
        if (video.current) video.current.srcObject = stream;
    });

    return (
        <div className='peer m-1'>
            <p className='white-text peer-header'>{stream.id}</p>
            {muted?<img className='peer-video' src={mutedUser} alt='user without video'/>:<video ref={video} className='peer-video' autoPlay muted={isUser} />}
        </div>
    )
}


export default VideoConference;