import React, { useContext, useEffect, useState } from "react";
import { UserPicture } from "../Pictures";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faEllipsis } from '@fortawesome/free-solid-svg-icons';
import { DeleteMessageForEveryone, DeleteMessageForYourself, GetPresignedGetRequests } from "../../requests/Messages";
import { actionTypes, StorageContext } from '../../ChatStorage';

const Message = ({ message, userID }) => {
// const Message = ({ message, picture, user }) => {

    let time = new Date(message.created);
    let displayedTime = time.getHours() + ":" + (time.getMinutes()<10?'0':'') + time.getMinutes();
    let isDeleted = message.text === "" && message.files.length === 0;

    const right = (
        <li className="chat-right">
            <div className="chat-hour">{displayedTime} <span className="fa fa-check-circle"></span></div>
            <MessageContent message={message} side="right" />
            <div className="chat-avatar">
                <UserPicture pictureUrl={message.Member.userID} />
                <div className="chat-name">{message.Member.username}</div>
            </div>
            <MessageOptions side="right" messageID={message.messageID} groupID={message.Member.groupID} isDeleted={isDeleted}/>
        </li>
    );

    const left = (
        <li className="chat-left">
            <MessageOptions side="left" messageID={message.messageID} groupID={message.Member.groupID} isDeleted={isDeleted}/>
            <div className="chat-avatar">
                <UserPicture pictureUrl={message.Member.userID} />
                <div className="chat-name">{message.Member.username}</div>
            </div>
            <MessageContent message={message} side="left" />
            <div className="chat-hour">{displayedTime} <span className="fa fa-check-circle"></span></div>
        </li>
    )

    return (
        <div>
            {message.Member.userID===userID?right:left}
        </div>
    )
}

const MessageContent = ({message, side}) => {

    const [fileUrls, setFileUrls] = useState(null);

    useEffect(() => {
        const fetchUrls = async () => {
            try {
                let result = await GetPresignedGetRequests(message.Member.groupID, message.files.map(file => {
                    return file.key;
                }))
    
                setFileUrls(result.data);
            } catch(err) {
                console.log(err.response.data);
            }
        }

        if (message.files && message.files.length > 0) fetchUrls();
    }, []);

    let messageText = message.text;
    let isDeleted = message.text === "" && message.files.length === 0;
    let hasText = message.text !== "";
    if (isDeleted) messageText=<div className="italic">Message deleted</div>

    return (
        <div className="d-flex flex-column justify-content-center">
            {hasText||isDeleted?<div className={"d-flex align-items-center justify-content-"+(side==="right")?"end":"start"}>
                <div className="chat-text d-flex justify-content-end">{messageText}</div>
            </div>:null}
            {message.files?<div className="d-flex flex-column">
                {fileUrls?fileUrls.map((item) => {
                    return <MessageFile key={item.key} url={item.url} />
                }):null}
            </div>:null}
        </div>
    );
}

const MessageFile = ({ url }) => {
    return (
        <img src={url} style={{height: '200px', width: '200px', border: '1px solid'}} alt="sample"/>
    )
}

const MessageOptions = ({messageID, isDeleted, side}) => {

    const [, dispatch] = useContext(StorageContext);

    const DeleteForYourself = async () => {
        try {
            if (!isDeleted) {
               let response = await DeleteMessageForYourself(messageID);
               dispatch({type: actionTypes.DELETE_MESSAGE, payload: { messageID: response.data.messageID, groupID: response.data.Member.groupID }});
            }
        } catch(err) {
            alert(err.response.data.err);
            return;
        }
    }

    const DeleteForEveryone = async () => {
    try {
        if (!isDeleted) await DeleteMessageForEveryone(messageID)
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