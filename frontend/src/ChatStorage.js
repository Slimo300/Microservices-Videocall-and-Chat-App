import React, {createContext, useReducer} from "react";

export const StorageContext = createContext({});

const initialState = {
    groups: [],
    notifications: [],
    user: {},
};

export const actionTypes = {
    LOGIN: "LOGIN",
    LOGOUT: "LOGOUT",
    SET_GROUPS: "SET_GROUPS",
    NEW_GROUP: "NEW_GROUP",
    DELETE_GROUP: "DELETE_GROUP",
    ADD_MEMBER: "ADD_MEMBER",
    DELETE_MEMBER: "DELETE_MEMBER",
    SET_MESSAGES: "SET_MESSAGES",
    ADD_MESSAGE: "ADD_MESSAGE",
    ADD_MESSAGES: "ADD_MESSAGES",
    SET_NOTIFICATIONS: "SET_NOTIFICATIONS",
    ADD_NOTIFICATION: "ADD_NOTIFICATION",
    DELETE_NOTIFICATION: "DELETE_NOTIFICATION",
    RESET_COUNTER: "RESET_COUNTER",
    NEW_PROFILE_PICTURE: "NEW_PROFILE_PICTURE",
    DELETE_PROFILE_PICTURE: "DELETE_PROFILE_PICTURE",
    NEW_GROUP_PROFILE_PICTURE: "NEW_GROUP_PROFILE_PICTURE",
    DELETE_GROUP_PROFILE_PICTURE: "DELETE_GROUP_PROFILE_PICTURE",
}

function reducer(state, action) {
    switch (action.type) {
        case actionTypes.LOGIN:
            return Login(state, action.payload);
        case actionTypes.LOGOUT:
            return Logout();
        case actionTypes.SET_GROUPS:
            return SetGroups(state, action.payload);
        case actionTypes.NEW_GROUP:
            return NewGroup(state, action.payload);
        case actionTypes.DELETE_GROUP:
            return DeleteGroup(state, action.payload);
        case actionTypes.ADD_MEMBER:
            return AddMemberToGroup(state, action.payload);
        case actionTypes.DELETE_MEMBER:
            return DeleteMemberFromGroup(state, action.payload);
        case actionTypes.SET_MESSAGES:
            return SetMessages(state, action.payload);
        case actionTypes.ADD_MESSAGE:
            return AddMessage(state, action.payload);
        case actionTypes.ADD_MESSAGES:
            return AddMessages(state, action.payload);
        case actionTypes.SET_NOTIFICATIONS:
            return SetInvites(state, action.payload);
        case actionTypes.ADD_NOTIFICATION:
            return AddInvite(state, action.payload);
        case actionTypes.DELETE_NOTIFICATION:
            return DeleteInvite(state, action.payload);
        case actionTypes.RESET_COUNTER:
            return ResetCounter(state, action.payload);
        case actionTypes.DELETE_NOTIFICATION:
            return DeleteNotification(state, action.payload);
        case actionTypes.NEW_PROFILE_PICTURE:
            return NewProfilePicture(state, action.payload);
        case actionTypes.DELETE_PROFILE_PICTURE:
            return DeleteProfilePicture(state);
        case actionTypes.NEW_GROUP_PROFILE_PICTURE:
            return NewGroupProfilePicture(state, action.payload);
        case actionTypes.DELETE_GROUP_PROFILE_PICTURE:
            return DeleteGroupProfilePicture(state, action.payload);
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
    return initialState;
}

function SetGroups(state, payload) {
    let newState = {...state};
    newState.groups = payload;
    for (let i = 0; i < newState.groups.length; i++ ) {
        newState.groups[i].messages = [];
        newState.groups[i].unreadMessages = 0;
    }
    return newState;
}

function NewGroup(state, payload) {
    let newState = {...state};
    payload.messages = [];
    payload.unreadMessages = 0;
    newState.groups = [...newState.groups, payload];
    return newState;
}

function DeleteGroup(state, payload) {
    let newState = {...state};
    newState.groups = newState.groups.filter( (item) => { return item.ID !== payload } );
    return newState;
}

function AddMemberToGroup(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.group_id) {
            newState.groups[i].Members.push(payload);
            return newState;
        }
    }
    throw new Error("Group not found");
}

function DeleteMemberFromGroup(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.group_id) {
            newState.groups[i].Members = newState.groups[i].Members.filter((item)=>{return item.ID !== payload.ID});
            return newState
        }
    }
    throw new Error("Group not found");
}

function SetMessages(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.group) {
            newState.groups[i].messages = payload.messages.reverse()
            return newState;
        }
    }
    throw new Error("Received messages don't belong to any of your groups");
}

function AddMessage(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.message.group) {
            newState.groups[i].messages = [...newState.groups[i].messages, payload.message];
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
        if (newState.groups[i].ID === payload.group) {
            newState.groups[i].messages = [...payload.messages.reverse(), ...newState.groups[i].messages];
            return newState;
        }
    }
    throw new Error("Received messages don't belong to any of your groups");
}

function SetInvites(state, payload) {
    let newState = {...state};
    newState.notifications = payload;
    return newState;
}

function AddInvite(state, payload) {
    let newState = {...state};
    newState.notifications = [...newState.notifications, payload];
    return newState;
}

function DeleteInvite(state, payload) {
    let newState = {...state};
    newState.notifications = newState.notifications.filter( (item) => { return item.ID !== payload } );
    return newState;
}

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

function DeleteNotification(state, payload) {
    let newState = {...state};
    newState.notifications = newState.notifications.filter( (item) => { return item.ID !== payload } );
    return newState;
}

function NewProfilePicture(state, payload) {
    let newState = {...state};
    newState.user.pictureUrl = payload;
    return newState;
}

function DeleteProfilePicture(state) {
    let newState = {...state};
    newState.user.pictureUrl = "";
    return newState;
}

function NewGroupProfilePicture(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload.groupID) {
            newState.groups[i].pictureUrl = payload.newURI;
            return newState;
        }
    }
    return newState;
}

function DeleteGroupProfilePicture(state, payload) {
    let newState = {...state};
    for (let i = 0; i < newState.groups.length; i++) {
        if (newState.groups[i].ID === payload) {
            newState.groups[i].pictureUrl = "";
            return newState;
        }
    }
    return newState;
}