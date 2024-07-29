"use client"

import { useEffect, useState } from "react"

export default function Logs(){
    const [log, setLogs] = useState("")

    useEffect(() =>{
        const response = fetch("/api/kitsune/logs").then((res) => (
            res.json()).then((json) => {setLogs(json)})
        )
    }, [])

    return (
    <>
        <h2 className="text-white text-3xl pl-5 pt-5">Logs</h2>
        <div className="bg-kc2-dark-gray mt-10 mx-4 p-4 rounded-md h-[600px] whitespace-pre-wrap overflow-scroll scrollbar-hide text-white">
            {log}
        </div>
    </>
    )
}