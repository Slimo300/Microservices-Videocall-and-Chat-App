import React, { useContext } from "react";
import { actionTypes, StorageContext } from "../ChatStorage";
import { GroupPicture } from "./Pictures";

export const GroupLabel = (props) => {

    const [, dispatch] = useContext(StorageContext);
    const change = () => {
        dispatch({type: actionTypes.RESET_COUNTER, payload: props.group.ID});
        props.setCurrent(props.group);
    };
    
    return (
        <li className="person" onClick={change}>
            <div className="user">
                <GroupPicture pictureUrl={props.group.pictureUrl}/>
            </div>
            <p className="name-time">
                <span className="name">{props.group.name}</span>
            </p>
            {props.group.unreadMessages>0?<span className="badge badge-primary float-right">{props.group.unreadMessages}</span>:null}
        </li>
    );
}
