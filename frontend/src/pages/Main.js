import React, { useContext, useEffect, useState } from "react";
import {Navigate} from "react-router-dom";
import { actionTypes, StorageContext } from "../ChatStorage";
import Chat from "../components/Chat";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { GroupLabel } from "../components/GroupLabel";
import { ModalCreateGroup } from "../components/modals/CreateGroup";
import {GetGroups, GetInvites} from "../requests/Groups";
import { GetUser } from "../requests/Users";
import { GetWebsocket } from "../requests/Ws";
import { LoadMessages } from "../requests/Messages";
import { faPlus } from "@fortawesome/free-solid-svg-icons";
import { ModalUserProfile } from "./Profile";

const Main = (props) => {

    return (
        <div>
            {window.localStorage.getItem("token") === null? <Navigate to="/login" />:<AuthMain {...props}/>}
        </div>
    );
}

const AuthMain = (props) => {

    const [state, dispatch] = useContext(StorageContext);
    const [current, setCurrent] = useState({}); // current group
    const [toggler, setToggler] = useState(false);
    function toggleToggler(){
        setToggler(!toggler);
    }

    const [createGrShow, setCreateGrShow] = useState(false);
    const toggleCreateGroup = () => {
        setCreateGrShow(!createGrShow);
    }

    // Getting user data, groups and invites and setting websocket connection
    useEffect(() => {
        const fetchData = async () => {
            const userResult = await GetUser();
            dispatch({type: actionTypes.LOGIN, payload: userResult.data});
            const groupsResult = await GetGroups();
            if (groupsResult.status === 200) {
                dispatch({type: actionTypes.ADD_GROUPS, payload: groupsResult.data});
            }
            const invitesResult = await GetInvites(state.invites.length);
            if (invitesResult.status === 200) {
                dispatch({type: actionTypes.ADD_INVITES, payload: invitesResult.data});
            }
            let websocket = await GetWebsocket();
            props.setWs(websocket);
        };

        fetchData();
    }, [dispatch]);

    if (props.ws !== undefined) props.ws.onmessage = (e) => {
        const msgJSON = JSON.parse(e.data);
        console.log(msgJSON);
        if (msgJSON.type !== undefined) {
            switch (msgJSON.type) {
                case "DELETE_GROUP":
                    dispatch({type: actionTypes.DELETE_GROUP, payload: msgJSON.payload});
                    break;
                case "UPDATE_MEMBER":
                    console.log("UPDATE_MEMBER");
                    dispatch({type: actionTypes.UPDATE_MEMBER, payload: msgJSON.payload});
                    break;
                case "DELETE_MEMBER":
                    dispatch({type: actionTypes.DELETE_MEMBER, payload: msgJSON.payload});
                    break;
                case "ADD_MEMBER":
                    dispatch({type: actionTypes.ADD_MEMBER, payload: msgJSON.payload});
                    break;
                case "ADD_INVITE":
                    dispatch({type: actionTypes.ADD_INVITE, payload: msgJSON.payload});
                    break;
                case "UPDATE_INVITE":
                    dispatch({type: actionTypes.UPDATE_INVITE, payload: msgJSON.payload});
                    break;
                case "DELETE_MESSAGE":
                    dispatch({type: actionTypes.DELETE_MESSAGE, payload: msgJSON.payload});
                    break;
                default:
                    console.log("Unexpected action from websocket: ", msgJSON.type);
            }
            return;
        }
        if (msgJSON.groupID === current.ID) { // add message to state
            console.log("current: ", current.ID);
            dispatch({type: actionTypes.ADD_MESSAGE, payload: {message: msgJSON, current: true}})
            toggleToggler();
        } else {
            console.log("Not current: ", current.ID);
            dispatch({type: actionTypes.ADD_MESSAGE, payload: {message: msgJSON, current: false}})
        }
    }

    // getting messages from specific group
    useEffect(()=>{
        (
            async () => {
                if (current.ID !== undefined && current.messages.length === 0) {
                    let messages = await LoadMessages(current.ID.toString(), 0);
                    if (messages.status === 204) {
                        return;
                    }
                    dispatch({type: actionTypes.ADD_MESSAGES, payload: {messages: messages.data, groupID: current.ID}})
                    toggleToggler();
                }
            }
        )();
    }, [current, dispatch]);

    return (
        <div className="container" >
            <div className="content-wrapper">
                <div className="row gutters">
                    <div className="col-xl-12 col-lg-12 col-md-12 col-sm-12 col-12">
                        <div className="card m-0">
                            <div className="row no-gutters">
                                <div className="col-xl-4 col-lg-4 col-md-4 col-sm-3 col-3" style={{height: '85vh', overflow: 'scroll'}}>
                                    <button className="btn btn-primary mt-3 ml-3" onClick={toggleCreateGroup}><FontAwesomeIcon icon={faPlus} className="mr-3"/>New Group</button>
                                    <hr />
                                    <div className="users-container">
                                        <ul className="users">
                                            {state.groups.length!==0?state.groups.map(item => {return <GroupLabel key={item.ID} setCurrent={setCurrent} group={item}/>}):null}
                                        </ul>
                                    </div>
                                </div>
                                <Chat group={current} ws={props.ws} user={state.user} setCurrent={setCurrent} toggler={toggler}/>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
          <ModalCreateGroup show={createGrShow} toggle={toggleCreateGroup}/>
          <ModalUserProfile show={props.profileShow} toggle={props.toggleProfile} user={state.user} />
        </div>
    )
}
export default Main;