import express, { Request, Response } from 'express';
import authenticateJWT from '../middlewear/authentication'
import { AxiosResponse } from 'axios';
import { axiosBase } from '../config'
import { UserPatchRequestBody } from '../types/types';
import cors from 'cors';

const router = express.Router()
router.use(cors());

if (process.env.USE_JWT === 'true') {
    router.use(authenticateJWT)
}

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

router.get('/:id', async (req: Request, res: Response) => {

    const userId = req.params.id;
    const headers = {
        'UUID':`${req.email}`
    }
    
    let response: AxiosResponse<User>
    try {
        response = await axiosBase.get<User>(`${USER_BASE_PATH}/${userId}`, { headers });
    } catch(error) {
        res.status(error.response.status).json({data: error.response.data})
        return
    }
    const user: User = response.data;

    res.status(response.status).json({ data: user });
});

// axios.interceptors.request.use(request => {
//     console.log('Starting Request', request);
//     return request;
//   });

router.patch('/:id', async (req:  Request<{id: string}, object, UserPatchRequestBody>, res: Response) => {
    const updatedUserBody: UserPatchRequestBody = req.body;
    const userId = req.params.id;
    const headers = {
        'UUID':`${req.email}`
    }
    
    let response: AxiosResponse<User>

    try{
        response = await axiosBase.patch<User>(`${USER_BASE_PATH}/${userId}`, updatedUserBody, {headers});
    } catch(error) {
        res.status(error.response.status).json({data: error.response.data})
        return
    }
    const user: User = response.data;

    res.status(response.status).json({ data: user });
});

router.delete('/:id', async (req: Request, res: Response) => {

    const userId = req.params.id;
    const headers = {
        'UUID':`${req.email}`
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