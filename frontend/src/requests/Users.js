import axiosObject, {userService} from "./Setup";

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
    return await axiosObject.post(userService+"/register", {
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
    return await axiosObject.get(userService+"/verify-account/"+code);
}

export async function Login(email, password) {
    if (email.trim() === "") {
        throw new Error("Email cannot be blank");
    }
    if (password.trim() === "") {
        throw new Error("Password cannot be blank");
    }
    return await axiosObject.post(userService+"/login", {
        email: email,
        password: password,
    }, {
        withCredentials: true,
    });
}


export async function GetUser() {
    return await axiosObject.get(userService+'/user');
}

export async function Logout() {
    return await axiosObject.post(userService+"/signout", {}, {
        withCredentials: true,
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

    return await axiosObject.put(userService+"/change-password", {
        "oldPassword": oldPassword,
        "newPassword": newPassword,
        "repeatPassword": repeatPassword,
    });
}

export async function UpdateProfilePicture(image) {
    return await axiosObject.post(userService+"/set-image", image, {
        headers: {
            'Content-Type': 'multipart/form-data',
        }
    })
}

export async function DeleteProfilePicture() {
    return await axiosObject.delete(userService+"/delete-image");
}

export async function ForgotPassword(email) {
    return await axiosObject.get(userService+"/forgot-password?email="+email);
}

export async function ResetForgottenPassword(resetCode, newPassword, repeatPassword) {
    if (newPassword === "") {
        throw new Error("password cannot be blank");
    }
    if (newPassword.length <  6) {
        throw new Error("password must be at least 6 characters long");
    }
    if (repeatPassword !== newPassword) {
        throw new Error("Passwords don't match");
    }

    return await axiosObject.patch(userService+"/reset-password/"+resetCode, {
        "newPassword": newPassword,
        "repeatPassword": repeatPassword,
    });

}