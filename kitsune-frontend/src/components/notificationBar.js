"use client"

import { useGlobalState } from "@/state/application"
import { useEffect, useState } from "react"
import useSWR from 'swr'

export default function NotificationBar({ popupTime, refreshRate }) {
    const notificationQueue = useGlobalState((state) => state.notificationQueue)
    const popNotification = useGlobalState((state) => state.popNotification)
    const pushNotification = useGlobalState((state) => state.pushNotification)
    const [notificationQueueState, setNotificationQueueState] = useState([])

    useEffect(() => {
        if (notificationQueue.length > 0 && notificationQueueState.length === 0){
            // If there are notifications in the queue and the component state queue is empty,
            // pop the notification from the Zustand store and set it to the component state queue
            const notification = popNotification()
            setNotificationQueueState([notification])
        }
    }, [notificationQueue, notificationQueueState, popNotification])

    useEffect(() => {
        if (notificationQueueState.length > 0){
            // If there are notifications in the component state queue,
            // remove the first notification after popupTime seconds
            const timer = setTimeout(() => {
                setNotificationQueueState((prevQueue) => prevQueue.slice(1))
            }, popupTime)

            return () => clearTimeout(timer)
        }
    }, [notificationQueueState, popupTime])

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

    const { data, error, isLoading} = useSWR('/api/kitsune/notifications', fetcher, { refreshInterval: refreshRate })
    useEffect(() =>{
        if (data){
            data.forEach(notification => {
                pushNotification({text: notification.Message, type: notification.NType})
            });
        }
    }, [data, pushNotification])
    


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