import React, { useEffect, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faMicrophone, faMicrophoneSlash, faVideo, faVideoSlash, faPhoneSlash, faShareFromSquare } from "@fortawesome/free-solid-svg-icons";

import { actionTypes } from "./RTCStreams";
import { USER_PEER_VIDEO_ELEMENT_ID } from "./PeerVideo";

const VIDEO_ACTIVE = "VideoActive";
const VIDEO_SCREENSHARE = "VideoScreenshare";
const VIDEO_INACTIVE = "VideoInactive";
const AUDIO_ACTIVE = "AudioActive";
const AUDIO_INACTIVE = "AudioInactive";

const Toolbar = ({userStream, audioSender, videoSender, peerConnection, ws, dispatch, setEnded}) => {
    const [audioState, setAudioState] = useState("");
    const [videoState, setVideoState] = useState("");

    useEffect(() => {
        if (audioSender.current && audioSender.current.track) setAudioState(AUDIO_ACTIVE);
        else setAudioState(AUDIO_INACTIVE);
        if (videoSender.current && videoSender.current.track) setVideoState(VIDEO_ACTIVE);
        else setVideoState(VIDEO_INACTIVE);
    }, [audioSender, videoSender]);

    const ToggleAudio = async () => {
        if (audioState !== AUDIO_ACTIVE) {
            const track = (await navigator.mediaDevices.getUserMedia({audio: true})).getAudioTracks()[0];

            userStream.current.addTrack(track);

            if (!audioSender.current) {
                audioSender.current = peerConnection.current.addTrack(track, userStream.current);
                ws.current.send(JSON.stringify({event: "renegotiate"}));
            } else {
                audioSender.current.replaceTrack(track);
                ws.current.send(JSON.stringify({event: "mute_yourself", data: JSON.stringify({actionType: "enable", kind: "audio"})}));
            }

            setAudioState(AUDIO_ACTIVE);
        } else {
            userStream.current.getAudioTracks()[0].stop();
            userStream.current.removeTrack(userStream.current.getAudioTracks()[0]);
            
            audioSender.current.replaceTrack(null);
            ws.current.send(JSON.stringify({event: "mute_yourself", data: JSON.stringify({actionType: "disable", kind: "audio"})}));

            setAudioState(AUDIO_INACTIVE);
        }

        document.getElementById(USER_PEER_VIDEO_ELEMENT_ID).dispatchEvent(new CustomEvent("streamchange", { detail: userStream.current }));
    };

    const ToggleVideo = async () => {
        if (videoState !== VIDEO_ACTIVE) {
            const track = (await navigator.mediaDevices.getUserMedia({video: true})).getVideoTracks()[0];

            if (videoState === VIDEO_SCREENSHARE) {
                userStream.current.getVideoTracks()[0].stop();
                userStream.current.removeTrack(userStream.current.getVideoTracks()[0]);
            }
            userStream.current.addTrack(track);

            if (!videoSender.current) {
                videoSender.current = peerConnection.current.addTrack(track, userStream.current);
                ws.current.send(JSON.stringify({event: "renegotiate"}));
            } else {
                videoSender.current.replaceTrack(track);
            }

            ws.current.send(JSON.stringify({event: "mute_yourself", data: JSON.stringify({actionType: "enable", kind: "video"})}));
            setVideoState(VIDEO_ACTIVE);

        } else {
            userStream.current.getVideoTracks()[0].stop();
            userStream.current.removeTrack(userStream.current.getVideoTracks()[0]);
    
            videoSender.current.replaceTrack(null);
            ws.current.send(JSON.stringify({event: "mute_yourself", data: JSON.stringify({actionType: "disable", kind: "video"})}));

            setVideoState(VIDEO_INACTIVE);
        }

        document.getElementById(USER_PEER_VIDEO_ELEMENT_ID).dispatchEvent(new CustomEvent("streamchange", { detail: userStream.current }));
    };

    const ToggleScreenShare = async () => {
        if (videoState !== VIDEO_SCREENSHARE) {
            const track = (await navigator.mediaDevices.getDisplayMedia({video: true})).getVideoTracks()[0];
            if (videoState === VIDEO_ACTIVE) {
                userStream.current.getVideoTracks()[0].stop();
                userStream.current.removeTrack(userStream.current.getVideoTracks()[0]);
            }
            userStream.current.addTrack(track);
    
            if (!videoSender.current) {
                videoSender.current = peerConnection.current.addTrack(track, userStream.current);
                ws.current.send(JSON.stringify({event: "renegotiate"}));
            } else {
                videoSender.current.replaceTrack(track);
            }
            if (videoState === VIDEO_INACTIVE) {
                ws.current.send(JSON.stringify({event: "mute_yourself", data: JSON.stringify({actionType: "enable", kind: "video"})}));
            }
            setVideoState(VIDEO_SCREENSHARE);

        } else {
            userStream.current.getVideoTracks()[0].stop();
            userStream.current.removeTrack(userStream.current.getVideoTracks()[0]);

            videoSender.current.replaceTrack(null);
            ws.current.send(JSON.stringify({event: "mute_yourself", data: JSON.stringify({actionType: "disable", kind: "video"})}));

            setVideoState(VIDEO_INACTIVE);
        }

        document.getElementById(USER_PEER_VIDEO_ELEMENT_ID).dispatchEvent(new CustomEvent("streamchange", { detail: userStream.current }));
    };


    const EndCall = () => {

        peerConnection.current.close();
        ws.current.close();

        dispatch({type: actionTypes.END_SESSION});

        userStream.current.getTracks().forEach((track) => {
            track.stop();
        });
        userStream.current = null;

        setEnded(true);
    };

    return (
        <div id="toolbar" className='d-flex justify-content-around rounded p-1'>

                <button className={"btn shadow rounded-circle "+(audioState===AUDIO_ACTIVE?"btn-secondary":"btn-danger")} type="button" onClick={ToggleAudio}>
                    <FontAwesomeIcon icon={audioState===AUDIO_ACTIVE?faMicrophone:faMicrophoneSlash} size='xl'/>
                </button>

                <button className={"btn shadow rounded-circle "+(videoState===VIDEO_ACTIVE?"btn-secondary":"btn-danger")} type="button" onClick={ToggleVideo}>
                    <FontAwesomeIcon icon={videoState===VIDEO_ACTIVE?faVideo:faVideoSlash} size='xl'/>
                </button>

                <button className={"btn shadow rounded-circle "+(videoState!==VIDEO_SCREENSHARE?"btn-secondary":"btn-danger")} type="button" onClick={ToggleScreenShare}>
                    <FontAwesomeIcon icon={faShareFromSquare} size='xl'/>
                </button>

                <button className="btn btn-danger shadow rounded-circle" type="button" onClick={EndCall}>
                    <FontAwesomeIcon icon={faPhoneSlash} size='xl' />
                </button>
            </div>
    )
}

export default Toolbar;