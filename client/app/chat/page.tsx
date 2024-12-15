"use client";

import React, { useState, useRef, useContext, useEffect } from 'react'
import ChatBody from '@/components/chat_body';
import { useRouter } from 'next/navigation'
import { WebSocketContext } from '@/modules/websocket_provider';
import { AuthContext } from '@/modules/auth_provider';
import { API_URL } from '@/constants';
import autosize from 'autosize'

export type Message = {
    content: string
    client_id: string 
    username: string 
    room_id: string
    type: 'rece' | 'self'
}

const page = () => {
    const [messages, setMessages] = useState<Array<Message>>([])
    const textarea = useRef<HTMLTextAreaElement>(null)
    const { conn } = useContext(WebSocketContext)
    const [ users, setUsers ] = useState<Array<{ username: string }>>([])
    const { user } = useContext(AuthContext)
    const router = useRouter()

    // Get clients in the room
    useEffect(() => {
        if (conn === null) {
            router.push('/')
            return
        }

        const roomId = conn.url.split('/')[5]
        async function getUsers() {
            try {
                const res = await fetch(`${API_URL}/ws/getClients/${roomId}`, {
                    method: 'GET',
                    headers: { 'Content-Type': 'application/json'},
                })

                const data = await res.json()
                console.log('data: '+ JSON.stringify(data))

                setUsers(data)
            } catch(err) {
                console.log(err)
            }
        }

        getUsers()
    }, [])

    // Handle websocket connection
    useEffect(() => {
        if (textarea.current) {
            autosize(textarea.current)
        }

        if (conn === null) {
            router.push('/')
            return
        }

        // 4 different stages of connection 
        conn.onmessage = (message) => {
            // handle the message that receive from the backend
            const m: Message = JSON.parse(message.data)
            if (m.content == 'A new user has joined the room') {
                setUsers([...users, { username: m.username }])
            }

            if (m.content == 'user left the chat') {
                const deleteUser = users.filter((user) => user.username != m.username)
                setUsers([...deleteUser])
                setMessages([...messages, m])
                return
            }

            user?.username == m.username ? (m.type = 'self') : (m.type = 'rece')
            setMessages([...messages, m])

        }

        conn.onclose = () => {}
        conn.onerror = () => {}
        conn.onopen = () => {}

    }, [textarea, messages, conn, users])

    const sendMessage = () => {
        if (!textarea.current?.value) return

        // check connection
        if (conn === null) {
            router.push('/')
            return
        }

        conn.send(textarea.current.value)
        textarea.current.value = ''
    }

    return (
        <div className='flex flex-col w-full'>
            <div className='p-4 md:mx-6 mb-14'>
            <ChatBody data={messages}/>
            </div>
            <div className='fixed bottom-0 mt-4 w-full'>
                <div className='flex md:flex-row px-4 py-2 bg-grey md:mx-4 rounded-md'>
                    <div className='flex w-full mr-4 rounded-md border border-blue-500'>
                        <textarea 
                            ref={textarea}
                            placeholder='type your message here'
                            className='w-full h-10 p-2 rounded-md focus:outline-none'
                            style={{ resize: 'none' }}
                        />
                    </div>
                    <div className='flex items-center'>
                        <button className='p-2 rounded-md bg-blue-500 text-white' onClick={sendMessage}>
                            Send
                        </button>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default page