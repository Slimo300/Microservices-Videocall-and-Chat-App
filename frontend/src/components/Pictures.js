import React from "react";

import group from "../statics/images/group.jpg";
import user from "../statics/images/user.png";

export const UserPicture = ({ userID, hasPicture, imageID }) => {
    return (
        <img className="rounded-circle img-thumbnail" alt="user" id={imageID}
            src={hasPicture?window._env_.STORAGE_URL+"/"+userID:user}
            onError={({ currentTarget }) => {
                currentTarget.onerror = null; 
                currentTarget.src=user;
            }}
        />
    );
} 
export const GroupPicture = ({ groupID, hasPicture, imageID }) => {
    return (
        <img className="rounded-circle img-thumbnail" alt="group" id={imageID}
            src={hasPicture?window._env_.STORAGE_URL+"/"+groupID:group}
            onError={({ currentTarget }) => {
                currentTarget.onerror = null; 
                currentTarget.src=group;
            }}
        />
    );
} 