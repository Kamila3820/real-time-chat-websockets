import React from 'react'

const index = () => {
  return (
    <div className='flex item-center justify-center min-w-full min-h-screen'>
        <form className='flex flex-col md:w-1/5'>
            <input 
              placeholder='email' 
              className='p-3 mt-8 rouded-md border-2 border-border-grey focus:outline-none focus:border-blue'/>

            <input 
              type='password' 
              placeholder='password' 
              className='p-3 mt-8 rouded-md border-2 border-border-grey focus:outline-none focus:border-blue'/>

            <button>login</button>
        </form>
    </div>
  )
}

export default index