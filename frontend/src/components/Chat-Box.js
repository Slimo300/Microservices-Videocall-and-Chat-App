import React from "react";
import Message from "./Message";

export const ChatBox = (props) => {

    const GetMemberPicture = (group, userID) => {
        for (let i = 0; i < group.Members.length; i++) {
            if (group.Members[i].User.ID === userID) {
                return group.Members[i].User.pictureUrl;
            }
        }
        return "";
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
        <ul className="chat-box d-flex flex-column-reverse" style={{'overflow-y': 'scroll', 'height': '65vh'}}>
            {props.group.messages===undefined?null:props.group.messages.map((item) => {
            return <div className="d-flex flex-column justify-content-end" ref={props.scrollRef}>
                    <Message key={item.ID} message={item} user={props.user.ID} picture={GetMemberPicture(props.group, item.userID)} />
                    {shouldDisplayDate(new Date(item.created), lastMessageDate)?<NewDate time={dateToDisplay} />:null}
                </div>})}
            {props.group.messages[props.group.messages.length-1]===undefined?null:<div className="d-flex flex-column justify-content-end" ref={props.scrollRef}>
                <NewDate time={props.group.messages[props.group.messages.length-1].created} />
            </div>}
        </ul>
    );
}

const NewDate = (props) => {
    let time = new Date(props.time)
    let displayedTime = time.getDate()+"."+(time.getMonth()+1)+"."+time.getFullYear()+" "+time.getHours()+":"+(time.getMinutes()<10?'0':'') + time.getMinutes();
    return (
        <div className="d-flex flex-column justify-content-center align-self-center text-secondary">
            {displayedTime}
        </div>
    );
}