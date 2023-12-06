import React, { useRef, useEffect, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEllipsisVertical, faMicrophoneSlash } from "@fortawesome/free-solid-svg-icons";
import Switch from "react-switch";

import mutedUser from "../../statics/images/video-user.png";

export const USER_PEER_VIDEO_ELEMENT_ID = "userPeerVideo"

export const PeerVideo = ({stream, username, memberID, muting, ws}) => {
    const video = useRef(null);
    const [, setToggle] = useState(false);

    const createStreamFromImg = (audioTrack) => {
        let canvas = document.createElement("canvas");
        canvas.width = 400;
        canvas.height = 300;
        const ctx = canvas.getContext("2d");

        const image = new Image();
        image.src = mutedUser;
        image.onload = () => {
            ctx.drawImage(image, 0, 0, canvas.width, canvas.height);
            const videoTrack = canvas.captureStream().getVideoTracks()[0];
            if (audioTrack && audioTrack instanceof MediaStreamTrack) video.current.srcObject = new MediaStream([videoTrack, audioTrack])
            else video.current.srcObject = new MediaStream([videoTrack]);
        };
    }

    useEffect(() => {
        const setupRef = () => {

            if (!video.current) return;
            if (!(stream instanceof MediaStream)) { createStreamFromImg(); return; };
            if (stream.getVideoTracks().length === 0) { createStreamFromImg(stream.getAudioTracks()[0]); return; }

            video.current.srcObject = stream;
        };
        
        stream.onremovetrack = () => {
            setToggle(toggle => !toggle);
        };

        setupRef();
    });

    return (
        <div className='peer m-1'>
            <h4 className='white-text peer-footer'>{username?username:stream.id}</h4>
            <div className="dropdown peer-options">
                <button className="peer-options-button p-1" type="button" id="dropdownMenuButton" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                    <FontAwesomeIcon icon={faEllipsisVertical} size='m' />
                </button>
                <div className="bg-dark dropdown-menu" onClick={(e) => {e.stopPropagation()}} aria-labelledby="dropdownMenuButton">
                    <div className="dropdown-item muting-menu-text d-flex column justify-content-between">
                        <PrivateSwitch ws={ws} memberID={memberID} kind={"audio"} />
                        <div className="pl-2">Toggle Audio</div>
                    </div>
                    <div className="dropdown-item muting-menu-text d-flex column justify-content-between">
                        <PrivateSwitch ws={ws} memberID={memberID} kind={"video"} />
                        <div className="pl-2">Toggle Video</div>
                    </div>
                    {muting?<div className="dropdown-item muting-menu-text d-flex column justify-content-between">
                        <GlobalSwitch ws={ws} memberID={memberID} kind={"audio"} />
                        <div className="pl-2">Toggle Audio (for everyone)</div>
                    </div>:null}
                    {muting?<div className="dropdown-item muting-menu-text d-flex column justify-content-between">
                        <GlobalSwitch ws={ws} memberID={memberID} kind={"video"} />
                        <div className="pl-2">Toggle Video (for everyone)</div>
                    </div>:null}
                </div>
            </div>
            {!(stream instanceof MediaStream) || (stream instanceof MediaStream && stream.getAudioTracks().length === 0)?<h5 className="mute-symbol"><FontAwesomeIcon icon={faMicrophoneSlash} size='m' /></h5>:null}
            <video ref={video} className='peer-video' autoPlay />
        </div>
    )
}

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

    return (
        <div className='peer m-1' id={USER_PEER_VIDEO_ELEMENT_ID}>
            <h4 className='white-text peer-footer'>{username?username:peerStream.id}</h4>
            {peerStream instanceof MediaStream && peerStream.getAudioTracks().length === 0?<h5 className="mute-symbol p-1"><FontAwesomeIcon icon={faMicrophoneSlash} size='l' /></h5>:null}
            <video ref={video} className='peer-video' autoPlay muted={true} />
        </div>
    )
}

const ACTION_ENABLE = "enable";
const ACTION_DISABLE = "disable";

const PrivateSwitch = ({ ws, memberID, kind }) => {
    const [toggled, setToggled] = useState(false);

    const ToggleSwitch = () => {
        if (toggled) ws.current.send(JSON.stringify({event: "mute_for_yourself", data: JSON.stringify({actionType: ACTION_ENABLE, memberID, kind})}));
        else ws.current.send(JSON.stringify({event: "mute_for_yourself", data: JSON.stringify({actionType: ACTION_DISABLE, memberID, kind})}));
        setToggled(toggled => !toggled);
    }
 
    return (
        <div>
            <Switch onChange={ToggleSwitch} checked={toggled} />
        </div>
    )
};

const GlobalSwitch = ({ ws, memberID, kind }) => {  
    const [toggled, setToggled] = useState(false);
  
    useEffect(() => {
        document.getElementById(memberID+":"+kind).addEventListener("track_muted", ev => {
            console.log(ev.detail);
            switch(ev.detail.actionType) {
                case "disable":
                    setToggled(true);
                    break;
                case "enable":
                    setToggled(false);
                    break;
            }
        })
    })

    const ToggleSwitch = () => {
        if (toggled) ws.current.send(JSON.stringify({event: "mute_for_everyone", data: JSON.stringify({actionType: ACTION_ENABLE, memberID, kind})}));
        else ws.current.send(JSON.stringify({event: "mute_for_everyone", data: JSON.stringify({actionType: ACTION_DISABLE, memberID, kind})}));
    }
  
    return (
        <div id={memberID+":"+kind}>
            <Switch onChange={ToggleSwitch} checked={toggled} />
        </div>
    )
  }
  