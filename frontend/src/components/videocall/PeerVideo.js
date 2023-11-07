import React, { useRef, useEffect, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEllipsisVertical, faMicrophoneSlash } from "@fortawesome/free-solid-svg-icons";

import mutedUser from "../../statics/images/video-user.png";

export const USER_PEER_VIDEO_ELEMENT_ID = "userPeerVideo"

const PeerVideo = ({stream, username}) => {
    const video = useRef(null);

    useEffect(() => {
        const setupRef = () => {
            if (!video.current || !(stream instanceof MediaStream)) return;

            if (stream.getVideoTracks().length === 0) {
                let canvas = document.createElement("canvas");
                canvas.width = 400;
                canvas.height = 300;
                const ctx = canvas.getContext("2d");

                const image = new Image();
                image.src = mutedUser;
                image.onload = () => {
                    ctx.drawImage(image, 0, 0, canvas.width, canvas.height);

                    const videoTrack = canvas.captureStream().getVideoTracks()[0];
                    const audioTrack = stream.getAudioTracks()[0];

                    if (audioTrack) video.current.srcObject = new MediaStream([audioTrack, videoTrack]);
                    else video.current.srcObject = new MediaStream([videoTrack]);
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
            <div className="dropdown peer-options">
                <button className="peer-status-icon p-1" type="button" id="dropdownMenuButton" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                    <FontAwesomeIcon icon={faEllipsisVertical} size='m' />
                </button>
                <div className="dropdown-menu" aria-labelledby="dropdownMenuButton">
                    <button className="dropdown-item"> Option #1</button>
                    <button className="dropdown-item"> Option #2</button>
                    <div className="dropdown-divider"></div>
                    <button className="dropdown-item"> Option #1</button>
                    <button className="dropdown-item"> Option #2</button>
                </div>  
            </div>
            {stream instanceof MediaStream && stream.getAudioTracks().length === 0?<h5 className="mute-symbol"><FontAwesomeIcon icon={faMicrophoneSlash} size='m' /></h5>:null}
            <video ref={video} className='peer-video' autoPlay />
        </div>
    )
}

export default PeerVideo;

export const UserPeerVideo = ({ stream, username }) => {
    const video = useRef(null);

    const [peerStream, setPeerStream] = useState(null);

    useEffect(() => {
        setPeerStream(stream);

        document.getElementById(USER_PEER_VIDEO_ELEMENT_ID).addEventListener("streamchange", ev => setPeerStream(new MediaStream(ev.detail.getTracks())));
    }, [stream]);

    useEffect(() => {
        const setupRef = () => {
            if (!video.current || !(peerStream instanceof MediaStream)) return;

            if (peerStream.getVideoTracks().length === 0) {
                let canvas = document.createElement("canvas");
                canvas.width = 400;
                canvas.height = 300;
                const ctx = canvas.getContext("2d");

                const image = new Image();
                image.src = mutedUser;
                image.onload = () => {
                    ctx.drawImage(image, 0, 0, canvas.width, canvas.height);

                    const videoTrack = canvas.captureStream().getVideoTracks()[0];
                    const audioTrack = peerStream.getAudioTracks()[0];

                    if (audioTrack) video.current.srcObject = new MediaStream([audioTrack, videoTrack]);
                    else video.current.srcObject = new MediaStream([videoTrack]);
                };

                return;
            }

            video.current.srcObject = peerStream;
        };
        
        setupRef();
    });

    if (peerStream) console.log(peerStream instanceof MediaStream && peerStream.getAudioTracks().length === 0);

    return (
        <div className='peer m-1' id={USER_PEER_VIDEO_ELEMENT_ID}>
            <h4 className='white-text peer-footer'>{username?username:peerStream.id}</h4>
            {peerStream instanceof MediaStream && peerStream.getAudioTracks().length === 0?<h5 className="mute-symbol p-1"><FontAwesomeIcon icon={faMicrophoneSlash} size='l' /></h5>:null}
            <video ref={video} className='peer-video' autoPlay muted={true} />
        </div>
    )
}