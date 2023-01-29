import React, {useEffect, useState } from "react";
import GroupMenu from "./GroupMenu";
import { ModalAddUser } from "../modals/AddUser";
import { ModalDeleteGroup } from "../modals/DeleteGroup";
import { ModalMembers } from "../modals/GroupMembers";
import { ModalGroupOptions } from "../modals/GroupOptions";
import { ModalLeaveGroup } from "../modals/LeaveGroup";
import { ChatBox } from "./Chat-Box";
import ChatInput from "./Chat-Input";

const Chat = (props) => {

    const [member, setMember] = useState({});

    // add user to group modal
    const [addUserShow, setAddUserShow] = useState(false);
    const toggleAddUser = () => {
        setAddUserShow(!addUserShow);
    };
    // delete group modal
    const [delGrShow, setDelGrShow] = useState(false);
    const toggleDelGroup = () => {
        setDelGrShow(!delGrShow);
    };
    // members modal
    const [membersShow, setMembersShow] = useState(false);
    const toggleMembers = () => {
        setMembersShow(!membersShow);
    };
    const [leaveGrShow, setLeaveGroupShow] = useState(false);
    const toggleLeaveGroup = () => {
        setLeaveGroupShow(!leaveGrShow);
    };
    const [optionsShow, setOptionsShow] = useState(false);
    const toggleOptions = () => {
        setOptionsShow(!optionsShow);
    };

    // getting group membership
    useEffect(()=>{
        (
            async () => {
                if (props.group.ID === undefined) {
                    return
                }
                for (let i = 0; i < props.group.Members.length; i++) {
                    if (props.group.Members[i].userID === props.user.ID ) {
                        setMember(props.group.Members[i]);
                        return;
                    }
                }
                throw new Error("No member matches user");
            }
        )();
    }, [props.group, props.user.ID]);

    // function for sending message when submit

    let load;
    if (props.group.ID === undefined) {
        load = <h1 className="text-center">Select a group to chat!</h1>;
    } else {
        load = (
            <div className="col-xl-8 col-lg-8 col-md-8 col-sm-9 col-9">
                <div className="selected-user row">
                    <span className="mr-auto mt-4">To: <span className="name">{props.group.name}</span></span>
                    <div className="dropdown">
                        <button className="btn btn-primary dropdown-toggle" type="button" id="dropdownMenuButton" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            Settings
                        </button>
                        <GroupMenu member={member} toggleOptions={toggleOptions} toggleDel={toggleDelGroup} toggleAdd={toggleAddUser} toggleMembers={toggleMembers} toggleLeave={toggleLeaveGroup}/>
                    </div>
                </div>
                <div className="chat-container d-flex flex-column justify-content-end" style={{'height': '80vh'}}>
                    <ChatBox group={props.group} user={props.user} />
                    <ChatInput ws={props.ws} group={props.group} user={props.user}/>
                </div>
                <ModalDeleteGroup show={delGrShow} toggle={toggleDelGroup} group={props.group} setCurrent={props.setCurrent}/>
                <ModalLeaveGroup show={leaveGrShow} toggle={toggleLeaveGroup} member={member} group={props.group} setCurrent={props.setCurrent}/>
                <ModalAddUser show={addUserShow} toggle={toggleAddUser} group={props.group}/>
                <ModalMembers show={membersShow} toggle={toggleMembers} group={props.group} member={member} />
                <ModalGroupOptions show={optionsShow} toggle={toggleOptions} group={props.group} />
            </div>
        );
    }
    return load;
}

export default Chat;
