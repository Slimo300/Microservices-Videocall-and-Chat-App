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

    return (
        <ul className="chat-box d-flex flex-column-reverse" style={{'overflow-y': 'scroll', 'height': '65vh'}}>
            {props.group.messages===undefined?null:props.group.messages.map((item) => {
            return <div className="d-flex flex-column justify-content-end" ref={props.scrollRef}>
                    <Message key={item.ID} message={item} user={props.user.ID} picture={GetMemberPicture(props.group, item.userID)} />
                </div>})}
        </ul>
    );
}