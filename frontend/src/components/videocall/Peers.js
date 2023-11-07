import React from "react";

import PeerVideo from "./PeerVideo";

const Peers = ({ userStream, RTCStreams }) => {
    return (
        <div className='d-flex flex-wrap justify-content-center align-items-center'>
            {userStream.current?<PeerVideo stream={userStream.current} username={localStorage.getItem("username")} isUser={true} />:null}
            {RTCStreams?RTCStreams.map(item => {
                return <PeerVideo key={item.stream.id} stream={item.stream} username={item.username} isUser={false} />
            }):null}
        </div> 
    )
}

export default Peers;