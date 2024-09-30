import './App.css';
import './Pages/HomePage/HomePage';
import HomePage from './Pages/HomePage/HomePage';
import NavBar from './Pages/NavBar/NavBar';
import CreateProject from './Pages/CreateProject/CreateProject';
import ProjectPage from './Pages/ProjectPage/ProjectPage';
import ErrorPage from './error-page';
import User from './Pages/User/User'


import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import Explore from './Pages/Explore/Explore';

const router = createBrowserRouter([
  {
    path: "/",
    element: <HomePage/>,
    errorElement: <ErrorPage/>
  },
  {
    path: "/explore",
    element:<Explore/>,
    errorElement: <ErrorPage/>
  },
  {
    path: "/about",
    element: <div>about</div>,
    errorElement: <ErrorPage/>
  },
  {
    path: "/project/:id",
    element: <ProjectPage/>,
    errorElement: <ErrorPage/>
  },
  {
    path: "/create_project",
    element: <CreateProject/>,
    errorElement: <ErrorPage/>
  },
  {
    path: "/user/:id",
    element: <User/>,
    errorElement: <ErrorPage/>
  }
]);

function App() {

  return (
    <div>
      <NavBar/>
      <RouterProvider router={router}/>
    </div>
    
  );
}

export default App;
