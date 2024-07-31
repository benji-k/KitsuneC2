"use client"

import { useEffect, useState } from "react"

export default function Logs() {
    const [log, setLogs] = useState("")
    const [fetchError, setFetchError] = useState(false)

    useEffect(() => {
        fetch("/api/kitsune/logs").then((res) => (
            res.status === 200 ?
                res.json().then((json) => { setLogs(json) })
                :
                setFetchError(true)
        ))
    }, [])

    return (
        <>
            <h2 className="text-white text-3xl pl-5 pt-5">Logs</h2>
            <div className="bg-kc2-dark-gray mt-10 mx-4 p-4 rounded-md h-[600px] whitespace-pre-wrap overflow-scroll scrollbar-hide">
                {fetchError ?
                    <div className="text-red-300">Error fetching data from server</div>
                    :
                    <div className="text-white">
                        {log}
                    </div>
                }
            </div>
        </>
    )
}