"use client";

import React, { useState, useContext, useEffect } from 'react';
import { API_URL } from '../../constants';
import { useRouter } from 'next/navigation';
import { AuthContext, UserInfo } from '@/modules/auth_provider';

const Index = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const { authenticated } = useContext(AuthContext);

  const router = useRouter(); // Correctly define router here.

  useEffect(() => {
    if (authenticated) {
      router.push('/');
      return;
    }
  }, [authenticated, router]); // Add router to the dependency array.

  const submitHandler = async (e: React.SyntheticEvent) => {
    e.preventDefault();

    try {
      const res = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      const data = await res.json();
      if (res.ok) {
        const user: UserInfo = {
          username: data.username,
          id: data.id,
        };

        localStorage.setItem('user_info', JSON.stringify(user));
        router.push('/'); // Correctly use router here.
      }
    } catch (err) {
      console.log(err);
    }
  };

  return (
    <div className='flex items-center justify-center min-w-full min-h-screen bg-gray-100'>
      <form className='flex flex-col md:w-1/5 bg-white p-6 rounded-lg shadow-lg'>
        <div className='text-3xl font-bold text-center'>
          <span className='text-blue-500'>welcome!</span>
        </div>

        <input
          placeholder='Email'
          className='p-3 mt-4 rounded-md border-2 border-gray-300 focus:outline-none focus:border-blue-500'
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />

        <input
          type='password'
          placeholder='Password'
          className='p-3 mt-4 rounded-md border-2 border-gray-300 focus:outline-none focus:border-blue-500'
          value={password}
          onChange={(p) => setPassword(p.target.value)}
        />

        <button
          className='p-3 mt-6 rounded-md bg-blue-500 font-bold text-white hover:bg-blue-600'
          type='submit'
          onClick={submitHandler}
        >
          Login
        </button>
      </form>
    </div>
  );
};

export default Index;
