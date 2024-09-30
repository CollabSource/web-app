import { Request } from "express";
import { ParsedQs } from 'qs'

export interface UserRegisterRequestBody {
    name: string;
    password: string;
    email: string;
    description: string;
    skills: string[];
}

export interface UserPatchRequestBody {
    name : string;
    description : string;
    skills : string[];
}

export interface Project {
    id : string
    ownerId: string
    name : string
    description : string
    category : string
    tags : string[]
    seeking: string[]
  }


export interface PaginatedResponseBody<T> {
	items: T[]
	page: number
	hasNext: boolean
}

export interface IRequest<P = { id: string }, ResBody = unknown, ReqBody = unknown, ReqQuery extends ParsedQs = ParsedQs> extends Request<P, ResBody, ReqBody, ReqQuery> {
  id: string;
}