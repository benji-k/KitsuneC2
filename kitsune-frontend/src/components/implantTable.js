
export default function ImplantTable({ implantsData }) {

    return (
        <div className="m-5 mt-3">
        <div className="mx-auto h-64 bg-kc2-light-gray overflow-scroll">

            <table className="table-auto w-full">
                <thead>
                    <tr className="bg-kc2-dark-gray h-10">
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
                        implantsData && implantsData.map((implant) => (
                            <tr key={implant.id} className="odd:bg-kc2-light-gray even:bg-[#434254] h-10 hover:bg-slate-500 cursor-pointer">
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.id}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.externalAddress}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.hostname}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.username}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.uid}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.gid}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.os}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.arch}</td>
                                <td className="text-xs text-white px-3 whitespace-nowrap">{implant.lastSeen} seconds ago</td>
                            </tr>
                        ))
                    }
                </tbody>
            </table>
        </div>
        </div>
    )
}