
export const actionTypes = {
    NEW_STREAM: "NEW_STREAM",
    DELETE_STREAM: "DELETE_STREAM",

    SET_USERNAME: "SET_USERNAME",

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

        case actionTypes.SET_USERNAME:
            return SetUsername(state, action.payload);
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

const SetUsername = (state, payload) => {
    console.log("Set Username");
    let newState = [...state];

    for (let i = 0; i < state.length; i++) {
        if (newState[i].stream.id === payload.streamID) {
            newState[i].username = payload.username;
            return newState;
        }
    }
    newState.push({username: payload.username, stream: {id: payload.streamID}});

    return newState;
};