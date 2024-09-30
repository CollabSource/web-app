import React, { useContext } from 'react';
import { AuthContext } from '../../../context/AuthContext';

const SignOutButton: React.FC = () => { 

    const { authDispatch } = useContext(AuthContext)

    const signOut = () => {
        localStorage.removeItem("auth_context_state");
        authDispatch({ type: 'LOG_OUT' })
    }

    return(
        <div>
            <button id='sign-out-button' onClick={() => signOut()}>Sign Out</button>
        </div>
    );
}

export default SignOutButton;