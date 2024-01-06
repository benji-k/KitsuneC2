import Image from 'next/image'
import { FaUser } from 'react-icons/fa'
import { FaLock } from 'react-icons/fa'

export default function Login() {

    return (
        <div className="min-h-screen px-7 bg-gradient-to-br from-kc2-wine-purple from-10% to-zinc-600">
            <div className="max-w-lg mx-auto py-36 min-h-screen flex flex-col items-center">
                <Image src="/fox.png" width={120} height={120} alt='Kitsune logo' className='mb-36'></Image>
                <div className='flex flex-col gap-4 mb-16 lg:mb-40 w-full'>
                    <div className='flex gap-4 items-center mr-2'>
                        <FaUser size={17} color='#cad5e199'/>
                        <input placeholder='Username' type='text' 
                        className='w-full bg-transparent border-b-[1px] border-slate-300 border-opacity-60
                        text-slate-300 outline-none p-2'></input>
                    </div>
                    <div className='flex gap-4 items-center mr-2'>
                        <FaLock size={17} color='#cad5e199'/>
                        <input placeholder='Password' type='password' 
                        className='w-full bg-transparent border-b-[1px] border-slate-300 border-opacity-60
                        text-slate-300 outline-none p-2'></input>
                    </div>
                    <button className='bg-slate-100 w-full rounded-lg py-2 mt-4'>Login</button>
                </div>
            </div>
        </div>
    )
}