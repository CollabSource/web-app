import { ChangeEvent, FC, useContext, useState } from "react";
import { User } from "../../types/user";
import axiosBase from "../../config/axiosConfig";
import { AxiosHeaders, AxiosResponse } from "axios";
import './User.css'
import { AuthContext } from "../../context/AuthContext";

interface UserSettingsForm {
  name        : string
  email       : string
  description : string
  skills      : string
}

interface PostData  {
  name        : string
  email       : string
  description : string
  skills      : string[]
}

interface Props {
  user: User
}
const PrivateProfile:FC<Props> = ({ user }) => {
  const { token } = useContext(AuthContext)
  const [formData, setFormData] = useState<UserSettingsForm>({
    name       : user?.name || '',
    email: (user as User)?.email || '',
    description   : user?.description || '',
    skills       : user?.skills?.join(',') || '',
  })
  
  const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData({           
        ...formData,
        [name]: value,
      })
  }

  const onSubmit = async () => {
    // if (checkFormError(noErrorObj, formData, setFormFieldError)) {
    //   return
    // }
    const tagArr = formData.skills.split(',').map(e => e.trim())

    const postData : PostData = { 
      name: formData.name,
      email: formData.email,
      description: formData.description,
      skills: tagArr,
    }
    const headers = new AxiosHeaders()
    headers.Authorization = token
    console.log(postData)
    try {
      const response: AxiosResponse<string> = await axiosBase.patch<string>(`/user/${user._id}`, postData, {headers});
    } catch (err) {
      // setSubmitError("An error occurred")
      // console.log(postData)
      console.log(err)
    }
  }

  return (
  <div>
    <h2> Profile Settings </h2>
    <hr/>
    <div className="form-container">

      <label className='formLabel'>Name:</label>
      <input className='text-input' type="text"  name="name" value={formData.name} onChange={handleChange}/>
        {/* <div className='errorMessage'>
            {formFieldError.nameErr}
        </div> */}
      <label className='formLabel'>Email:</label>
      <input className='text-input' type="text"  name="email" value={formData.email} onChange={handleChange}/>

      <label className='formLabel'>About me:</label>
      <textarea className='text-textarea' name="description" value={formData.description} onChange={handleChange}/>

      <label className='formLabel'>Skills:</label>
      <input className='text-input' type="text"  name="skills" value={formData.skills} onChange={handleChange}/>
    </div>

    <div className='button-cont-profile'>
      <button type="submit" onClick={onSubmit}>Save</button>
    </div>
  </div>
  )
}

export default PrivateProfile