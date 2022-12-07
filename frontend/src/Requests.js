const groupsService = 'http://localhost:8081/api';
const messageService = 'http://localhost:8082/api';
const userService = 'http://localhost:8083/api';
const wsService = 'ws://localhost:8084';

let axios = require('axios').default;
axios.defaults.headers.common['Content-Type'] = "application/json";

export async function Register(email, username, password, rpassword) {
    if (email.trim() === "") {
        throw new Error("Email can't be blank");
    }
    if (username.trim() === "") {
        throw new Error("Username can't be blank");
    }
    if (password.trim() === "") {
        throw new Error("Password can't be blank");
    }
    if (password !== rpassword) {
        throw new Error("Passwords don't match");
    }
    return await axios.post(userService+"/register", {
        username: username, 
        email: email,
        password: password,
        rpassword: rpassword,
    })
}

export async function VerifyAccount(code) {
    if (code === "") {
        throw new Error("code can't be blank");
    }
    return await axios.get(userService+"/verify-account/"+code);
}

export async function Login(email, password) {
    if (email.trim() === "") {
        throw new Error("Email cannot be blank");
    }
    if (password.trim() === "") {
        throw new Error("Password cannot be blank");
    }
    return await axios.post(userService+"/login", {
        email: email,
        password: password,
    }, {
        withCredentials: true,
    });
}


export async function GetUser() {
    return await axios.get(userService+'/user');
}

export async function GetInvites(offset) {
    return await axios.get(groupsService+"/invites?num=8&offset="+offset);
}

export async function GetGroups() {
    return await axios.get(groupsService+"/group");
}

export async function Logout() {
    return await axios.post(userService+"/signout", {}, {
        withCredentials: true,
    });
}

export async function LoadMessages(groupID, offset) {
    return await axios.get(messageService+"/group/"+groupID+"/messages?num=8&offset="+offset);
}

export async function CreateGroup(name) {
    return await axios.post(groupsService+"/group", {
        "name": name,
    })
}

export async function DeleteGroup(groupID) {
    return await axios.delete(groupsService+"/group/"+groupID);
}

export async function SendGroupInvite(username, groupID) {
    return await axios.post(groupsService+"/invites", {
        "target": username,
        "group": groupID
    });
}

export async function RespondGroupInvite(inviteID, answer) {
    return await axios.put(groupsService+"/invites/"+inviteID, {
        "answer": answer
    })
}

export async function DeleteMember(groupID, memberID) {
    return await axios.delete(groupsService+"/group/"+groupID+"/member/"+memberID);
}

export async function SetRights(groupID, memberID, adding, deletingMessages, deletingMembers, setting) {
    return await axios.patch(groupsService+"/group/"+groupID+"/member/"+memberID, {
        "adding": adding,
        "deletingMessages": deletingMessages,
        "deletingMembers": deletingMembers,
        "setting": setting
    });
}

export async function ChangePassword(oldPassword, newPassword, repeatPassword) {
    if (newPassword === "") {
        throw new Error("password cannot be blank");
    }
    if (newPassword.length <  6) {
        throw new Error("password must be at least 6 characters long");
    }
    if (repeatPassword !== newPassword) {
        throw new Error("Passwords don't match");
    }

    return await axios.put(userService+"/change-password", {
        "oldPassword": oldPassword,
        "newPassword": newPassword,
    });
}

export async function UpdateProfilePicture(image) {
    return await axios.post(userService+"/set-image", image, {
        headers: {
            'Content-Type': 'multipart/form-data',
        }
    })
}

export async function DeleteProfilePicture() {
    return await axios.delete(userService+"/delete-image");
}

export async function UpdateGroupProfilePicture(imageForm, groupID) {
    return await axios.post(groupsService+"/group/"+groupID+"/image", imageForm, {
        headers: {
            'Content-Type': 'multipart/form-data',
        }
    });
}

export async function DeleteGroupProfilePicture(groupID) {
    return await axios.delete(groupsService+"group/"+groupID+"/image");
}

export async function GetWebsocket() {
    let access_token = window.localStorage.getItem("token")
    let socket = new WebSocket(wsService+'/ws?authToken='+access_token);
    socket.onopen = () => {
        console.log("Websocket openned");
    };
    socket.onclose = () => {
        console.log("Websocket closed");
    };
    return socket;
}

async function refreshAccessToken() {
    let response = await axios.post(userService+"/refresh", {}, {
        withCredentials: true,
    })
    console.log(response);
    if (response.accessToken !== undefined) {
        window.localStorage.setItem("token", response.accessToken);
    }
}

// Request interceptor for API calls
axios.interceptors.request.use(
    async config => {
        let accessToken = window.localStorage.getItem("token")
        config.headers = { 
            'Authorization': `Bearer ${accessToken}`,
            'Accept': 'application/json',
        }
        return config;
    },
    error => {
      Promise.reject(error)
  });
  
  // Response interceptor for API calls
  axios.interceptors.response.use((response) => {
    return response
  }, async function (error) {
    const originalRequest = error.config;
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      await refreshAccessToken();
      let access_token = window.localStorage.getItem("token");
      axios.defaults.headers.common['Authorization'] = 'Bearer ' + access_token;
      return axios(originalRequest);
    }
    return Promise.reject(error);
  });