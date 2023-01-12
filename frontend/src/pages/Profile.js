import React, { useContext, useState } from "react";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import {ChangePassword, UpdateProfilePicture, DeleteProfilePicture }from "../requests/Users";
import { actionTypes, StorageContext } from '../ChatStorage';
import { UserPicture } from "../components/Pictures";

export const ModalUserProfile = (props) => {

    const [, dispatch] = useContext(StorageContext);

    const [oldPassword, setOldPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [repeatPassword, setRepeatPassword] = useState("");

    const [file, setFile] = useState({});

    const [message, setMessage] = useState("");

    const changePassword = async(e) => {
            e.preventDefault();
            let response;

            try {
                response = await ChangePassword(oldPassword, newPassword, repeatPassword);
                if (response.status === 200) {
                    setMessage("Password changed");
                }
            } catch(err) {
                if (err.response !== undefined) setMessage(err.response.data.err);
                else setMessage(err.message);
            }
            setRepeatPassword("");
            setNewPassword("");
            setOldPassword("");

            document.getElementById("oldpassword").value= "";
            document.getElementById("newpassword").value= "";
            document.getElementById("rpassword").value= "";

            setTimeout(function() {
                setMessage("");
            }, 3000);
    }

    const updatePicture = async(e) => {
        e.preventDefault();

        let data = new FormData();
        data.append("avatarFile", file);
    
        let response = await UpdateProfilePicture(data);
        if (response.status === 200) {
            setMessage("Image uploaded succesfully");
            console.log(response);
            dispatch({type: actionTypes.SET_PROFILE_PICTURE, payload: response.data.newUrl});
            let timestamp = new Date().getTime();
            document.getElementById("profilePicture").src = "https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.user.pictureUrl+"?"+timestamp;
            document.getElementById("customFile").value= null;

        } else {
            setMessage(response.data.err);
        }
        setTimeout(function() {
            setMessage("");
        }, 3000);
    };

    const deletePicture = async() => {

        let response = await DeleteProfilePicture();
        if (response.status === 200) {
            setMessage("Image deleted successfully");
            dispatch({type: actionTypes.SET_PROFILE_PICTURE, payload: ""});
            let timestamp = new Date().getTime();
            document.getElementById("profilePicture").src = "https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.user.pictureUrl+"?"+timestamp;
        } else {
            setMessage(response.data.err);
        }
        setTimeout(function() {
            setMessage("");
        }, 3000);
    };

    return (
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    User Profile
                </ModalHeader>
                <ModalBody>
                    <div className="container">
                        <div className="row d-flex justify-content-center">
                            <div className="text-center card-box">
                                <div className="member-card">
                                    {message}
                                    <div className="mx-auto profile-image-holder">
                                        <UserPicture pictureUrl={props.user.pictureUrl} imageID="profilePicture"/>
                                    </div>
                                    <div>
                                        <h4>{props.name}</h4>
                                    </div>
                                    <hr />
                                    <h3>Change profile picture</h3>
                                    <form encType="multipart/form-data" action="">
                                        <input type="file" className="form-control" id="customFile" onChange={e => setFile(e.target.files[0])}  />
                                        <div className="text-center mt-2">
                                            <button className="btn btn-primary text-center w-100" onClick={updatePicture}>Upload</button>
                                        </div>
                                    </form>
                                    <div className="text-center mt-4">
                                        <button className="btn btn-danger text-center w-100" onClick={deletePicture}>Delete Picture</button>
                                    </div>
                                    <hr />
                                    <form className="mt-4">
                                        <h3> Change Password </h3>
                                        <div className="mb-3 text-center">
                                            <label htmlFor="pass" className="form-label">Old Password</label>
                                            <input name="oldpassword" type="password" className="form-control" id="oldpassword" onChange={(e) => setOldPassword(e.target.value)} />
                                        </div>
                                        <div className="mb-3 text-center">
                                            <label htmlFor="pass" className="form-label">New Password</label>
                                            <input name="newpassword" type="password" className="form-control" id="newpassword" onChange={(e) => setNewPassword(e.target.value)} />
                                        </div>
                                        <div className="mb-3 text-center">
                                            <label htmlFor="pass" className="form-label">Repeat Password</label>
                                            <input name="rpassword" type="password" className="form-control" id="rpassword" onChange={(e) => setRepeatPassword(e.target.value)} />
                                        </div>
                                            
                                        <div className="text-center">
                                            <button className="btn btn-primary text-center w-100" onClick={changePassword}>Change password</button>
                                        </div>
                                    </form>
                                </div>
                            </div>
                        </div>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 