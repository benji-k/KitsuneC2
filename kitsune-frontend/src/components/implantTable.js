"use client"

import useSWR from 'swr'
import ReactLoading from 'react-loading';
import { useDashboardState } from '@/state/application';

export default function ImplantTable({ refreshRate }) {
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
    const selectImplant = useDashboardState((state) => (state.selectImplant))
    const selectedImplants = useDashboardState((state) => (state.selectedImplants))

    const { data, error, isLoading } = useSWR('/api/kitsune/implants', fetcher, { refreshInterval: refreshRate })

    if (error) return (
        <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll flex justify-center items-center text-lg">
            <p className='text-red-300'>Error fetching data from server: {error.info}</p>
        </div>
    )

    if (isLoading) return (
        <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll flex justify-center items-center">
            <ReactLoading type="spinningBubbles" color="#cccccc" height={75} width={75} />
        </div>
    )

    if (data.length === 0) return (
        <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll flex justify-center items-center text-lg">
            <p className='text-slate-300'>No active implants</p>
        </div>
    )

    if (data) return (
        <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll">
            <table className="table-auto w-full">
                <thead>
                    <tr className="bg-kc2-dark-gray h-10">
                        <th></th>
                        <th className="text-left text-xs text-white px-3">ID</th>
                        <th className="text-left text-xs text-white px-3">External Address</th>
                        <th className="text-left text-xs text-white px-3">Hostname</th>
                        <th className="text-left text-xs text-white px-3">Username</th>
                        <th className="text-left text-xs text-white px-3">UID</th>
                        <th className="text-left text-xs text-white px-3">GID</th>
                        <th className="text-left text-xs text-white px-3">OS</th>
                        <th className="text-left text-xs text-white px-3">Arch</th>
                        <th className="text-left text-xs text-white px-3">Last Seen</th>
                    </tr>
                </thead>
                <tbody>
                    {
                        data && data.map((implant) => (
                            <tr key={implant.Id} className="odd:bg-kc2-light-gray even:bg-[#434254] h-10 hover:bg-slate-500 cursor-pointer"
                             onClick={()=>{selectImplant(implant.Id)}}>
                                <td className='px-3'><input type='checkbox' onChange={()=>{}} checked={selectedImplants.includes(implant.Id)}></input></td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.Id}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.Public_ip}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.Hostname}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.Username}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.Uid}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.Gid}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.Os}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.Arch}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{Math.floor(Date.now() / 1000) - parseInt(implant.Last_checkin)} seconds ago</td>
                            </tr>
                        ))
                    }
                </tbody>
            </table>
        </div>
    )
}