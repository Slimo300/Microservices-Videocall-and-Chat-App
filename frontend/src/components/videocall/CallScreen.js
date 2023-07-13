import React, { useState, useEffect } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faMicrophone, faMicrophoneSlash, faVideo, faVideoSlash, faPhoneSlash } from "@fortawesome/free-solid-svg-icons";

import PeerVideo from "./PeerVideo";
import MediaButton from "./MediaButton";
import ScreenShareButton from "./ScreenShareButton";

const CallScreen = ({endSession, dataChannel, stream, video, audio, RTCStreams}) => {

    const [ended, setEnded] = useState(false);
    const [userStream, setUserStream] = useState(null);

    useEffect(() => {
        setUserStream(stream);
    }, [stream]);

    useEffect(() => {
        if (!userStream ) return;

        if (dataChannel) dataChannel.onmessage = e => {
            console.log("Message received: %v", e.data);
        };
    }, [dataChannel, userStream]);

    
    if (dataChannel) dataChannel.onmessage = e => {
        console.log("Message received: %v", e.data);
    }

    const EndCall = () => {
        endSession();
        userStream.getTracks().forEach((track) => {
            track.stop();
        });

        setEnded(true);
        setUserStream(null);
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
                <MediaButton isActive={true} mediaRef={audio} activeIcon={faMicrophone} inactiveIcon={faMicrophoneSlash} />
                <MediaButton isActive={true} mediaRef={video} activeIcon={faVideo} inactiveIcon={faVideoSlash} />
                
                <ScreenShareButton video={video} userStream={userStream} setUserStream={setUserStream}/>
                <button className="btn btn-danger shadow rounded-circle" type="button" onClick={EndCall}>
                    <FontAwesomeIcon icon={faPhoneSlash} size='xl' />
                </button>
            </div>
            <div className='d-flex flex-wrap justify-content-center align-items-center'>
                {userStream?<PeerVideo dataChannel={dataChannel} stream={userStream} isUser={true} />:null}
                {Object.keys(RTCStreams).map(streamID => {
                    return <PeerVideo dataChannel={dataChannel} key={streamID} stream={RTCStreams[streamID]} />
                })}
            </div> 
        </div>
    );
};

export default CallScreen;