import Navbar from "@/components/navbar"
import ImplantTable from "@/components/implantTable"
import NewTaskBtn from "@/components/newTaskBtn"
import TaskSelectBtn from "@/components/taskSelectBtn"
import PendingTaskTable from "@/components/pendingTaskTable"


export default function Dashboard() {
    const testImplant1 = {
        "id": "e358efa489f58062f10dd7316b65649e",
        "externalAddress": "192.168.1.200",
        "hostname": "Debian",
        "username": "root",
        "uid": "0",
        "gid": "0",
        "os": "Linux",
        "lastSeen": "20013"
    }
    const testImplant2 = {
        "id": "a358efa489f58062f10dd7316b65649e",
        "externalAddress": "192.168.1.200",
        "hostname": "Debian",
        "username": "root",
        "uid": "0",
        "gid": "0",
        "os": "Linux",
        "lastSeen": "20013"
    }
    const testTask1 = {
        "id": "a358efa489f58062f10dd7316b65649e",
        "module": "exec",
        "arguments": "bash -c echo pwned"
    }
    const testTask2 = {
        "id": "v358efa489f58062f10dd7316b65649e",
        "module": "exec",
        "arguments": "bash -c echo pwned"
    }


    return (
        <>
            <h2 className="text-white text-3xl pl-5 pt-5">Implants</h2>
            <ImplantTable
                implantsData={[testImplant1, testImplant2]}
            />
            <div className="flex justify-between items-center">
                <div className="flex items-center">
                    <h2 className="text-white text-3xl pl-5 ">Tasks</h2>
                    <TaskSelectBtn />
                </div>
                <NewTaskBtn />
            </div>
            <PendingTaskTable pendingTasksData={[testTask1, testTask2]} />
        </>
    )
}