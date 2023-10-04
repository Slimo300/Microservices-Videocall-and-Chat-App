import React, { useContext } from "react";
import { UserPicture } from "../Pictures";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faEllipsis } from '@fortawesome/free-solid-svg-icons'
import {DeleteMessageForEveryone, DeleteMessageForYourself} from "../../requests/Messages";
import { actionTypes, StorageContext } from "../../ChatStorage";

const Message = ({ message, picture, user }) => {

    let time = new Date(message.created);
    let displayedTime = time.getHours() + ":" + (time.getMinutes()<10?'0':'') + time.getMinutes();

    const right = (
        <li className="chat-right">
            <div className="chat-hour">{displayedTime} <span className="fa fa-check-circle"></span></div>
            <MessageContent message={message} fileUrl={picture} side="right" />
            <div className="chat-avatar">
                <UserPicture pictureUrl={picture} />
                <div className="chat-name">{message.nick}</div>
            </div>
            <MessageOptions side="right" messageID={message.messageID} groupID={message.groupID}/>
        </li>
    );

    const left = (
        <li className="chat-left">
            <MessageOptions side="left" messageID={message.messageID} groupID={message.groupID}/>
            <div className="chat-avatar">
                <UserPicture pictureUrl={picture} />
                <div className="chat-name">{message.nick}</div>
            </div>
            <MessageContent message={message} fileUrl={picture} side="left" />
            <div className="chat-hour">{displayedTime} <span className="fa fa-check-circle"></span></div>
        </li>
    )

    return (
        <div>
            {message.userID===user?right:left}
        </div>
    )
}

const MessageContent = (props) => {
    let messageText = props.message.text;
    let isDeleted = props.message.text === "" && props.message.files.length === 0;
    let hasText = props.message.text !== "";
    if (isDeleted) {
        messageText=<div className="italic">Message deleted</div>
    }
    let messageHolderClassName = (props.side==="right")?"d-flex flex-row align-items-center justify-content-end":"d-flex flex-row align-items-center justify-content-start"

    return (
        <div className="d-flex flex-column justify-content-center">
            {hasText||isDeleted?<div className={messageHolderClassName}>
                <div className="chat-text d-flex justify-content-end">{messageText}</div>
            </div>:null}
            {props.message.files===undefined||props.message.files===null?null:<div className="d-flex flex-column">
            {props.message.files.map((item) => {
                return <MessageFile key={item.key} file={item} />
            })}
            </div>}
        </div>
    );
}

const MessageFile = (props) => {
    return (
    <img src={"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.file.key} style={{height: '200px', width: '200px', border: '1px solid'}} alt="sample"/>
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
                <button type="button" className="btn btn-light dropdown-item" onClick={DeleteForYourself}><p style={{float: props.side}}>Delete for yourself</p></button>
                <button type="button" className="btn btn-light dropdown-item" onClick={DeleteForEveryone}><p style={{float: props.side}}>Delete for everyone</p></button>
            </div>
        </div>
    )
}

export default Message;