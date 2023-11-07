import React, { useState } from "react";

import { PeerVideo, UserPeerVideo } from "./PeerVideo";
import Toolbar from "./Toolbar";

const CallScreen = ({peerConnection, userStream, videoSender, audioSender, ws, RTCStreams, dispatch, username, muting}) => {

    const [ended, setEnded] = useState(false);

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
            <Toolbar userStream={userStream} audioSender={audioSender} videoSender={videoSender} peerConnection={peerConnection} ws={ws} dispatch={dispatch} setEnded={setEnded} />
            <div className='d-flex flex-wrap justify-content-center align-items-center'>
                {userStream.current?<UserPeerVideo stream={userStream.current} username={username}/>:null}
                {RTCStreams?RTCStreams.map(item => {
                    return <PeerVideo key={item.stream.id} stream={item.stream} username={item.username} memberID={item.memberID} muting={muting} ws={ws}/>
                }):null}
            </div> 
        </div>
    );
};

export default CallScreen;