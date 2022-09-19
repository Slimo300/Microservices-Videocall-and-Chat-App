import React from "react";

const Message = (props) => {
    const right = (
        <li className="chat-right">
            <div className="chat-hour">{props.time} <span className="fa fa-check-circle"></span></div>
            <div className="chat-text">{props.message}</div>
            <div className="chat-avatar">
                <img className="rounded-circle img-thumbnail"
                    src={"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.picture}
                    onError={({ currentTarget }) => {
                        currentTarget.onerror = null; 
                        currentTarget.src="https://erasmuscoursescroatia.com/wp-content/uploads/2015/11/no-user.jpg";
                    }}
                />
                <div className="chat-name">{props.name}</div>
            </div>
        </li>
    );

    const left = (
        <li className="chat-left">
            <div className="chat-avatar">
                <img className="rounded-circle img-thumbnail"
                    src={"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.picture}
                    onError={({ currentTarget }) => {
                        currentTarget.onerror = null; 
                        currentTarget.src="https://erasmuscoursescroatia.com/wp-content/uploads/2015/11/no-user.jpg";
                    }}
                />
                <div className="chat-name">{props.name}</div>
            </div>
            <div className="chat-text">{props.message}</div>
            <div className="chat-hour">{props.time} <span className="fa fa-check-circle"></span></div>
        </li>
    )

    return (
        <div>
            {props.member===props.user?right:left}
        </div>
    )
}

export default Message;