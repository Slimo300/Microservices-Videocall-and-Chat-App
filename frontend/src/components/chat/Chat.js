import React, {useEffect, useState } from "react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faVideo, faPhone } from '@fortawesome/free-solid-svg-icons'

import GroupMenu from "./GroupMenu";
import { ModalAddUser } from "../modals/AddUser";
import { ModalDeleteGroup } from "../modals/DeleteGroup";
import { ModalMembers } from "../modals/GroupMembers";
import { ModalGroupOptions } from "../modals/GroupOptions";
import { ModalLeaveGroup } from "../modals/LeaveGroup";
import { ChatBox } from "./Chat-Box";
import ChatInput from "./Chat-Input";
import { GetWebRTCAccessCode } from "../../requests/Ws";

const Chat = ({group, user, ws, setCurrent}) => {

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
                if (group.ID === undefined) {
                    return
                }
                for (let i = 0; i < group.Members.length; i++) {
                    if (group.Members[i].userID === user.ID ) {
                        setMember(group.Members[i]);
                        return;
                    }
                }
                throw new Error("No member matches user");
            }
        )();
    }, [group, user.ID]);

    // function for sending message when submit

    const JoinVideoCall = async () => {
        try {
            let accessCode = await GetWebRTCAccessCode(group.ID);
            window.open("call/"+group.ID+"?accessCode="+accessCode+"&initialVideo=true&initialAudio=true", "_blank", 'directories=no,titlebar=no,toolbar=no,location=no,status=no,menubar=no,scrollbars=no,resizable=no'); // 'directories=no,titlebar=no,toolbar=no,location=no,status=no,menubar=no,scrollbars=no,resizable=no'

        } catch(err) {
            alert(err);
        }
    }

    const JoinCall = async () => {
        try {
            let accessCode = await GetWebRTCAccessCode(group.ID);
            window.open("call/"+group.ID+"?accessCode="+accessCode+"&initialAudio=true", "_blank", 'directories=no,titlebar=no,toolbar=no,location=no,status=no,menubar=no,scrollbars=no,resizable=no'); // 'directories=no,titlebar=no,toolbar=no,location=no,status=no,menubar=no,scrollbars=no,resizable=no'

        } catch(err) {
            alert(err);
        }
    }

    let load;
    if (group.ID === undefined) {
        load = <h1 className="text-center">Select a group to chat!</h1>;
    } else {
        load = (
            <div className="col-xl-8 col-lg-8 col-md-8 col-sm-9 col-9">
                <div className="selected-user row">
                    <span className="mr-auto mt-4">To: <span className="name">{group.name}</span></span>
                    <button className="btn btn-primary mt-3 mr-1 mb-3" type="button" onClick={JoinCall}>
                        <FontAwesomeIcon icon={faPhone} />
                    </button>
                    <button className="btn btn-primary mt-3 mr-1 mb-3" type="button" onClick={JoinVideoCall}>
                        <FontAwesomeIcon icon={faVideo} />
                    </button>
                    <div className="dropdown">
                        <button className="btn btn-primary dropdown-toggle" type="button" id="dropdownMenuButton" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            Settings
                        </button>
                        <GroupMenu member={member} toggleOptions={toggleOptions} toggleDel={toggleDelGroup} toggleAdd={toggleAddUser} toggleMembers={toggleMembers} toggleLeave={toggleLeaveGroup}/>
                    </div>
                </div>
                <div className="chat-container d-flex flex-column justify-content-end" style={{'height': '80vh'}}>
                    <ChatBox group={group} user={user} />
                    <ChatInput ws={ws} group={group} user={user} member={member}/>
                </div>
                <ModalDeleteGroup show={delGrShow} toggle={toggleDelGroup} group={group} setCurrent={setCurrent}/>
                <ModalLeaveGroup show={leaveGrShow} toggle={toggleLeaveGroup} member={member} group={group} setCurrent={setCurrent}/>
                <ModalAddUser show={addUserShow} toggle={toggleAddUser} group={group}/>
                <ModalMembers show={membersShow} toggle={toggleMembers} group={group} member={member} />
                <ModalGroupOptions show={optionsShow} toggle={toggleOptions} group={group} />
            </div>
        );
    }
    return load;
}

export default Chat;
