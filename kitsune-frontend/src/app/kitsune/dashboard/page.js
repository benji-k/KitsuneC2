import ImplantTable from "@/components/implantTable"
import NewTaskBtn from "@/components/newTaskBtn"
import TaskSelectBtn from "@/components/taskSelectBtn"
import TaskTable from "@/components/TaskTable"

export default function Dashboard() {

    return (
        <>
            <h2 className="text-white text-3xl pl-5 pt-5">Implants</h2>
            <ImplantTable 
                refreshRate={3000}
            />
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