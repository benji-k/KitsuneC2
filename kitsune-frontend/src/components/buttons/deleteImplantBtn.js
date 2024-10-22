"use client"

import { useDashboardState } from "@/state/application"
import { useGlobalState } from "@/state/application"
import { useState } from "react"

export default function DeleteImplantBtn() {
    const selectedImplants = useDashboardState((state) => state.selectedImplants)
    const setSelectedImplants = useDashboardState((state) => state.setSelectedImplants)
    const setConfirmationWindowOpen = useDashboardState((state) => state.setConfirmationWindowOpen)
    const setConfirmationWindowText = useDashboardState((state) => state.setConfirmationWindowText)
    const setConfirmationWindowCallbackFn = useDashboardState((state) => state.setConfirmationWindowCallbackFn)
    const pushNotification = useGlobalState((state) => state.pushNotification)
    const [showToolTip, setShowToolTip] = useState(false)
    const [loading, setLoading] = useState(false)


    const deleteImplant = async function () {
        setLoading(true)
        const formData = new FormData()
        formData.append("implants", selectedImplants)

        try {
            const request = await fetch("/api/kitsune/implants/remove", {
                method: "POST",
                body: formData,
            })
            if (request.status === 500) {
                const err = await request.json()
                const errText = JSON.stringify(err)
                pushNotification({ text: errText, type: "ERROR" })
            } else if (request.status === 200){
                setSelectedImplants([])
            }
        } catch (e) {
            pushNotification({ text: e.message, type: "ERROR" })
        } finally {
            setLoading(false)
        }
    }

    const confirmDeleteImplant = function () {
        setConfirmationWindowText("This removes the implant and permanently deletes all (pending) tasks and their associated results. It is up to you to kill the implant first before removing it. Are you sure?")
        setConfirmationWindowCallbackFn(deleteImplant)
        setConfirmationWindowOpen(true)
    }

    return (

        <div className="relative text-white text-sm z-10" onMouseLeave={() => { setShowToolTip(false) }} onMouseEnter={() => { selectedImplants.length === 0 && setShowToolTip(true) }}>
            {showToolTip ?
                <div className="absolute bg-slate-900 p-2 text-sm rounded-lg bottom-6 right-8 w-52">
                    Please select at least 1 implant
                </div>
                :
                <></>
            }

            {
                loading ?
                    <div className="bg-red-300 rounded-md mr-5 p-2 px-6">Deleting...</div>
                :
                    <button className={`${selectedImplants.length === 0 ? "bg-slate-300" : "bg-red-500"} rounded-md mr-5 p-2
                    px-6`} onClick={() => { confirmDeleteImplant() }} disabled={selectedImplants.length === 0}
                    >Delete
                    </button>
            }

        </div>

    )
}