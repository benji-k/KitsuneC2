"use client"

import { useState } from "react"
import { useDashboardState } from "@/state/application"
import { useGlobalState } from "@/state/application"
import { Tasks } from "@/constants/tasks"

export default function ResultWindow() {
    const setResultWindowOpen = useDashboardState((state) => state.setResultWindowOpen)
    const taskResult = useDashboardState((state) => state.taskResult)
    const pushNotification = useGlobalState((state) => state.pushNotification)
    const [downloading, setDownloading] = useState(false)

    const b64ToArguments = (b64Str) => {
        try{
            const obj = JSON.parse(atob(b64Str))
            delete obj.TaskId
            return JSON.stringify(obj)
        } catch{
            return "{}"
        }
    }

    const b64ToResults = (b64Str) => {
        try {
            const obj = JSON.parse(atob(b64Str))
            const error = obj.Error

            if (error){
                return "Error: " + error
            }

            delete obj.TaskId
            
            if (Object.keys(obj).length === 1){
                return obj[Object.keys(obj)[0]]
            } else if (Object.keys(obj).length > 1){
                return JSON.stringify(obj)
            } else {
                return "No output"
            }
        } catch{
            return "No output"
        }
    }

    const getTaskName = (task) => {
        const filteredTasks = Tasks.filter((t) => t.taskType === task.Task_type);
        return filteredTasks[0] ? filteredTasks[0].taskName : "?";
    };

    const getDownloadedFileName = function (task){
        try{
            const filePath = JSON.parse(b64ToArguments(task.Task_data)).Destination
            return filePath.replace(/^.*[\\/]/, '')
        } catch{
            return "test"
        }        
    }

    //For results that downloaded a file , we download it using this function. So the flow of the file is as follows:
    // implant -> K2C-server -> K2C-frontend -> client
    const downloadRemoteFile = async function(){
        setDownloading(true)
        try{
            const response = await fetch("/api/kitsune/file/download?taskId="+taskResult.Task_id)
            if (response.status === 200){
                const blob = await response.blob()
                const url = window.URL.createObjectURL(blob)
                const link = document.createElement("a")
                const fileName = getDownloadedFileName(taskResult)
                link.href = url
                link.setAttribute("download", fileName)
                document.body.appendChild(link)
                link.click()
                document.body.removeChild(link)
            } else{
                const err = await response.json()
                const errText = Object.values(err.error)
                pushNotification({text: errText, type:"ERROR"})
            }
        } catch (e){
            pushNotification({text: e.message, type:"ERROR"})
        } finally{
            setDownloading(false)
        }
    }

    return (
        <div className="fixed flex justify-center items-center top-0 left-0 h-full w-full bg-[#9B9B9B]/[.54] z-50">
            <div className="bg-kc2-dark-gray flex flex-col items-center w-full px-3 mx-5 pt-4 rounded-2xl max-w-6xl max-h-[600px] scrollbar-hide overflow-scroll">
            <h2 className="text-white mb-2 self-start mt-3 pb-2 text-xl w-full border-b-2 border-b-slate-200 border-opacity-30">Task info:</h2>
            <div className="text-white text-sm flex flex-col self-start max-w-full">
                <div className="flex gap-16">
                        <p className="text-kc2-soap-pink w-24 shrink-0">Task ID:</p>
                        <p className="overflow-scroll">{taskResult.Task_id}</p>
                </div>
                <div className="flex gap-16">
                        <p className="text-kc2-soap-pink w-24 shrink-0">Implant ID:</p>
                        <p className="overflow-scroll">{taskResult.Implant_id}</p>
                </div>
                <div className="flex gap-16">
                    <p className="text-kc2-soap-pink w-24 shrink-0">Task type:</p>
                    <p className="overflow-scroll">{getTaskName(taskResult)}</p>
                </div>
                <div className="flex gap-16">
                    <p className="text-kc2-soap-pink w-24 shrink-0">Task arguments:</p>
                    <p className="overflow-scroll max-h-60">{b64ToArguments(taskResult.Task_data)}</p>
                </div>    
            </div>

            <h2 className="text-white self-start mt-5 mb-4 pb-2 text-xl w-full border-b-2 border-b-slate-200 border-opacity-30">Task result:</h2>
            <div className="bg-black w-full rounded-lg min-h-[100px] mb-4 overflow-scroll scrollbar-hide text-white text-sm p-4 whitespace-pre-wrap">
                {
                    b64ToResults(taskResult.Task_result)
                }
            </div>
            
            <div className="flex">
                {taskResult.Task_type === 19 && b64ToResults(taskResult.Task_result).includes("Wrote file to") && <button className="bg-green-500 mr-4 text-white rounded-md px-8 py-1 mb-3 md:self-end"
                onClick={() => {downloadRemoteFile()}}>Download File
                </button>} 
                <button className="bg-[#F96B6B] text-white rounded-md px-8 py-1 mb-3 md:self-end"
                    onClick={() => { setResultWindowOpen(false) }}>Exit
                </button>
            </div>
            
            </div>
        </div>
    )
}