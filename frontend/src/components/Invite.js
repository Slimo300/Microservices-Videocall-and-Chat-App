import React, { useContext } from "react";
import { actionTypes, StorageContext } from "../ChatStorage";
import {RespondGroupInvite} from "../requests/Groups";
import { UserPicture, GroupPicture } from "./Pictures";

const Invite = (props) => {
    if (props.invite.status === 1){ 
        return (
            <AwaitingInvite {...props}/>  
        );
    }
    return (
        <AnsweredInvite {...props} />
    );
    
};

const AwaitingInvite = ({ invite, userID }) => {

    const [, dispatch] = useContext(StorageContext);

    const Respond = async (answer) => {
        let response;
        try {
            response = await RespondGroupInvite(invite.ID, answer);
            if (response.data.invite !== undefined) dispatch({type: actionTypes.UPDATE_INVITE, payload: response.data.invite});
            if (response.data.group !== undefined) dispatch({type: actionTypes.ADD_GROUP, payload: response.data.group});
            
        } catch (err) {
            if (err.response.data.err !== undefined) alert(err.response.data.err);
            else alert(err.message);
        }
    };

    let isUserATarget = false;
    if (userID === invite.targetID) {
        isUserATarget = true;
    }

    return (
        <div className="dropdown-item invite">
            <div className="list-group-item list-group-item-info d-flex row justify-content-around">
                {isUserATarget?<InviteImage pictureUrl={invite.issuer.pictureUrl} isUser={true}/>:null}
                <div className="chat-name align-self-center">{isUserATarget?invite.issuer.username:"You"}</div>
                <div className="align-self-center">invited </div>
                {isUserATarget?null:<InviteImage pictureUrl={invite.target.pictureUrl} isUser={true}/>}
                <div className="chat-name align-self-center">{isUserATarget?"You":invite.target.username}</div>
                <div className="align-self-center">to </div>
                <div className="chat-name align-self-center">{invite.group.name}</div>
                <InviteImage pictureUrl={invite.group.pictureUrl} isUser={false}/>
                {isUserATarget?<button className="btn-primary h-50 align-self-center" type="button" onClick={() => {Respond(true)}}>Accept</button>:null}
                {isUserATarget?<button className="btn-secondary h-50 align-self-center" type="button" onClick={() => {Respond(false)}}>Decline</button>:null}
            </div>
        </div>
    )
};

const AnsweredInvite = ({invite, userID}) => {

    let isUserATarget = false;
    if (userID === invite.targetID) {
        isUserATarget = true;
    }

    let action;
    if (invite.status === 2) {
        action = "accepted";
    } else if (invite.status === 3) {
        action = "rejected";
    } else {
        throw new Error("Message with wrong status");
    }

    return (
        <div className="dropdown-item invite">
            <div className="list-group-item list-group-item-info d-flex row justify-content-around">
                {isUserATarget?null:<InviteImage pictureUrl={invite.target.pictureUrl} isUser={true}/>}
                <div className="chat-name align-self-center">{isUserATarget?"You":invite.target.username}</div>
                <div className="align-self-center">{action} </div>
                {isUserATarget?<InviteImage pictureUrl={invite.issuer.pictureUrl} isUser={true}/>:null}
                <div className="chat-name align-self-center">{isUserATarget?invite.issuer.username:"your"}</div>
                <div className="align-self-center">invite to </div>
                <div className="chat-name align-self-center">{invite.group.name}</div>
                <InviteImage pictureUrl={invite.group.pictureUrl} isUser={false}/>
            </div>
        </div>
    )
};

const InviteImage = ({ pictureUrl, isUser}) => {
    let image = <GroupPicture pictureUrl={pictureUrl} />;
    if (isUser) {
        image = <UserPicture pictureUrl={pictureUrl} />
    }
    return (
        <div className="chat-avatar image-holder-invite">
            {image}
        </div>
    );
} 

export default Invite;