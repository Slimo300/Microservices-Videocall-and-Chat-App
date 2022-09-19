import React, { useContext } from "react";
import { actionTypes, StorageContext } from "../ChatStorage";

export const GroupLabel = (props) => {

    const [, dispatch] = useContext(StorageContext);
    const change = () => {
        dispatch({type: actionTypes.RESET_COUNTER, payload: props.group.ID});
        props.setCurrent(props.group);
    };
    
    return (
        <li className="person" onClick={change}>
            <div className="user">
                <img className="rounded-circle img-thumbnail"
                    src={"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.group.pictureUrl}
                    onError={({ currentTarget }) => {
                        currentTarget.onerror = null; 
                        currentTarget.src="https://cdn.icon-icons.com/icons2/3005/PNG/512/people_group_icon_188185.png";
                    }}
                />
            </div>
            <p className="name-time">
                <span className="name">{props.group.name}</span>
            </p>
            {props.group.unreadMessages>0?<span className="badge badge-primary float-right">{props.group.unreadMessages}</span>:null}
        </li>
    );
}
