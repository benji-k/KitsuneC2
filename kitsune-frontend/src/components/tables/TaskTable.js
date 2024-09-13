"use client"

import useSWR from 'swr'
import ReactLoading from 'react-loading';
import { Tasks } from '@/constants/tasks';
import { useDashboardState } from '@/state/application';
import { useGlobalState } from "@/state/application"
import { FaTrash } from "react-icons/fa6";

export default function TaskTable({ refreshRate }) {
    const fetcher = async url => {
        const res = await fetch(url)

        if (!res.ok) {
            const error = new Error('An error occurred while fetching the data.')
            const errReason = await res.json()
            error.info = errReason.error
            error.status = res.status
            throw error
        }

        return res.json()
    }

    const pushNotification = useGlobalState((state) => state.pushNotification)
    const deletePendingTask = async function (taskId, implantId) {
        const formData = new FormData()
        formData.append("implantId", implantId)
        formData.append("taskId", taskId)

        try {
            const response = await fetch("/api/kitsune/tasks/remove", {
                method: "POST",
                body: formData,
            })

            if (response.status === 500) {
                const err = await response.json()
                const errText = Object.values(err.error)
                pushNotification({ text: errText, type: "ERROR" })
            }
        } catch (e) {
            pushNotification({ text: e, type: "ERROR" })
        }
    }

    const showCompletedTasks = useDashboardState((state) => (state.showCompletedTasks))
    const tasksDataUrl = showCompletedTasks ? "/api/kitsune/tasks?completed=true" : "/api/kitsune/tasks?completed=false"
    const selectedImplants = useDashboardState((state) => (state.selectedImplants))
    const setResultWindowOpen = useDashboardState((state) => (state.setResultWindowOpen))
    const setTaskResult = useDashboardState((state) => state.setTaskResult)

    const { data, error, isLoading } = useSWR(tasksDataUrl, fetcher, { refreshInterval: refreshRate })

    if (error) return (
        <div className="m-5 mt-3">
            <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll flex justify-center items-center text-lg scrollbar-hide">
                <p className='text-red-300'>Error fetching data from server: {error.info}</p>
            </div>
        </div>
    )

    if (isLoading) return (
        <div className="m-5 mt-3">
            <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll flex justify-center items-center scrollbar-hide">
                <ReactLoading type="spinningBubbles" color="#cccccc" height={75} width={75} />
            </div>
        </div>
    )

    if (data.filter((t) => selectedImplants.includes(t.Implant_id)).length === 0) return (
        <div className="m-5 mt-3">
            <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll flex justify-center items-center text-lg scrollbar-hide">
                <p className='text-slate-300'>No tasks for selected implant(s)</p>
            </div>
        </div>
    )



    const getTaskName = (task) => {
        const filteredTasks = Tasks.filter((t) => t.taskType === task.Task_type);
        return filteredTasks[0] ? filteredTasks[0].taskName : "?";
    };

    const b64ToArguments = (base64Arg) => {
        try {
            const jsonObj = JSON.parse(atob(base64Arg))
            delete jsonObj.TaskId
            return JSON.stringify(jsonObj)
        } catch {
            return "{}"
        }
    }

    if (data) return (
        <div className="m-5 mt-3">
            <div className="mx-auto h-80 bg-kc2-dark-gray overflow-scroll scrollbar-hide">

                <table className="table-auto w-full">
                    <thead>
                        <tr className="bg-kc2-light-gray h-10">
                            <th className="text-left text-xs text-white px-3">ID</th>
                            <th className="text-left text-xs text-white px-3">Module</th>
                            <th className="text-left text-xs text-white px-3">Arguments</th>
                            {showCompletedTasks && <th className="text-left text-xs text-white px-3">Result</th>}
                            {!showCompletedTasks && <th></th>}
                        </tr>
                    </thead>
                    <tbody>
                        {
                            data && data.filter((task) => selectedImplants.includes(task.Implant_id)).map((task) => (
                                <tr key={task.Task_id} className="odd:bg-kc2-dark-gray even:bg-black h-10">
                                    <td className="text-xs text-white px-3 whitespace-nowrap">{task.Task_id}</td>
                                    <td className="text-xs text-white px-3 whitespace-nowrap">{getTaskName(task)}</td>
                                    <td className="text-xs text-white px-3 whitespace-nowrap overflow-hidden max-w-64"><div className='overflow-scroll'>{b64ToArguments(task.Task_data)}</div></td>
                                    {showCompletedTasks &&
                                        <td className="text-xs text-kc2-soap-pink px-3 whitespace-nowrap cursor-pointer hover:underline"
                                            onClick={() => {
                                                setTaskResult(task)
                                                setResultWindowOpen(true)
                                            }}>Show result</td>}
                                    {!showCompletedTasks &&
                                        <td className='px-3'>
                                            <div>
                                                <FaTrash size={15} color='#F96B6B' className='cursor-pointer' onClick={() => (deletePendingTask(task.Task_id, task.Implant_id))} />
                                            </div>
                                        </td>
                                    }
                                </tr>
                            ))
                        }
                    </tbody>
                </table>
            </div>
        </div>
    )
}
