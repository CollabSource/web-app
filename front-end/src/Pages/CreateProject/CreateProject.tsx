import React, { ChangeEvent, useContext, useState} from 'react'
import { AuthContext } from '../../context/AuthContext'
import { ProjectBase } from "../../types/project"
import axiosBase from '../../config/axiosConfig'
import './CreateProject.css'
import { useNavigate } from 'react-router-dom'
import { AxiosResponse } from 'axios'
import checkFormError  from './errorHandling'
import PromptAccount from '../Modals/PromptAccount'
import LogInModal from '../Modals/LogInModal'
import SignUpModal from '../Modals/SignUpModal'


export interface FormFieldsError {
  nameErr: string
  descriptionErr: string
  categoryErr: string
  tagsErr: string
  seekingErr: string
}

export interface ProjectForm {
  name: string
  description: string
  category: string
  tags: string,
  seeking: string
}

interface PostData extends ProjectBase {
  ownerId: string
}

const CreateProject:React.FC = () => {
  const navigate = useNavigate()
  const { loggedIn, userID } = useContext(AuthContext)
  const [showPromptAccount, setShowPromptAccount ] = useState<Boolean>(false)
  const [showLogin, setShowLogin ] = useState<Boolean>(false)
  const [showSignUp, setShowSignUp ] = useState<Boolean>(false)
  
  const [formData, setFormData] = useState<ProjectForm>({
    name       : '',
    description: '',
    category   : '',
    tags       : '',
    seeking    : '',
  })
  const noErrorObj: FormFieldsError = {nameErr: "", descriptionErr: "", categoryErr: "", tagsErr: "", seekingErr: ''}
  const [formFieldError, setFormFieldError] = useState<FormFieldsError>(noErrorObj);
  const [submitError, setSubmitError] = useState<String>("")

  const categories = ['business', 'software engineering', 'art'] // TODO get from backend?

  const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData({           
        ...formData,
        [name]: value,
      });
  };

  const selectCategory = (e: ChangeEvent<HTMLSelectElement>) => {
    let newState = formData;
    newState.category = e.target.value;
    setFormData(newState);
  }

  const onSubmit = async () => {

    if (!loggedIn) {
      setShowPromptAccount(true)
      return
    }
    if (checkFormError(noErrorObj, formData, setFormFieldError)) {
      return
    }
    const tagArr = formData.tags.split(',').map(e => e.trim())
    const seekingArr = formData.seeking.split(',').map(e => e.trim())

    const postData : PostData = { 
      name:        formData.name, 
      description: formData.description, 
      category:    formData.category, 
      tags:        tagArr,
      seeking:     seekingArr, 
      ownerId:     userID
    }
    try {
      const response: AxiosResponse<string> = await axiosBase.post<string>(`/project`, postData);
      const id = response.data;
      navigate(`/project/${id}`, {state: {_id: id, ...postData}})
    } catch (err) {
      setSubmitError("An error occurred")
      console.log(postData)
      console.log("hello error")
    }
  }

  return (
    <div>
      <div className='formBody'>
        <div className='inputContainer'>
          <h2> Create a new project</h2>
          <label className='formLabel'>Name:</label>
          <input type="text"  name="name" value={formData.name} onChange={handleChange}/>
          <div className='errorMessage'>
              {formFieldError.nameErr}
          </div>
      
        <label className='formLabel'>Description:</label>
          <textarea className='description' name="description" value={formData.description} onChange={handleChange}/>
          <div className='errorMessage'>
              {formFieldError.descriptionErr}
          </div>

          <label className='formLabel'>Category:</label>
          <select onChange={selectCategory}>
          <option value="" selected disabled hidden>Select one</option>
            {categories.map((value) => (
                <option key={value} value={value}> {value} </option>
              ))
            }
          </select>
          <div className='errorMessage'>
              {formFieldError.categoryErr}
          </div>
          
          <label className='formLabel'>Tags:</label>
          <input type="text" name="tags" placeholder="e.g., gardening, 3d printing, education" value={formData.tags} onChange={handleChange}/>
          <div className='errorMessage'>
              {formFieldError.tagsErr}
          </div>

          <label className='formLabel'>Roles wanted:</label>
          <input type="text" name="seeking" placeholder="e.g., sous chef, project manager, guitarist" value={formData.seeking} onChange={handleChange}/>
          <div className='errorMessage'>
              {formFieldError.seekingErr}
          </div>
          <div className='buttonContainer'>
            <button type="submit" onClick={onSubmit}>Create</button>
          </div>
          <div className='errorMessage'>
            {submitError}
          </div>
        </div>
      </div>
      {showPromptAccount &&
        <PromptAccount 
          setShowPromptAccount={setShowPromptAccount}
          setShowLogIn={setShowLogin}
          setShowSignUp={setShowSignUp}
          />
      }
      {
        showLogin &&
        <LogInModal SetShowLogIn={setShowLogin}/>
      }
      {
        showSignUp &&
        <SignUpModal setShowSignUp={setShowSignUp}/>
      }
    </div>

  )
}

export default CreateProject