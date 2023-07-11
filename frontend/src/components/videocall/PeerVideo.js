import React, {useRef, useState, useEffect} from "react";

import mutedUser from "../../statics/images/video-user.png";

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

export default PeerVideo;