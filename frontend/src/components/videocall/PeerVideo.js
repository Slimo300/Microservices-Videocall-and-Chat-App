import React, {useRef, useEffect} from "react";

import mutedUser from "../../statics/images/video-user.png";

const PeerVideo = ({stream, isUser, username, isVideoMuted}) => {
    const video = useRef(null);

    useEffect(() => {
        if (video.current) video.current.srcObject = stream;
    });

    return (
        <div className='peer m-1'>
            <p className='white-text peer-header'>{username?username:stream.id}</p>
            {isVideoMuted?<img className='peer-video' src={mutedUser} alt='user without video'/>:<video ref={video} className='peer-video' autoPlay muted={isUser} />}
        </div>
    )
}

export default PeerVideo;