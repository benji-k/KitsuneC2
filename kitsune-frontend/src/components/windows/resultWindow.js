import { useDashboardState } from "@/state/application"
import { Tasks } from "@/constants/tasks"

export default function ResultWindow() {
    const setResultWindowOpen = useDashboardState((state) => state.setResultWindowOpen)
    const taskResult = useDashboardState((state) => state.taskResult)

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
            const result = obj.Result 
            const error = obj.Error

            if (error){
                return "Error: " + error
            }

            if (result){
                return result
            }

            return "No output"
        } catch{
            return "No output"
        }
    }

    const getTaskName = (task) => {
        const filteredTasks = Tasks.filter((t) => t.taskType === task.Task_type);
        return filteredTasks[0] ? filteredTasks[0].taskName : "?";
    };

    return (
        <div className="fixed flex justify-center items-center top-0 left-0 h-full w-full bg-[#9B9B9B]/[.54] z-50">
            <div className="bg-kc2-dark-gray flex flex-col items-center w-full px-3 mx-5 pt-4 rounded-2xl max-w-6xl max-h-[600px] scrollbar-hide overflow-scroll">
            <h2 className="text-white mb-2 self-start mt-3 pb-2 text-xl w-full border-b-2 border-b-slate-200 border-opacity-30">Task info:</h2>
            <div className="text-white text-sm flex flex-col self-start">

                <div className="flex">
                        <p className="w-32 text-kc2-soap-pink">Task ID:</p>
                        <p>{taskResult.Task_id}</p>
                </div>
                <div className="flex">
                        <p className="w-32 text-kc2-soap-pink">Implant ID:</p>
                        <p>{taskResult.Implant_id}</p>
                </div>
                <div className="flex">
                    <p className="w-32 text-kc2-soap-pink">Task type:</p>
                    <p>{getTaskName(taskResult)}</p>
                </div>
                <div className="flex">
                    <p className="w-32 text-kc2-soap-pink">Task arguments:</p>
                    <p>{b64ToArguments(taskResult.Task_data)}</p>
                </div>    
            </div>

            <h2 className="text-white self-start mt-5 mb-4 pb-2 text-xl w-full border-b-2 border-b-slate-200 border-opacity-30">Task result:</h2>
            <div className="bg-black w-full rounded-lg min-h-[100px] mb-4 overflow-scroll scrollbar-hide text-white text-sm p-4 whitespace-pre-wrap">
                {
                    b64ToResults(taskResult.Task_result)
                }
            </div>
            
            <button className="bg-[#F96B6B] text-white rounded-md px-8 py-1 mb-3 md:self-end"
                onClick={() => { setResultWindowOpen(false) }}>Exit</button>
            </div>
        </div>
    )
}