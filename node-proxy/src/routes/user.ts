import express, { Response } from 'express';
import authenticateJWT from '../middlewear/authentication'
import { AxiosResponse } from 'axios';
import { axiosBase } from '../config'
import { UserPatchRequestBody } from '../types/types';
import { Project, IRequest } from '../types/types';
import cors from 'cors';

const router = express.Router()
router.use(cors());

router.use(authenticateJWT)

const USER_BASE_PATH = '/api/v1/user'

// TODO remove password and make a different endpoint
interface User {
    _id: string;
    name: string;
    email: string;
    description: string;
    password: string;
    skills: string[];
} 

interface PublicUser {
    _id: string;
    name: string;
    description: string;
    skills: string[];
} 

router.get('/:id', async (req: IRequest, res: Response) => {
    const userId = req.params.id;
    const headers = {
        'UUID':`${req.id}`
    }
    let response: AxiosResponse<User | PublicUser>
    try {
        response = await axiosBase.get<User | PublicUser>(`${USER_BASE_PATH}/${userId}`, { headers });
    } catch(error) {
        res.status(error.response.status).json({data: error.response.data})
        return
    }
    const user: User | PublicUser = response.data;

    res.status(response.status).json({ data: user });
});

router.get('/projects/:id', async (req: IRequest, res: Response) => {
    const userId = req.params.id;
    const headers = {
        'UUID':`${req.id}`
    }
    const page = req.query.page
    const size = req.query.size
    let response: AxiosResponse<Project[]>
    try {
        response = await axiosBase.get<Project[]>(`/api/v1/user_to_project/${userId}?page=${page}&size=${size}`, { headers });
    } catch(error) {
        res.status(error.response.status).json({data: error.response.data})
        return
    }
    const projects: Project[] = response.data;

    res.status(response.status).json(projects);
});

router.patch('/:id', async (req:  IRequest<{ id: string }, object, UserPatchRequestBody>, res: Response) => {
    console.log("yuooo")
    const updatedUserBody: UserPatchRequestBody = req.body;
    const userId = req.params.id;
    const headers = {
        'UUID': `${req.id}`
    }
    
    let response: AxiosResponse<User>
    console.log(updatedUserBody)
    try{
        response = await axiosBase.patch(`${USER_BASE_PATH}/${userId}`, updatedUserBody, {headers});
    } catch(error) {
        res.status(error.response.status).json({data: error.response.data})
        return
    }
    const user: User = response.data;

    res.status(response.status).json({ data: user });
});

router.delete('/:id', async (req: IRequest, res: Response) => {

    const userId = req.params.id;
    const headers = {
        'UUID':`${req.id}`
    }
    
    let response: AxiosResponse<string> 

    try {
        response = await axiosBase.delete<string>(`${USER_BASE_PATH}/${userId}`, { headers });
    } catch(error) {
       res.status(error.response.status).json({data: error.response.data})
       return
    }

    res.status(response.status).json({ data: response.data });
});
  

export default router