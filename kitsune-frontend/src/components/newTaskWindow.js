"use client"

import { useState } from "react"
import { useDashboardState } from "@/state/dashboard"
import { Tasks } from "@/constants/tasks"

export default function NewTaskWindow() {
    const [selectedTask, setSelectedTask] = useState(Tasks.at(0).taskName)
    const setNewTaskWindowOpen = useDashboardState((state) => state.setNewTaskWindowOpen)
    const selectedImplants = useDashboardState((state) => state.selectedImplants)
    const pushNotification = useDashboardState((state) => state.pushNotification)

    const [taskArguments, setTaskArguments] = useState({})
    const [addingTask, setAddingTask] = useState(false)

    const handleArgumentChange = function(name, value){
        setTaskArguments((prevState) => ({
            ...prevState,
            [name] : value
        }))
    }

    const submitTask = async function(e){
        e.preventDefault()
        setAddingTask(true)
        const taskType = Tasks.find((t)=>(t.taskName === selectedTask)).taskType

        const postData = {
            "implants" : selectedImplants,
            "taskType" : taskType,
            "arguments" : taskArguments
        }

        try{
            const success = await fetch("/api/kitsune/tasks/add", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(postData),
            })
            if (success.status === 500){
                const err = await success.json()
                const errText = Object.values(err.error)[0]
                pushNotification({text: errText, type: "ERROR"})
            } else if (success.status === 200){
                setNewTaskWindowOpen(false)
            }
        } catch(e){
            pushNotification({text: e, type:"ERROR"})
        } finally{
            setAddingTask(false)
        }
        
    }

    return (
        <div className="fixed flex justify-center items-center top-0 left-0 h-full w-full bg-[#9B9B9B]/[.54] z-20">
            <form onSubmit={submitTask} className="bg-kc2-dark-gray flex flex-col items-center w-full px-3 mx-5 pt-4 rounded-2xl max-w-6xl max-h-[600px] overflow-scroll">

                <select className="bg-kc2-light-gray text-white text-center rounded-md w-full max-w-2xl"
                    onChange={(e) => {
                        setSelectedTask(e.target.value)
                        setTaskArguments({})
                        }}>
                    {
                        Tasks.map((t) => (
                            <option value={t.taskName} key={t.taskName}>{t.taskName}</option> //Don't ask me why, but setting value=t.taskType crashes app.
                        ))
                    }
                </select>

                <h2 className="text-white self-start mt-3 pb-2 text-xl w-full border-b-2 border-b-slate-200 border-opacity-30">Description:</h2>
                <p className="text-white self-start mb-3">{Tasks.filter((t) => (t.taskName === selectedTask)).at(0).description}</p>

                <h2 className="text-white self-start mt-3 pb-2 text-xl w-full border-b-2 border-b-slate-200 border-opacity-30">Arguments:</h2>

                <div className="bg-kc2-light-gray text-white w-full my-3 rounded-md pl-3 pb-3">
                    <p className="mb-3 mt-2 text-lg border-b-[1px] border-slate-400 pb-1">Required</p>
                    <div className="px-4 flex flex-col gap-2">
                        {Tasks.find((t) => (t.taskName === selectedTask)).args.map((arg) => (
                            !arg.optional &&
                                <Argument name={arg.name} tooltip={arg.tooltip} type={arg.type} key={arg.name} required={true} onChange={(val)=>{handleArgumentChange(arg.apiName, val)}} />
                        ))}
                    </div>
                    <p className="mb-3 mt-4 text-lg border-b-[1px] border-slate-400 pb-1">Optional</p>
                    <div className="px-4 flex flex-col gap-2">
                        {Tasks.find((t) => (t.taskName === selectedTask)).args.map((arg) => (
                            arg.optional &&
                                <Argument name={arg.name} tooltip={arg.tooltip} type={arg.type} key={arg.name} required={false} onChange={(val)=>{handleArgumentChange(arg.apiName, val)}}/>
                        ))}
                    </div>
                </div>

                <h2 className="text-white self-start mt-3 pb-2 text-xl w-full border-b-2 border-b-slate-200 border-opacity-30">Affected Implants:</h2>
                <div className="self-start text-kc2-soap-pink mb-4">
                    {selectedImplants.length === 0 ? 
                    <p>No implants selected.</p>
                     :
                    selectedImplants.map((i) => (
                        <p key={i}>{i}</p>
                    ))
                    }
                </div>

                <div className="flex w-full justify-center md:justify-end items-center gap-10 mb-4">
                    <button className="bg-[#F96B6B] text-white rounded-md px-8 py-1"
                        onClick={() => { setNewTaskWindowOpen(false) }}>Cancel</button>

                    {addingTask ? 
                        <div className="bg-[#abe4b1] text-white rounded-md px-8 py-1">Adding...</div>
                        :
                        <input type="submit" value="Execute" className="bg-[#0EC420] text-white rounded-md px-8 py-1"></input>
                    }
                    
                </div>
            </form>
        </div>
    )
}

function Argument({ name, type, tooltip, required, onChange }) {
    const handleChange = function(e){
        onChange(e.target.value)
    }

    const switchInputType = function (t) {
        switch (t) {
            case String:
                return "text"
            case Number:
                return "number"
            case Boolean:
                return "checkbox"
            case "file":
                return "file"
            default:
                return "text"
        }
    }

    const [showToolTip, setShowToolTip] = useState(false)

    return (
        <div className="relative">
            {showToolTip ?
                <div className="absolute bg-slate-900 p-2 text-sm rounded-lg bottom-4 left-8 ">
                    {tooltip}
                </div>
                :
                <></>
            }

            <div className="flex justify-between">
                <p onMouseEnter={() => { setShowToolTip(true) }} onMouseLeave={() => { setShowToolTip(false) }} className="w-32">{name}</p>
                <input onChange={handleChange} required={required} type={switchInputType(type)} className="bg-kc2-dark-gray cursor-pointer rounded-md outline-none grow overflow-scroll file:bg-slate-700 file:border-0 file:text-white"></input>
            </div>
        </div>
    )
}