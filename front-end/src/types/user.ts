export interface User {
  _id: string;
  name: string;
  email: string;
  description: string;
  skills: string[];
} 

export interface PublicUser {
  _id: string;
  name: string;
  description: string;
  skills: string[];
} 

export type UserType = User | PublicUser | null