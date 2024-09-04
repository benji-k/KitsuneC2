"use client"

import useSWR from 'swr'
import ReactLoading from 'react-loading';
import { FaTrash } from "react-icons/fa6";
import { useGlobalState } from "@/state/application"

export default function ListenerTable({ refreshRate }) {
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

    const pushNotification = useGlobalState((state) => state.pushNotification)
    const deleteListener = async function (id) {
        const postData = {
            "id": id,
        }

        try {
            const response = await fetch("/api/kitsune/listeners/remove", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(postData),
            })

            if (response.status === 500) {
                const err = await response.json()
                const errText = Object.values(err.error)
                pushNotification({ text: errText, type: "ERROR" })
            }
        } catch (e) {
            pushNotification({ text: e, type: "ERROR" })
        }
    }

    const { data, error, isLoading } = useSWR("/api/kitsune/listeners", fetcher, { refreshInterval: refreshRate })

    if (error) return (
        <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll scrollbar-hide flex justify-center items-center text-lg">
            <p className='text-red-300'>Error fetching data from server: {error.info}</p>
        </div>
    )

    if (isLoading) return (
        <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll scrollbar-hide flex justify-center items-center">
            <ReactLoading type="spinningBubbles" color="#cccccc" height={75} width={75} />
        </div>
    )

    if (data.length === 0) return (
        <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll scrollbar-hide flex justify-center items-center text-lg">
            <p className='text-slate-300'>No active listeners</p>
        </div>
    )

    if (data) return (
        <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll scrollbar-hide">
            <table className="table-auto w-full">
                <thead>
                    <tr className="bg-kc2-dark-gray h-10">
                        <th className="text-left text-xs text-white px-3">ID</th>
                        <th className="text-left text-xs text-white px-3">Network Interface</th>
                        <th className="text-left text-xs text-white px-3">Port</th>
                        <th className="text-left text-xs text-white px-3">Type</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    {
                        data && data.map((listener, index) => (
                            <tr key={index} className="odd:bg-kc2-light-gray even:bg-[#434254] h-10">
                                <td className="text-xs text-white px-3 whitespace-nowrap">{index}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{listener.Network}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{listener.Port}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{listener.Type}</td>
                                <td className='px-3'>
                                    <div>
                                        <FaTrash size={15} color='#F96B6B' className='cursor-pointer' onClick={() => (deleteListener(index))} />
                                    </div>
                                </td>
                            </tr>
                        ))
                    }
                </tbody>
            </table>
        </div>
    )

}