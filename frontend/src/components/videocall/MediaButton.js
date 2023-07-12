import React, { useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const MediaButton = ({activeIcon, inactiveIcon, isActive, mediaRef, ws}) => {
    const [active, setActive] = useState(isActive);

    const ToggleButton = () => {
        if (!ws.current) {
            console.error("ws cannot be undefined");
        }

        if (!active) {
            mediaRef.current.sender.replaceTrack(mediaRef.current.track);
        } else {
            mediaRef.current.sender.replaceTrack(null);
        }
        
        setActive(!active);
    }

    if (!mediaRef || !activeIcon || !inactiveIcon ) throw new Error("function, or icons not specified"); 

    return (
        <button className={"btn shadow rounded-circle "+(active?"btn-secondary":"btn-danger")} type="button" onClick={ToggleButton}>
            <FontAwesomeIcon icon={active?activeIcon:inactiveIcon} size='xl'/>
        </button>
    );
}

export default MediaButton;
