
export default function PendingTaskTable({pendingTasksData}){

    return(
        <div className="m-5 mt-3">
        <div className="mx-auto h-80 bg-kc2-dark-gray overflow-scroll">

            <table className="table-auto w-full">
                <thead>
                    <tr className="bg-kc2-light-gray h-10">
                        <th className="text-left text-xs text-white px-3">ID</th>
                        <th className="text-left text-xs text-white px-3">Module</th>
                        <th className="text-left text-xs text-white px-3">Arguments</th>
                    </tr>
                </thead>
                <tbody>
                    {
                        pendingTasksData && pendingTasksData.map((task) => (
                            <tr key={task.id} className="odd:bg-kc2-dark-gray even:bg-black h-10">
                                <td className="text-xs text-white px-3 whitespace-nowrap">{task.id}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{task.module}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{task.arguments}</td>
                            </tr>
                        ))
                    }
                </tbody>
            </table>
        </div>
        </div>
    )
}