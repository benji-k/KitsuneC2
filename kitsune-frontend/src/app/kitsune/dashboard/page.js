"use client"

import ImplantTable from "@/components/tables/implantTable"
import NewTaskBtn from "@/components/buttons/newTaskBtn"
import DeleteImplantBtn from "@/components/buttons/deleteImplantBtn"
import TaskSelectBtn from "@/components/buttons/taskSelectBtn"
import TaskTable from "@/components/tables/TaskTable"
import NewTaskWindow from "@/components/windows/newTaskWindow"
import ResultWindow from "@/components/windows/resultWindow"
import ConfirmationWindow from "@/components/windows/confirmationWindow"
import { useDashboardState } from "@/state/application"

export default function Dashboard() {
    const newTaskWindowOpen = useDashboardState((state) => state.newTaskWindowOpen)
    const resultWindowOpen = useDashboardState((state) => state.resultWindowOpen)
    const confirmationWindowOpen = useDashboardState((state) => state.confirmationWindowOpen)
    return (
        <>
            {newTaskWindowOpen && <NewTaskWindow />}
            {resultWindowOpen && <ResultWindow />}    
            {confirmationWindowOpen && <ConfirmationWindow />}   
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
                <div className="flex">
                    <DeleteImplantBtn />
                    <NewTaskBtn />
                </div>
            </div>
            <TaskTable refreshRate={3000} />
        </>
    )
}