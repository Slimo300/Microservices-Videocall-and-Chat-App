import React, {useState} from 'react';
import { Modal, ModalHeader, ModalBody } from 'reactstrap';
import {ForgotPassword} from '../../requests/Users';

export const ModalForgotPassword = (props) => {

    const [msg, setMsg] = useState("");
    const [email, setEmail] = useState("");


    const submit = async() => {
        try {
            let response = await ForgotPassword(email);
            if (response.status === 200) {
                setMsg("We sent password reseting link to your email");
            }
        } catch(err) {
            setMsg(err.response.data.err);
        }
        setTimeout(function () {    
            setMsg("");
            props.toggle();
        }, 2000);
    }

    return (
        <Modal id="buy" tabIndex="-1" role="dialog" isOpen={props.show} toggle={props.toggle}>
            <div role="document">
                <ModalHeader toggle={props.toggle} className="bg-dark text-primary text-center">
                    Forgot Password?
                </ModalHeader>
                <ModalBody>
                <div>
                        {msg!==""?<h5 className="mb-4 text-danger">{msg}</h5>:null}
                        <form onSubmit={submit}>
                            <div className="form-group">
                                <label htmlFor="email">Type in your email:</label>
                                <input name="name" type="text" className="form-control" id="gr_name" onChange={(e)=>{setEmail(e.target.value)}}/>
                            </div>
                            <div className="form-row text-center">
                                <div className="col-12 mt-2">
                                    <button type="submit" className="btn btn-dark btn-large text-primary">Send</button>
                                </div>
                            </div>
                        </form>
                    </div>
                </ModalBody>
            </div>
        </Modal>
    );
} 