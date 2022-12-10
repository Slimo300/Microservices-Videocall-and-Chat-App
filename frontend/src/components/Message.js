import React, { useContext } from "react";
import { UserPicture } from "./Pictures";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faEllipsis } from '@fortawesome/free-solid-svg-icons'
import {DeleteMessageForEveryone, DeleteMessageForYourself} from "../requests/Messages";
import { actionTypes, StorageContext } from "../ChatStorage";

const Message = (props) => {

    let time = new Date(props.message.created);
    let displayedTime = time.getHours() + ":" + time.getMinutes();
    let messageText = props.message.text;
    if (props.message.text === "") {
        messageText=<div style={{"font-style": "italic"}}>Message deleted</div>
    }

    const right = (
        <li className="chat-right">
            <div className="chat-hour">{displayedTime} <span className="fa fa-check-circle"></span></div>
            <div className="chat-text d-flex align-items-center">{messageText}</div>
            <div className="chat-avatar">
                <UserPicture pictureUrl={props.picture} />
                <div className="chat-name">{props.message.nick}</div>
            </div>
            <MessageOptions side="right" messageID={props.message.messageID} groupID={props.message.groupID}/>
        </li>
    );

    const left = (
        <li className="chat-left">
            <MessageOptions side="left" messageID={props.message.messageID} groupID={props.message.groupID}/>
            <div className="chat-avatar">
                <UserPicture pictureUrl={props.picture} />
                <div className="chat-name">{props.message.nick}</div>
            </div>
            <div className="chat-text d-flex align-items-center">{messageText}</div>
            <div className="chat-hour">{displayedTime} <span className="fa fa-check-circle"></span></div>
        </li>
    )

    return (
        <div>
            {props.message.userID===props.user?right:left}
        </div>
    )
}

const MessageOptions = (props) => {

    const [, dispatch] = useContext(StorageContext);

    const DeleteForYourself = async () => {
        let response;
        try {
            response = await DeleteMessageForYourself(props.groupID, props.messageID);
        } catch(err) {
            alert(err.response.data.err);
            return;
        }
        dispatch({type: actionTypes.DELETE_MESSAGE, payload: response.data});
    }

    const DeleteForEveryone = async () => {
    let response;
    try {
        response = await DeleteMessageForEveryone(props.groupID, props.messageID)
      } catch(err) {
        alert(err.response.data.err);
        return;
      }
      dispatch({type: actionTypes.DELETE_MESSAGE, payload: response.data});
    }

    return (
        <div className="btn-group dropup" style={{"height": "48px"}} >
            <button type="button" className="btn btn-outline-light" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                <FontAwesomeIcon icon={faEllipsis} />
            </button>
            <div className="dropdown-menu">
                <button type="button" class="btn btn-light dropdown-item" onClick={DeleteForYourself}><p style={{float: props.side}}>Delete for yourself</p></button>
                <button type="button" class="btn btn-light dropdown-item" onClick={DeleteForEveryone}><p style={{float: props.side}}>Delete for everyone</p></button>
            </div>
        </div>
    )
}

export default Message;