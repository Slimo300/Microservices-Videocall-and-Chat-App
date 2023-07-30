import React, { useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faMicrophone, faMicrophoneSlash, faVideo, faVideoSlash, faPhoneSlash, faShareFromSquare } from "@fortawesome/free-solid-svg-icons";

import PeerVideo from "./PeerVideo";
import { AUDIO_ACTIVE, AUDIO_INACTIVE, VIDEO_ACTIVE, VIDEO_INACTIVE, VIDEO_SCREENSHARE } from "../../pages/Call";

const CallScreen = ({CallHandler, userStream, videoState, audioState, RTCStreams}) => {

    const [ended, setEnded] = useState(false);

    const EndCall = () => {
        CallHandler.EndCall();

        setEnded(true);
    };

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

                <button className={"btn shadow rounded-circle "+(audioState===AUDIO_ACTIVE?"btn-secondary":"btn-danger")} type="button" onClick={CallHandler.ToggleAudio}>
                    <FontAwesomeIcon icon={audioState===AUDIO_ACTIVE?faMicrophone:faMicrophoneSlash} size='xl'/>
                </button>

                <button className={"btn shadow rounded-circle "+(videoState===VIDEO_ACTIVE?"btn-secondary":"btn-danger")} type="button" onClick={CallHandler.ToggleVideo}>
                    <FontAwesomeIcon icon={videoState===VIDEO_ACTIVE?faVideo:faVideoSlash} size='xl'/>
                </button>

                <button className={"btn shadow rounded-circle "+(videoState!==VIDEO_SCREENSHARE?"btn-secondary":"btn-danger")} type="button" onClick={CallHandler.ToggleScreenShare}>
                    <FontAwesomeIcon icon={faShareFromSquare} size='xl'/>
                </button>

                <button className="btn btn-danger shadow rounded-circle" type="button" onClick={EndCall}>
                    <FontAwesomeIcon icon={faPhoneSlash} size='xl' />
                </button>
            </div>
            <div className='d-flex flex-wrap justify-content-center align-items-center'>
                {userStream?<PeerVideo stream={userStream} username={localStorage.getItem("username")} isVideoMuted={videoState===VIDEO_INACTIVE} isAudioMuted={audioState===AUDIO_INACTIVE} isUser={true} />:null}
                {RTCStreams?RTCStreams.map(item => {
                    return <PeerVideo key={item.stream.id} stream={item.stream} username={item.username} isVideoMuted={!item.videoEnabled} isAudioMuted={!item.audioEnabled} isUser={false} />
                }):null}
            </div> 
        </div>
    );
};

export default CallScreen;