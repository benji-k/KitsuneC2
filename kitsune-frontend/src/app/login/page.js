"use client"

import Image from 'next/image'
import { FaUser } from 'react-icons/fa'
import { FaLock } from 'react-icons/fa'
import { RxCrossCircled } from "react-icons/rx";
import { signIn } from 'next-auth/react'
import { useState } from 'react'
import { redirect, useSearchParams } from 'next/navigation'
import { useSession } from 'next-auth/react';

export default function Login() {
    const { data: session, status } = useSession()
    if(status === "authenticated"){
        redirect("/kitsune/dashboard")
    }

    const [submitting, setSubmitting] = useState(false)
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")


    return (
        <div className="min-h-screen px-7 bg-gradient-to-br from-kc2-wine-purple from-10% to-zinc-600">
            <div className="max-w-lg mx-auto py-36 min-h-screen flex flex-col items-center">
                <Image src="/fox.png" width={120} height={120} alt='Kitsune logo' className='mb-36'></Image>
                <form onSubmit={(e) => {
                    e.preventDefault()
                    signIn("credentials", {username: username, password: password, callbackUrl: "/kitsune/dashboard"} )     
                    }} className='flex flex-col gap-4 mb-16 lg:mb-40 w-full'>

                    <div className='flex gap-4 items-center mr-2'>
                        <FaUser size={17} color='#cad5e199'/>
                        <input onChange={e => setUsername(e.currentTarget.value)} placeholder='Username' type='text' 
                        className='w-full bg-transparent border-b-[1px] border-slate-300 border-opacity-60
                        text-slate-300 outline-none p-2'></input>
                    </div>
                    <div className='flex gap-4 items-center mr-2'>
                        <FaLock size={17} color='#cad5e199'/>
                        <input onChange={e => setPassword(e.currentTarget.value)} placeholder='Password' type='password' 
                        className='w-full bg-transparent border-b-[1px] border-slate-300 border-opacity-60
                        text-slate-300 outline-none p-2'></input>
                    </div>

                    <button type="submit" onClick={() => {
                        setSubmitting(true)
                    }} 
                    className={`${submitting ? "bg-slate-300" : "bg-slate-100"} rounded-lg py-2 mt-4`}>{submitting ? "Logging in..." : "Login"}
                    </button>

                    {useSearchParams().get("error") != null ? 
                        <div className='bg-red-400 flex items-center justify-center rounded-lg py-2 gap-1 text-white text-sm my-5 fade'>
                            <RxCrossCircled size={15} color='white'/>
                            Authentication failure
                        </div> 
                    : 
                    <></>}
                </form>
            </div>
        </div>
    )
}