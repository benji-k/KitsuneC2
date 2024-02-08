"use client"

import { useDashboardState } from "@/state/dashboard"
import { useEffect, useState } from "react"

export default function NotificationBar({ popupTime }) {
    const notificationQueue = useDashboardState((state) => state.notificationQueue)
    const popNotification = useDashboardState((state) => state.popNotification)
    const [notificationQueueState, setNotificationQueueState] = useState([])

    useEffect(() => {
        if (notificationQueue.length > 0 && notificationQueueState.length === 0){
            // If there are notifications in the queue and the component state queue is empty,
            // pop the notification from the Zustand store and set it to the component state queue
            const notification = popNotification()
            setNotificationQueueState([notification])
        }
    }, [notificationQueue, notificationQueueState])

    useEffect(() => {
        if (notificationQueueState.length > 0){
            // If there are notifications in the component state queue,
            // remove the first notification after 5 seconds
            const timer = setTimeout(() => {
                setNotificationQueueState((prevQueue) => prevQueue.slice(1))
            }, popupTime)

            return () => clearTimeout(timer)
        }
    }, [notificationQueueState])


    const currentNotification = notificationQueueState[0]

    const getNotificationColor = function(notificationType){
        switch(notificationType){
            case "ERROR":
                return "bg-red-400"
            case "INFO":
                return "bg-blue-400"
            case "SUCCESS":
                return "bg-green-400"
            default:
                return "bg-blue-400"
        }
    }
    
    return (
        <div className={`fixed w-full px-5 bottom-16 z-50 transition-opacity duration-300 ${notificationQueueState.length > 0 ? 'opacity-100' : 'opacity-0'}`}>
            {notificationQueueState.length > 0 && (
                <div className={`${getNotificationColor(currentNotification.type)} rounded-md py-2 px-2 flex justify-center`}>
                    <p className="text-white font-semibold text-center">{currentNotification && currentNotification.text}</p>
                </div>
            )}
        </div>
    )
}