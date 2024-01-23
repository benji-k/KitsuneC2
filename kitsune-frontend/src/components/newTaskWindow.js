"use client"

import { useState } from "react"
import { useDashboardState } from "@/state/dashboard"
import { Tasks } from "@/constants/tasks"

export default function NewTaskWindow() {
    const [selectedTask, setSelectedTask] = useState(Tasks.at(0).taskName)
    const setNewTaskWindowOpen = useDashboardState((state) => state.setNewTaskWindowOpen)
    const selectedImplants = useDashboardState((state) => state.selectedImplants)

    return (
        <div className="fixed flex justify-center items-center top-0 left-0 h-full w-full bg-[#9B9B9B]/[.54] z-50">
            <div className="bg-kc2-dark-gray flex flex-col items-center w-full px-3 mx-5 pt-4 rounded-2xl max-w-6xl max-h-[600px] overflow-scroll">

                <select className="bg-kc2-light-gray text-white text-center rounded-md w-full max-w-2xl"
                    onChange={(e) => setSelectedTask(e.target.value)}>
                    {
                        Tasks.map((t) => (
                            <option value={t.taskName} key={t.taskName}>{t.taskName}</option>
                        ))
                    }
                </select>

                <h2 className="text-white self-start mt-3 pb-2 text-xl w-full border-b-2 border-b-slate-200 border-opacity-30">Description:</h2>
                <p className="text-white self-start mb-3">{Tasks.filter((t) => (t.taskName === selectedTask)).at(0).description}</p>

                <h2 className="text-white self-start mt-3 pb-2 text-xl w-full border-b-2 border-b-slate-200 border-opacity-30">Arguments:</h2>

                <div className="bg-kc2-light-gray text-white w-full my-3 rounded-md pl-3 pb-3">
                    <p className="mb-3 mt-2 text-lg border-b-[1px] border-slate-400 pb-1">Required</p>
                    <div className="px-4 flex flex-col gap-2">
                        {Tasks.filter((t) => (t.taskName === selectedTask)).at(0).args.map((arg) => (
                            arg.optional ?
                                <></>
                                :
                                <Argument name={arg.name} tooltip={arg.tooltip} type={arg.type} key={arg.name} />
                        ))}
                    </div>
                    <p className="mb-3 mt-4 text-lg border-b-[1px] border-slate-400 pb-1">Optional</p>
                    <div className="px-4 flex flex-col gap-2">
                        {Tasks.filter((t) => (t.taskName === selectedTask)).at(0).args.map((arg) => (
                            arg.optional ?
                                <Argument name={arg.name} tooltip={arg.tooltip} type={arg.type} key={arg.name} />
                                :
                                <></>
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
                    <button className="bg-[#0EC420] text-white rounded-md px-8 py-1">Execute</button>
                </div>
            </div>
        </div>
    )
}

function Argument({ name, type, tooltip }) {
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
                <p onMouseEnter={() => { setShowToolTip(true) }} onMouseLeave={() => { setShowToolTip(false) }} className="w-24">{name}</p>
                <input type={switchInputType(type)} className="bg-kc2-dark-gray cursor-pointer rounded-md outline-none grow overflow-scroll file:bg-slate-700 file:border-0 file:text-white"></input>
            </div>
        </div>
    )
}