import React, { useContext, useState } from "react";
import { StorageContext } from "../../ChatStorage";
import { LoadMessages } from "../../requests/Messages";
import { actionTypes } from "../../ChatStorage";
import Message from "./Message";

export const ChatBox = ({ group, user }) => {

    const [allMessagesFlag, setAllMessagesFlag] = useState(false);
    const [, dispatch] = useContext(StorageContext);

    const GetMemberPicture = (group, userID) => {
        for (let i = 0; i < group.Members.length; i++) {
            if (group.Members[i].User.ID === userID) {
                return group.Members[i].User.pictureUrl;
            }
        }
        return "";
    };
    const loadMessages = async() => {
        let messages = await LoadMessages(group.ID.toString(), group.messages.length);
        if (messages.status === 204) {
            setAllMessagesFlag(true);
            return;
        }
        dispatch({type: actionTypes.ADD_MESSAGES, payload: {messages: messages.data, groupID: group.ID}}); 
    };

    // Date of the last message in chat-box
    let lastMessageDate = new Date(0);
    // Helper to hold lastMessageDate to be displayed when it is changed by 
    // shouldDisplayDate to current message date
    let dateToDisplay = new Date(0);

    const shouldDisplayDate = (currentMessageDate, previousMessageDate) => {
        dateToDisplay = lastMessageDate;
        lastMessageDate = currentMessageDate;
        let currentMessageTime = currentMessageDate.getTime();
        let previousMessageTime = previousMessageDate.getTime();

        if (previousMessageTime === 0) return false;
        
        return (previousMessageTime - currentMessageTime) > 3600000;
    }

    return (
        <div className="d-flex flex-column col p-0 vh-60">
            {!allMessagesFlag?<div className="text-center align-top"><p className="text-primary" style={{cursor: "pointer"}} onClick={loadMessages}>Load more messages</p></div>:null}         
            <ul className="d-flex flex-column-reverse col p-0 overflower">
                {group.messages===undefined?null:group.messages.map((item) => {
                return <div key={item.ID} className="d-flex flex-column justify-content-end">
                        <Message message={item} user={user.ID} picture={GetMemberPicture(group, item.userID)} />
                        {shouldDisplayDate(new Date(item.created), lastMessageDate)?<NewDate time={dateToDisplay} />:null}
                    </div>})}
                {group.messages[group.messages.length-1]===undefined?null:<div className="d-flex flex-column justify-content-end">
                    <NewDate time={group.messages[group.messages.length-1].created} />
                </div>}
            </ul>
        </div>
    );
}

const NewDate = (props) => {
    let time = new Date(props.time)
    let displayedTime = time.getDate()+"."+(time.getMonth()+1)+"."+time.getFullYear()+" "+time.getHours()+":"+(time.getMinutes()<10?'0':'') + time.getMinutes();
    return (
        <div className="d-flex flex-column justify-content-center align-self-center text-secondary my-3">
            {displayedTime}
        </div>
    );
}