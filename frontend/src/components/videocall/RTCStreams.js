
export const actionTypes = {
    NEW_STREAM: "NEW_STREAM",
    
    SET_USER_INFO: "SET_USERNAME",
    USER_DISCONNECTED: "USER_DISCONNECTED",

    END_SESSION: "END_SESSION",
};

export const RTCStreamsReducer = (state, action) => {
    switch (action.type) {
        case actionTypes.NEW_STREAM:
            return NewStream(state, action.payload);
        case actionTypes.USER_DISCONNECTED:
            return UserDisconnected(state, action.payload);

        case actionTypes.END_SESSION:
            return EndSession(state);

        case actionTypes.SET_USER_INFO:
            return SetUserInfo(state, action.payload);
        
        default:
            console.log("Unknown dispatch type: ", action.type);
    }
};

const NewStream = (state, payload) => {
    let newState = [...state];
    for (let i = 0; i < newState.length; i++) {
        if (newState[i].stream.id === payload.id) {
            newState[i].stream = payload;
            return newState;
        }
    }
    newState.push({stream: payload});
    return newState;
}

const EndSession = (state) => {
    state.forEach(rtcStream => {
        rtcStream.stream.getTracks().forEach( track => track.stop());
    });

    return [];
}

const SetUserInfo = (state, payload) => {
    let newState = [...state];

    for (let i = 0; i < state.length; i++) {
        if (newState[i].stream.id === payload.streamID) {
            newState[i].username = payload.username;
            newState[i].memberID = payload.memberID;
            return newState;
        }
    }
    newState.push({username: payload.username, memberID: payload.memberID, stream: {id: payload.streamID}});

    return newState;
};

const UserDisconnected = (state, payload) => {
    let newState = [...state];
    newState = newState.filter((stream) => { return stream.memberID !== payload });
    return newState
};