import React, { useContext } from "react";
import { actionTypes, StorageContext } from "../ChatStorage";
import {RespondGroupInvite} from "../Requests";

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

const AwaitingInvite = (props) => {

    const [state, dispatch] = useContext(StorageContext);

    const Respond = async (answer) => {
        let response;
        try {
            response = await RespondGroupInvite(props.invite.ID, answer);
            console.log(response.data);
            if (response.data.invite !== undefined) dispatch({type: actionTypes.UPDATE_INVITE, payload: response.data.invite});
            if (response.data.group !== undefined) dispatch({type: actionTypes.ADD_GROUP, payload: response.data.group});
            
        } catch (err) {
            if (err.response.data.err !== undefined) alert(err.response.data.err);
            else alert(err.message);
        }
    };

    let isUserATarget = false;
    if (state.user.ID === props.invite.targetID) {
        isUserATarget = true;
    }

    return (
        <div className="dropdown-item invite">
            <div className="list-group-item list-group-item-info d-flex row justify-content-around">
                {isUserATarget?<InviteImage pictureUrl={props.invite.issuer.pictureUrl}/>:null}
                <div className="chat-name align-self-center">{isUserATarget?props.invite.issuer.username:"You"}</div>
                <div className="align-self-center">invited </div>
                {isUserATarget?null:<InviteImage pictureUrl={props.invite.target.pictureUrl}/>}
                <div className="chat-name align-self-center">{isUserATarget?"You":props.invite.target.username}</div>
                <div className="align-self-center">to </div>
                <div className="chat-name align-self-center">{props.invite.group.name}</div>
                <InviteImage pictureUrl={props.invite.group.pictureUrl}/>
                {isUserATarget?<button className="btn-primary h-50 align-self-center" type="button" onClick={() => {Respond(true)}}>Accept</button>:null}
                {isUserATarget?<button className="btn-secondary h-50 align-self-center" type="button" onClick={() => {Respond(false)}}>Decline</button>:null}
            </div>
        </div>
    )
};

const AnsweredInvite = (props) => {

    const [state,] = useContext(StorageContext);

    let isUserATarget = false;
    if (state.user.ID === props.invite.targetID) {
        isUserATarget = true;
    }

    let action;
    if (props.invite.status === 2) {
        action = "accepted";
    } else if (props.invite.status === 3) {
        action = "rejected";
    } else {
        throw new Error("Message with wrong status");
    }

    return (
        <div className="dropdown-item invite">
            <div className="list-group-item list-group-item-info d-flex row justify-content-around">
                {isUserATarget?null:<InviteImage pictureUrl={props.invite.target.pictureUrl}/>}
                <div className="chat-name align-self-center">{isUserATarget?"You":props.invite.target.username}</div>
                <div className="align-self-center">{action} </div>
                {isUserATarget?<InviteImage pictureUrl={props.invite.issuer.pictureUrl}/>:null}
                <div className="chat-name align-self-center">{isUserATarget?props.invite.issuer.username:"your"}</div>
                <div className="align-self-center">invite to </div>
                <div className="chat-name align-self-center">{props.invite.group.name}</div>
                <InviteImage pictureUrl={props.invite.group.pictureUrl}/>
            </div>
        </div>
    )
};

const InviteImage = (props) => {
    return (
        <div className="chat-avatar image-holder-invite">
            <img className="rounded-circle img-thumbnail"
                src={"https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.pictureUrl}
                onError={({ currentTarget }) => {
                    currentTarget.onerror = null; 
                    currentTarget.src="https://cdn.icon-icons.com/icons2/3005/PNG/512/people_group_icon_188185.png";
                }}
            />
        </div>
    );
} 

export default Invite;