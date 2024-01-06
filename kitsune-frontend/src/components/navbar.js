import Image from 'next/image'
import Link from 'next/link'
import { FaUser } from 'react-icons/fa'
import { IoMdArrowDropdown } from "react-icons/io";

import {NavLinks} from '../constants/navbar'

export default function Navbar(){
    
    return (
        <>
            <div className="bg-kc2-light-gray px-10 py-3">
                <div className='flex justify-between items-center'>
                    <div className='flex items-center gap-3'>
                        <Image src="/fox.png" width={40} height={40} alt='Kitsune logo'></Image>
                        <h1 className='text-kc2-soap-pink text-2xl hidden md:block'>KitsuneC2</h1>
                    </div>

                    <div className='flex gap-16 text-[#ABABAB]'>
                        <div className='hidden md:flex gap-3 items-center'>
                            <FaUser size={15} color='#ABABAB'/>
                            <p>Welcome, USER</p>
                        </div>
                        <button>Logout</button>
                    </div>
                </div>
            </div>
            <div className="bg-kc2-dark-gray py-3 px-10  overflow-scroll border-b-2 border-black">
                <nav className='flex gap-8 justify-center md:justify-end'>
                {NavLinks.map((link) => (
                    <li key={link.label} className='list-none'>
                        <Link href={link.href} className='text-white'>{link.label}</Link>
                    </li>
                ))}
                </nav>
            </div>
        </>
    )
}