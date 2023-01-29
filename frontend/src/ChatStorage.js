import React, {createContext, useReducer} from "react";

export const StorageContext = createContext({});

const initialState = {
    groups: [],
    invites: [],
    user: {},
};

export const actionTypes = {
    LOGIN: "LOGIN",
    LOGOUT: "LOGOUT",

    ADD_GROUPS: "ADD_GROUPS",
    ADD_GROUP: "ADD_GROUP",
    DELETE_GROUP: "DELETE_GROUP",

    ADD_MEMBER: "ADD_MEMBER",
    UPDATE_MEMBER: "UPDATE_MEMBER",
    DELETE_MEMBER: "DELETE_MEMBER",

    ADD_MESSAGES: "ADD_MESSAGES",
    ADD_MESSAGE: "ADD_MESSAGE",
    DELETE_MESSAGE: "DELETE_MESSAGE",

    ADD_INVITES: "ADD_INVITES",
    ADD_INVITE: "ADD_INVITE",
    UPDATE_INVITE: "UPDATE_INVITE",

    RESET_COUNTER: "RESET_COUNTER",
    SET_PROFILE_PICTURE: "SET_PROFILE_PICTURE",
    SET_GROUP_PICTURE: "SET_GROUP_PICTURE",
}

function reducer(state, action) {
    switch (action.type) {
        case actionTypes.LOGIN:
            return Login(state, action.payload);
        case actionTypes.LOGOUT:
            return Logout();

        case actionTypes.ADD_GROUPS:
            return AddGroups(state, action.payload);
        case actionTypes.ADD_GROUP:
            return AddGroup(state, action.payload);
        case actionTypes.DELETE_GROUP:
            return DeleteGroup(state, action.payload);

        case actionTypes.ADD_MEMBER:
            return AddMemberToGroup(state, action.payload);
        case actionTypes.UPDATE_MEMBER:
            return UpdateMember(state, action.payload);
        case actionTypes.DELETE_MEMBER:
            return DeleteMemberFromGroup(state, action.payload);
            
        case actionTypes.ADD_MESSAGES:
            return AddMessages(state, action.payload);
        case actionTypes.ADD_MESSAGE:
            return AddMessage(state, action.payload);
        case actionTypes.DELETE_MESSAGE:
            return DeleteMessage(state, action.payload);
            
        case actionTypes.ADD_INVITES:
            return AddInvites(state, action.payload);
        case actionTypes.ADD_INVITE:
            return AddInvite(state, action.payload);
        case actionTypes.UPDATE_INVITE:
            return UpdateInvite(state, action.payload);

        case actionTypes.RESET_COUNTER:
            return ResetCounter(state, action.payload);

        case actionTypes.SET_PROFILE_PICTURE:
            return SetProfilePicture(state, action.payload);
        case actionTypes.SET_GROUP_PICTURE:
            return SetGroupPicture(state, action.payload);

        default:
            throw new Error("Action not specified");
    }
}

const ChatStorage = ({children}) => {

    const [state, dispatch] = useReducer(reducer, initialState);

    return (
        <StorageContext.Provider value={[state, dispatch]}>
            {children}
        </StorageContext.Provider>
    );
}
export default ChatStorage;

function Login(state, payload) {
    let newState = {...state};
    newState.user = payload;
    return newState;
}

function Logout() {
    return {
        groups: [],
        invites: [],
        user: {},
    };
}

// GROUP HANDLERS

function AddGroups(state, payload) {
    let newState = {...state};
    newState.groups = [...newState.groups, ...payload];
    for (let i = 0; i < newState.groups.length; i++ ) {
        newState.groups[i].messages = [];
        newState.groups[i].unreadMessages = 0;
    }
    return newState;
}

function AddGroup(state, payload) {
    let newState = {...state};
    payload.messages = [];
    payload.unreadMessages = 0;
    newState.groups.push(payload);
    return newState;
}

function DeleteGroup(state, payload) {
    let newState = {...state};
    newState.groups = newState.groups.filter( (item) => { return item.ID !== payload } );
    return newState;
}

// MEMBER HANDLERS

function AddMemberToGroup(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.groupID) {
            newState.groups[i].Members.push(payload);
            return newState;
        }
    }
    throw new Error("Group not found");
}

function UpdateMember(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.groupID) {
            for (let j = 0; j < newState.groups[i].Members.length; j++) {
                if (newState.groups[i].Members[j].ID === payload.ID) {
                    newState.groups[i].Members[j] = payload;
                    return newState;
                }
            }
        }
    }
    throw new Error("Member not found");
}

function DeleteMemberFromGroup(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.groupID) {
            newState.groups[i].Members = newState.groups[i].Members.filter((item)=>{return item.ID !== payload.ID});
            return newState;
        }
    }
    throw new Error("Group not found");
}

// MESSAGE HANDLERS

function AddMessage(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.message.groupID) {
            newState.groups[i].messages.unshift(payload.message)
            if (!payload.current) {
                newState.groups[i].unreadMessages += 1;
            }
            return newState;
        }
    }
    throw new Error("Received message don't belong to any of your groups");
}

function AddMessages(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.groupID) {
            newState.groups[i].messages = [...newState.groups[i].messages, ...payload.messages ];
            return newState;
        }
    }
    throw new Error("Received messages don't belong to any of your groups");
}

function DeleteMessage(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.groupID) {
            for (let j = 0; j < newState.groups[i].messages.length; j++) {
                if (newState.groups[i].messages[j].messageID === payload.messageID) {
                    newState.groups[i].messages[j].text = ""; 
                    newState.groups[i].messages[j].files = [];
                }
            }
            return newState;
        }
    }
    throw new Error("Message not found")
}

// INVITES HANDLERS

function AddInvites(state, payload) {
    let newState = {...state};
    newState.invites = payload;
    return newState;
}

function AddInvite(state, payload) {
    let newState = {...state};
    newState.invites.push(payload);
    return newState;
}

function UpdateInvite(state, payload) {
    let newState = {...state};
    newState.invites = newState.invites.filter((item)=>{return item.ID !== payload.ID});
    newState.invites.push(payload);
    return newState;
}

// COUNTER

function ResetCounter(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload) {
            newState.groups[i].unreadMessages = 0;
            return newState;
        }
    }
    throw new Error("No such group in storage");
}

// PICTURE HANDLERS

function SetProfilePicture(state, payload) {
    let newState = {...state};
    newState.user.pictureUrl = payload
    return newState;
}

function SetGroupPicture(state, payload) {
    let newState = {...state}
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.groupID) {
            newState.groups[i].pictureUrl = payload.newUrl;
            return newState
        }
    }
    return newState;
}
