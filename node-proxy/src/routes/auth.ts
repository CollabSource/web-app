import express, { Request, Response } from 'express';
import jwt from 'jsonwebtoken';
import { secretKey, axiosBase } from '../config'
import { UserRegisterRequestBody } from '../types/types'
import cors from 'cors';
import { AxiosResponse } from 'axios';

const router = express.Router()
router.use(cors());

const BASE_PATH = '/api/v1'
  
router.post('/register', async (req: Request<object, object, UserRegisterRequestBody>, res: Response) => {
    const reqBody: UserRegisterRequestBody  = req.body;

    if (!reqBody.email || !reqBody.password) {
        return res.status(400).json({ message: 'Username and password are required' });
    }

    try {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const response = await axiosBase.post(`${BASE_PATH}/register`, reqBody);

        res.status(201).json({ message: 'User registered successfully' });
   
    } catch (error) {
        if (error.response.status === 422) {
            res.status(error.response.status).json({ message: 'Failed to register user' });
        } else {
            res.status(error.response.status).json({ message: 'Something went wrong'});
        } 
    }  
});

router.post('/login', async (req: Request, res: Response) => {
    const { email, password }: { email: string; password: string } = req.body;

    // TODO make sure username is a valid email

    if (!email || !password) {
        return res.status(400).json({ message: 'Username and password are required' });
    }

    try {
        const response: AxiosResponse<string> = await axiosBase.post(`${BASE_PATH}/login`, { email, password });
        const id = response.data
        const token = jwt.sign({ id }, secretKey, { expiresIn: '3h' });
        res.status(200).json({ token, userId: response.data });
    } catch (error) {
        console.log(error)
        return res.status(error.response.status).json({ message: error.response.data });
    }  
});

export default router