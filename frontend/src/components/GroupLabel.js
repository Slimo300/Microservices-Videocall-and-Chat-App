import React, { useContext } from "react";
import { actionTypes, StorageContext } from "../ChatStorage";
import { GroupPicture } from "./Pictures";

export const GroupLabel = ({ group, setCurrent }) => {

    const [, dispatch] = useContext(StorageContext);
    const change = () => {
        dispatch({type: actionTypes.RESET_COUNTER, payload: group.ID});
        setCurrent(group);
    };
    
    return (
        <li className="person" onClick={change}>
            <div className="user">
                <GroupPicture pictureUrl={group.pictureUrl}/>
            </div>
            <p className="name-time">
                <span className="name">{group.name}</span>
            </p>
            {group.unreadMessages>0?<span className="badge badge-primary float-right">{group.unreadMessages}</span>:null}
        </li>
    );
}
