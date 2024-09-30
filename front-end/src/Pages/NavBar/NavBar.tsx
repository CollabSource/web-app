import React from 'react';
import './NavBar.css';
import { useContext } from 'react';
import SignUpButton from './Components/SignUpButton';
import NavBarButton from './Components/NavBarButton';
import LogInButton from './Components/LogInButton';
import { AuthContext } from '../../context/AuthContext';
import SignOutButton from './Components/SignOutButton';
import { ReactComponent as ProfileSVG } from "../../assets/svg/user-profile-gray-svgrepo-com.svg"
import { useNavigate } from 'react-router-dom';

const NavBar: React.FC = () => {
    // const navigate = useNavigate()
    const {loggedIn, userID} = useContext(AuthContext)

    return(
        <div id="nav-bar">
            <div id="left-nav-items">
                <a id="title" href="/">CollabSource</a>
                <NavBarButton text="Explore" pathToPage="/explore"/>
                <NavBarButton text="About" pathToPage=""/>
                <NavBarButton text="Other" pathToPage=""/>
            </div>
            {!loggedIn ?
                <div id="right-nav-items">
                    <LogInButton/>
                    <SignUpButton/>
                </div>
                :
                <div id="right-nav-items">
                    <a className="svg-cont-prof" href={`/user/${userID}`}>
                        <ProfileSVG className='prof-pic-nav'/>
                    </a>
                    <SignOutButton/>
                </div>
            }
        </div>
    );
}

export default NavBar;