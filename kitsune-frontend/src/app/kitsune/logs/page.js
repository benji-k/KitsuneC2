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

    const downloadLog = function () {
        const blobLog = new Blob([log], {
            type: 'text/plain'
        })
        const url = window.URL.createObjectURL(blobLog)
        const link = document.createElement("a")
        const fileName = "log.txt"
        link.href = url
        link.setAttribute("download", fileName)
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
    }

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
            {fetchError ?
                <div className="bg-slate-300 max-w-max p-3 px-7 m-4 rounded-lg text-white">Download</div>
                :
                <button onClick={downloadLog} className="bg-green-600 p-3 px-7 m-4 rounded-lg text-white">Download</button>
            }
        </>
    )
}