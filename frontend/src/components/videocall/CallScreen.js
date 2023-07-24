import React, { useState, useEffect } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faMicrophone, faMicrophoneSlash, faVideo, faVideoSlash, faPhoneSlash } from "@fortawesome/free-solid-svg-icons";

import PeerVideo from "./PeerVideo";
import MediaButton from "./MediaButton";
import ScreenShareButton from "./ScreenShareButton";
// import { actionTypes } from "./RTCStreams";

const CallScreen = ({CallHandler, peerConnection, ws, dispatch, stream, video, audio, RTCStreams}) => {

    const [ended, setEnded] = useState(false);
    const [userStream, setUserStream] = useState(null);

    useEffect(() => {
        setUserStream(stream);
    }, [stream]);

    const EndCall = () => {
        // peerConnection.current.close();
        // ws.current.close();

        // dispatch({type: actionTypes.END_SESSION});
        // userStream.getTracks().forEach((track) => {
        //     track.stop();
        // });

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
                <MediaButton isActive={true} mediaRef={audio} activeIcon={faMicrophone} inactiveIcon={faMicrophoneSlash} />
                <MediaButton isActive={true} mediaRef={video} activeIcon={faVideo} inactiveIcon={faVideoSlash} />
                
                <ScreenShareButton video={video} userStream={userStream} setUserStream={setUserStream}/>
                <button className="btn btn-danger shadow rounded-circle" type="button" onClick={EndCall}>
                    <FontAwesomeIcon icon={faPhoneSlash} size='xl' />
                </button>
            </div>
            <div className='d-flex flex-wrap justify-content-center align-items-center'>
                {userStream?<PeerVideo stream={userStream} isUser={true} username={localStorage.getItem("username")} isVideoMuted={false} />:null}
                {RTCStreams?RTCStreams.map(item => {
                    return <PeerVideo key={item.stream.id} stream={item.stream} username={item.username} isVideoMuted={false} />
                }):null}
            </div> 
        </div>
    );
};

export default CallScreen;