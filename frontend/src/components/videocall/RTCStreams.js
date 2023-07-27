
export const actionTypes = {
    NEW_STREAM: "NEW_STREAM",
    DELETE_STREAM: "DELETE_STREAM",

    SET_USER_INFO: "SET_USERNAME",

    TOGGLE_MUTE: "TOGGLE_MUTE",

    END_SESSION: "END_SESSION",
};

export const RTCStreamsReducer = (state, action) => {
    switch (action.type) {
        case actionTypes.NEW_STREAM:
            return NewStream(state, action.payload);
        case actionTypes.DELETE_STREAM:
            return DeleteStream(state, action.payload);

        case actionTypes.END_SESSION:
            return EndSession(state);

        case actionTypes.SET_USER_INFO:
            return SetUserInfo(state, action.payload);
        case actionTypes.TOGGLE_MUTE:
            return ToggleMute(state, action.payload);
        
        default:
            console.log("Unknown dispatch type: ", action.type);
    }
};

const NewStream = (state, payload) => {
    console.log("New Stream");
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

const DeleteStream = (state, payload) => {
    let newState = [...state];
    newState = newState.filter((stream) => { return stream.id === payload });
    return newState
};

const EndSession = (state) => {
    state.forEach(rtcStream => {
        rtcStream.stream.getTracks().forEach( track => track.stop());
    });

    return [];
}

const SetUserInfo = (state, payload) => {
    console.log("Set Username");
    let newState = [...state];

    for (let i = 0; i < state.length; i++) {
        if (newState[i].stream.id === payload.streamID) {
            newState[i].username = payload.username;
            if (payload.videoEnabled !== undefined) newState[i].videoEnabled = payload.videoEnabled;
            if (payload.audioEnabled !== undefined) newState[i].audioEnabled = payload.audioEnabled;
            return newState;
        }
    }
    newState.push({username: payload.username, stream: {id: payload.streamID}, videoEnabled: payload.videoEnabled, audioEnabled: payload.audioEnabled} );

    return newState;
};

const ToggleMute = (state, payload) => {
    console.log("Mute toggled");

    let newState = [...state];

    for (let i = 0; i < state.length; i++) {
        if (newState[i].stream.id === payload.streamID) {
            if (payload.videoEnabled !== undefined) newState[i].videoEnabled = payload.videoEnabled;
            if (payload.audioEnabled !== undefined) newState[i].audioEnabled = payload.audioEnabled;

            return newState;
        }
    }

    throw new Error("Mute called on unexisting object");
};