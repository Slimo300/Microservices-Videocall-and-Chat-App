import React from "react";
import { UserPicture } from "../Pictures";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faEllipsis } from '@fortawesome/free-solid-svg-icons'
import {DeleteMessageForEveryone, DeleteMessageForYourself} from "../../requests/Messages";

const Message = ({ message, picture, user }) => {

    let time = new Date(message.created);
    let displayedTime = time.getHours() + ":" + (time.getMinutes()<10?'0':'') + time.getMinutes();
    let isDeleted = message.text === "" && message.files.length === 0;

    const right = (
        <li className="chat-right">
            <div className="chat-hour">{displayedTime} <span className="fa fa-check-circle"></span></div>
            <MessageContent message={message} side="right" />
            <div className="chat-avatar">
                <UserPicture pictureUrl={picture} />
                <div className="chat-name">{message.Member.username}</div>
            </div>
            <MessageOptions side="right" messageID={message.messageID} groupID={message.Member.groupID} isDeleted={isDeleted}/>
        </li>
    );

    const left = (
        <li className="chat-left">
            <MessageOptions side="left" messageID={message.messageID} groupID={message.Member.groupID} isDeleted={isDeleted}/>
            <div className="chat-avatar">
                <UserPicture pictureUrl={picture} />
                <div className="chat-name">{message.Member.username}</div>
            </div>
            <MessageContent message={message} side="left" />
            <div className="chat-hour">{displayedTime} <span className="fa fa-check-circle"></span></div>
        </li>
    )

    return (
        <div>
            {message.Member.userID===user?right:left}
        </div>
    )
}

const MessageContent = ({message, side}) => {
    let messageText = message.text;
    let isDeleted = message.text === "" && message.files.length === 0;
    let hasText = message.text !== "";
    if (isDeleted) {
        messageText=<div className="italic">Message deleted</div>
    }
    let messageHolderClassName = (side==="right")?"d-flex flex-row align-items-center justify-content-end":"d-flex flex-row align-items-center justify-content-start"

    return (
        <div className="d-flex flex-column justify-content-center">
            {hasText||isDeleted?<div className={messageHolderClassName}>
                <div className="chat-text d-flex justify-content-end">{messageText}</div>
            </div>:null}
            {message.files===undefined||message.files===null?null:<div className="d-flex flex-column">
            {message.files.map((item) => {
                return <MessageFile key={item.key} file={item} />
            })}
            </div>}
        </div>
    );
}

const MessageFile = (props) => {
    return (
    <img src={window._env_.STORAGE_URL+"/"+props.file.key} style={{height: '200px', width: '200px', border: '1px solid'}} alt="sample"/>
    )
}

const MessageOptions = ({groupID, messageID, canDelete, isDeleted, side}) => {

    // const [, dispatch] = useContext(StorageContext);

    const DeleteForYourself = async () => {
        try {
            if (!isDeleted) await DeleteMessageForYourself(groupID, messageID);
        } catch(err) {
            alert(err.response.data.err);
            return;
        }
    }

    const DeleteForEveryone = async () => {
    try {
        if (!isDeleted) await DeleteMessageForEveryone(groupID, messageID)
      } catch(err) {
        alert(err.response.data.err);
        return;
      }
    }

    return (
        <div className="btn-group dropup" style={{"height": "48px"}} >
            <button type="button" className="btn btn-outline-light" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                <FontAwesomeIcon icon={faEllipsis} />
            </button>
            <div className="dropdown-menu">
                <button type="button" className="btn btn-light dropdown-item" onClick={DeleteForYourself}><p style={{float: side}}>Delete for yourself</p></button>
                <button type="button" className="btn btn-light dropdown-item" onClick={DeleteForEveryone}><p style={{float: side}}>Delete for everyone</p></button>
            </div>
        </div>
    )
}

export default Message;