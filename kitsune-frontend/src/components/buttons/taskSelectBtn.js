'use client'
import { useDashboardState } from "@/state/application"

export default function TaskSelectBtn() {
    const showCompletedTasks = useDashboardState((state) => state.setShowCompletedTasks)

    return (
        <select className="bg-kc2-light-gray text-white px-3 ml-4 py-1 rounded-md" onChange={(e) => 
        e.target.value === "completed" ? showCompletedTasks(true) : showCompletedTasks(false)}>
            <option value="pending">Pending</option>
            <option value="completed">Completed</option>
        </select>
    )
}