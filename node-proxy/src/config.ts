import axios from 'axios';

export const secretKey: string = process.env.JWT_SECRET_KEY; 

export const axiosBase = axios.create({
  baseURL: process.env.BACK_END_URL || "http://backend:8080", // Set from compose for local dev or production
});

