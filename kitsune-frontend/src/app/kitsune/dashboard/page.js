"use client"

import ImplantTable from "@/components/implantTable"
import NewTaskBtn from "@/components/newTaskBtn"
import TaskSelectBtn from "@/components/taskSelectBtn"
import TaskTable from "@/components/TaskTable"
import NewTaskWindow from "@/components/newTaskWindow"
import ResultWindow from "@/components/resultWindow"
import { useDashboardState } from "@/state/application"

export default function Dashboard() {
    const newTaskWindowOpen = useDashboardState((state) => state.newTaskWindowOpen)
    const resultWindowOpen = useDashboardState((state) => state.resultWindowOpen)
    return (
        <>
            {newTaskWindowOpen && <NewTaskWindow />}
            {resultWindowOpen && <ResultWindow />}       
            <h2 className="text-white text-3xl pl-5 pt-5">Implants</h2>
            <div className="m-5 mt-3">
                <ImplantTable 
                    refreshRate={3000}
                />
            </div>
            <div className="flex justify-between items-center">
                <div className="flex items-center">
                    <h2 className="text-white text-3xl pl-5 ">Tasks</h2>
                    <TaskSelectBtn />
                </div>
                <NewTaskBtn />
            </div>
            <TaskTable refreshRate={3000} />
        </>
    )
}