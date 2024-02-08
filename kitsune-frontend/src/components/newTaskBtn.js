"use client"

import { useDashboardState } from "@/state/dashboard"
import { useState } from "react"

export default function NewTaskBtn(){
    const setNewTaskWindowOpen = useDashboardState((state) => state.setNewTaskWindowOpen)
    const selectedImplants = useDashboardState((state) => state.selectedImplants)
    const [showToolTip, setShowToolTip] = useState(false) 

    return(
        <div className="relative text-white text-sm z-10" onMouseLeave={()=>{setShowToolTip(false)}} onMouseEnter={() =>{selectedImplants.length === 0 && setShowToolTip(true)}}>
            {showToolTip ?
                <div className="absolute bg-slate-900 p-2 text-sm rounded-lg bottom-6 right-8 w-52">
                Please select at least 1 implant
                </div>
             :
                <></>
            }
            
            <button className={`${selectedImplants.length === 0 ? "bg-slate-300":"bg-green-600"} rounded-md mr-5 p-2
             px-6`} onClick={()=>{setNewTaskWindowOpen(true)}} disabled={selectedImplants.length === 0}
             >+ Add Task</button>
        </div>
    )
}