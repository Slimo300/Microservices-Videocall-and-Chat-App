import React, { useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faShareFromSquare } from "@fortawesome/free-solid-svg-icons";

import useQuery from "../../hooks/useQuery";
import StartCall from "./StartCall";

const ScreenShareButton = ({userStream, setUserStream, video}) => {

    const [active, setActive] = useState(false);
    const mocking = useQuery().get("mock");

    const ToggleScreenShare = async () => {
        let videoTrack;

        userStream.getVideoTracks()[0].stop();
        if (active) {
            videoTrack = (await StartCall(mocking)).getVideoTracks()[0];
        } else {
            videoTrack = (await navigator.mediaDevices.getDisplayMedia({video: true, audio: false})).getVideoTracks()[0];
            videoTrack.onended = async () => {
                console.log("screen share track ended");

                let newTrack = (await navigator.mediaDevices.getUserMedia({video: true, audio: false})).getVideoTracks()[0];
                
                video.current.sender.replaceTrack(newTrack);

                const stream = new MediaStream([newTrack]);
                setUserStream(stream);

                setActive(false);
            };
        }

        video.current.sender.replaceTrack(videoTrack);

        const stream = new MediaStream([videoTrack]);
        setUserStream(stream);

        setActive(!active);
    };

    return (
        <button className={"btn shadow rounded-circle "+(active?"btn-secondary":"btn-danger")} type="button" onClick={ToggleScreenShare}>
            <FontAwesomeIcon icon={faShareFromSquare} size='xl'/>
        </button>
    );
};

export default ScreenShareButton;