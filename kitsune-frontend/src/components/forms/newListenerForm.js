"use client"

import { useState } from "react"
import { useGlobalState } from "@/state/application"

export default function NewListenerForm() {
    const [newListenerNetwork, setNewListenerNetwork] = useState("0.0.0.0")
    const [newListenerPort, setNewListenerPort] = useState("4444")
    const [newListenerType, setNewListenerType] = useState("TCP")
    const pushNotification = useGlobalState((state) => state.pushNotification)

    const AddListener = async function(e){
        e.preventDefault()
        
        const formData = new FormData()
        formData.append("network", newListenerNetwork)
        formData.append("port", newListenerPort)
        formData.append("type", newListenerType)

        try{
            const response = await fetch("/api/kitsune/listeners/add", {
                method: "POST",
                body: formData,
            })
            if (response.status === 500){
                const err = await response.json()
                const errText = Object.values(err.error)
                pushNotification({text: errText, type:"ERROR"})
            } else if (response.status === 200){
                pushNotification({text: "Successfully added listener", type:"SUCCESS"})
            }
        } catch(e){
            pushNotification({text: e.message, type:"ERROR"})
        } 
    }

    return (
        <div className="bg-kc2-dark-gray p-3 max-w-md rounded-lg">
            <form onSubmit={AddListener}>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">Network</p>
                        <input type="text" placeholder="0.0.0.0" className="bg-kc2-light-gray rounded-md outline-none text-white pl-1" onChange={(e) => {setNewListenerNetwork(e.target.value)}}></input>
                    </div>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">Port</p>
                        <input type="number" placeholder="4444" required className="bg-kc2-light-gray rounded-md outline-none text-white pl-1" onChange={(e) => {setNewListenerPort(e.target.value)}}></input>
                    </div>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">Type</p>
                        <select className="bg-kc2-light-gray text-white text-center rounded-md p-1" onChange={(e) => {setNewListenerType(e.target.value)}}>
                            <option>TCP</option>
                        </select>
                    </div>
                    
                    <input type="submit" value="Add Listener" className="bg-[#0EC420] text-white rounded-md px-8 py-1 mt-3"></input>
            </form>
        </div>
    )
}