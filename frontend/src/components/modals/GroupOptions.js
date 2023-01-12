import React, { useContext, useState } from "react";
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import {UpdateGroupProfilePicture, DeleteGroupProfilePicture} from "../../requests/Groups";
import { actionTypes, StorageContext } from '../../ChatStorage';
import { GroupPicture } from "../Pictures";

export const ModalGroupOptions = (props) => {

    const [, dispatch] = useContext(StorageContext);

    const [file, setFile] = useState({});

    const [message, setMessage] = useState("");

    const updatePicture = async(e) => {
        e.preventDefault();
        let data = new FormData();
        data.append("avatarFile", file);
    
        let response = await UpdateGroupProfilePicture(data, props.group.ID);

        if (response.status === 200) {
            setMessage("Image uploaded successfully");
            dispatch({type: actionTypes.SET_GROUP_PICTURE, payload: {newUrl: response.data.newUrl, groupID: props.group.ID}});
            let timestamp = new Date().getTime();
            document.getElementById("profilePicture").src = "https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.group.pictureUrl+"?"+timestamp;
            document.getElementById("customFile").value= null;
        } else {
            setMessage(response.data.err);
        }
        setTimeout(function() {
            setMessage("");
        }, 3000);
    };

    const deletePicture = async() => {
        let response = await DeleteGroupProfilePicture(props.group.ID)

        if (response.status === 200) {
            setMessage("Image deleted successfully");
            dispatch({type: actionTypes.DELETE_GROUP_PROFILE_PICTURE, payload: props.group.ID})
            let timestamp = new Date().getTime();
            document.getElementById("profilePicture").src = "https://chatprofilepics.s3.eu-central-1.amazonaws.com/"+props.group.ID+"?"+timestamp;
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
                    Group Profile
                </ModalHeader>
                <ModalBody>
                    <div class="container">
                        <div className="row d-flex justify-content-center">
                            <div className="text-center card-box">
                                <div className="member-card">
                                    {message}
                                    <div className="mx-auto profile-image-holder">
                                        <GroupPicture pictureUrl={props.group.pictureUrl} imageID="profilePicture"/>
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
                                </div>
                            </div>
                        </div>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 