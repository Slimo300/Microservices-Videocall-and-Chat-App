import React from "react";

import group from "../statics/images/group.jpg";
import user from "../statics/images/user.png";

export const UserPicture = ({ pictureUrl, imageID }) => {
    return (
        <img className="rounded-circle img-thumbnail" alt="user" id={imageID}
            src={(pictureUrl === "")?user:window._env_.STORAGE_URL+"/"+pictureUrl}
            onError={({ currentTarget }) => {
                currentTarget.onerror = null; 
                currentTarget.src=user;
            }}
        />
    );
} 
export const GroupPicture = ({ pictureUrl, imageID }) => {
    return (
        <img className="rounded-circle img-thumbnail" alt="group" id={imageID}
            src={(pictureUrl === "")?group:window._env_.STORAGE_URL+"/"+pictureUrl}
            onError={({ currentTarget }) => {
                currentTarget.onerror = null; 
                currentTarget.src=group;
            }}
        />
    );
} 