import React, {useRef, useEffect} from "react";

import mutedUser from "../../statics/images/video-user.png";

const PeerVideo = ({stream, username, isUser, isVideoMuted}) => {
    const video = useRef(null);

    useEffect(() => {
        const setupRef = () => {
            if (!video.current || !(stream instanceof MediaStream)) return;

            if (stream.getVideoTracks().length === 0 || isVideoMuted) {
                let canvas = document.createElement("canvas");
                canvas.width = 400;
                canvas.height = 300;
                const ctx = canvas.getContext("2d");

                const image = new Image();
                image.src = mutedUser;
                image.onload = () => {
                    ctx.drawImage(image, 0, 0, canvas.width, canvas.height);

                    const videoTrack = canvas.captureStream().getVideoTracks()[0];

                    video.current.srcObject = new MediaStream([stream.getAudioTracks()[0], videoTrack]);
                };

                return;
            }

            video.current.srcObject = stream;
        };
        
        setupRef();
    });

    return (
        <div className='peer m-1'>
            <h4 className='white-text peer-footer'>{username?username:stream.id}</h4>
            <video ref={video} className='peer-video' autoPlay muted={isUser} />
        </div>
    )
}

export default PeerVideo;