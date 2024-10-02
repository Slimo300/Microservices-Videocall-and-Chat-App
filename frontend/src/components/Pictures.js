import React from "react";

import group from "../statics/images/group.jpg";
import user from "../statics/images/user.png";

export const UserPicture = ({ userID, hasPicture, imageID }) => {
    return (
        <img className="rounded-circle img-thumbnail" alt="user" id={userID}
            src={!hasPicture?user:window._env_.STORAGE_URL+"/"+userID}
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
            src={!hasPicture?group:window._env_.STORAGE_URL+"/"+groupID}
            onError={({ currentTarget }) => {
                currentTarget.onerror = null; 
                currentTarget.src=group;
            }}
        />
    );
} 