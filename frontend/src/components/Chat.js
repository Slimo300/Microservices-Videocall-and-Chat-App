import React, { useContext, useEffect, useRef, useState } from "react";
import { StorageContext, actionTypes } from "../ChatStorage";
import APICaller from "../Requests";
import GroupMenu from "./GroupMenu";
import Message from "./Message";
import { ModalAddUser } from "./modals/AddUser";
import { ModalDeleteGroup } from "./modals/DeleteGroup";
import { ModalMembers } from "./modals/GroupMembers";
import { ModalGroupOptions } from "./modals/GroupOptions";
import { ModalLeaveGroup } from "./modals/LeaveGroup";

const Chat = (props) => {

    const [state, dispatch] = useContext(StorageContext);

    const [member, setMember] = useState({});
    const [msg, setMsg] = useState(""); // currently typed message

    const scrollRef = useRef();
    useEffect( () => {
        scrollRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [props.toggler] );

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
    const [allMessagesFlag, setAllMessagesFlag] = useState(false);

    const GetMemberPicture = (group, member) => {
        for (let i = 0; i < group.Members.length; i++) {
            if (group.Members[i].ID === member) {
                return group.Members[i].User.pictureUrl;
            }
        }
        return "";
    }

    // getting group membership
    useEffect(()=>{
        (
            async () => {
                if (props.group.ID === undefined) {
                    return
                }
                for (let i = 0; i < props.group.Members.length; i++) {
                    if (props.group.Members[i].user_id === state.user.ID ) {
                        setMember(props.group.Members[i]);
                        return;
                    }
                }
                throw new Error("No member matches user");
            }
        )();
    }, [props.group.ID, state.user.ID, props.group.Members]);

    // function for sending message when submit
    const sendMessage = (e) => {
        e.preventDefault();
        if (msg.trim() === "") return false;
        if (props.ws !== undefined) props.ws.send(JSON.stringify({
            "group": props.group.ID,
            "member": member.ID,
            "text": msg,
            "nick": member.nick
        }));
        document.getElementById("text-area").value = "";
        document.getElementById("text-area").focus();
    }

    const loadMessages = async() => {
        let messages = await APICaller.LoadMessages(props.group.ID.toString(), props.group.messages.length);
        if (messages.status === 204) {
            setAllMessagesFlag(true);
            return;
        }
        dispatch({type: actionTypes.ADD_MESSAGES, payload: {messages: messages.data, group: props.group.ID}}); 
    };

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
                        <GroupMenu toggleOptions={toggleOptions} toggleDel={toggleDelGroup} toggleAdd={toggleAddUser} toggleMembers={toggleMembers} toggleLeave={toggleLeaveGroup} member={member}/>
                    </div>
                </div>
                <div className="chat-container">
                    <ul className="chat-box" style={{height: '70vh', overflow: 'scroll'}}>
                        {!allMessagesFlag?<li className="text-center"><a className="text-primary" style={{cursor: "pointer"}} onClick={loadMessages}>Load more messages</a></li>:null}
                        {props.group.messages===undefined?null:props.group.messages.map((item) => {
                        return <div ref={scrollRef}>
                                <Message key={item.ID} time={item.created} message={item.text} name={item.nick} member={item.member} user={member.ID} picture={GetMemberPicture(props.group, item.member)} />
                            </div>})}
                    </ul>
                    <form id="chatbox" className="form-group mt-3 mb-0 d-flex column justify-content-center" onSubmit={sendMessage}>
                        <textarea autoFocus  id="text-area" className="form-control mr-1" rows="3" placeholder="Type your message here..." onChange={(e)=>{setMsg(e.target.value)}}></textarea>
                        <input className="btn btn-primary" type="submit" value="Send" />
                    </form>
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
