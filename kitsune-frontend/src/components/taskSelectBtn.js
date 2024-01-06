'use client'

import { IoMdArrowDropdown } from "react-icons/io";

export default function TaskSelectBtn(){

    return(
        <button className="bg-kc2-light-gray text-white px-3 ml-4
        rounded-md">
            <div className="flex items-center gap-3">
                <p>Pending</p>
                <IoMdArrowDropdown />
            </div>
        </button>
    )
}