import express, {Request, Response} from 'express';
import authenticateJWT from '../middlewear/authentication';
import { AxiosResponse } from 'axios';
import { axiosBase } from '../config';
import { Project, PaginatedResponseBody } from '../types/types';
import cors from 'cors';

const router = express.Router();
router.use(cors());

if (process.env.USE_JWT === 'true') {
  router.use(authenticateJWT)
}
const PROJECT_BASE_PATH = '/api/v1/projects'

interface ProjectFilter {
  categories:    string[]
	searchQuery :  string
}
router.post('/', async (req: Request, res: Response) => {
  const headers = {
    'userEmail': `${req.email}`
  }
  const page = req.query.page
  const size = req.query.size
  try {
    const projectQuery: ProjectFilter = req.body
    const response: AxiosResponse<PaginatedResponseBody<Project>> = await axiosBase.post<PaginatedResponseBody<Project>>(`${PROJECT_BASE_PATH}?page=${page}&size=${size}`, projectQuery, { headers });
    const project: PaginatedResponseBody<Project> = response.data;
    res.status(200).json(project)

  } catch (error) {
    if (error.response) {
      res.status(error.response.status).json({ message: error.response.data });
    } else {
      res.status(500).json({ message: error });
    }
    // res.status(error.response.status).json({ message: error.response.data})
  }
});
export default router
