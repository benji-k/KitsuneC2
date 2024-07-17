"use client"

import { useState } from "react"
import { useGlobalState } from "@/state/application"

export default function GenImplantForm() {
    const [implantOs, setImplantOs] = useState("linux")
    const [implantArch, setImplantArch] = useState("amd64")
    const [serverIp, setServerIp] = useState("127.0.0.1")
    const [serverPort, setServerPort] = useState(4444)
    const [implantName, setImplantName] = useState("")
    const [cbInterval, setCbInterval] = useState("10")
    const [cbJitter, setCbJitter] = useState("2")
    const [maxRetryCount, setMaxRetryCount] = useState("10")
    const [generating, setGenerating] = useState(false)
    const pushNotification = useGlobalState((state) => state.pushNotification)

    const handleFormSubmit = async function (e) {
        e.preventDefault()

        const postData = {
            "os": implantOs,
            "arch": implantArch,
            "serverIp": serverIp,
            "serverPort": serverPort,
            "name": implantName,
            "cbInterval": cbInterval,
            "cbJitter": cbJitter,
            "maxRetryCount": maxRetryCount
        }

        setGenerating(true)

        try{
            const response = await fetch("/api/kitsune/implants/generate", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(postData),
            })
    
            if (response.status === 200){
                const blob = await response.blob()
                const url = window.URL.createObjectURL(blob)
                const link = document.createElement("a")
                const fileName = `implant_${implantOs}_${implantArch}${implantOs === "windows" ? ".exe" : ""}`
                link.href = url
                link.setAttribute("download", fileName)
                document.body.appendChild(link)
                link.click()
                document.body.removeChild(link)
            } else{
                const err = await response.json()
                const errText = Object.values(err.error)[0]
                pushNotification({text: errText, type:"ERROR"})
            }
        } catch(e){
            pushNotification({text: e, type:"ERROR"})
        } finally{
            setGenerating(false)
        }
    }


    return (
        <div className="flex gap-x-48 gap-y-10 flex-wrap">
            <div className="bg-kc2-dark-gray p-3 max-w-md min-w-96 rounded-lg">
                <form onSubmit={handleFormSubmit}>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">OS</p>
                        <select className="bg-kc2-light-gray text-white text-center rounded-md p-1" onChange={(e) => { setImplantOs(e.target.value) }}>
                            <option>linux</option>
                            <option>windows</option>
                            <option>aix</option>
                            <option>android</option>
                            <option>darwin</option>
                            <option>dragonfly</option>
                            <option>freebsd</option>
                            <option>illumos</option>
                            <option>ios</option>
                            <option>js</option>
                            <option>netbsd</option>
                            <option>plan9</option>
                            <option>solaris</option>
                        </select>
                    </div>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">Architecture</p>
                        <select className="bg-kc2-light-gray text-white text-center rounded-md p-1" onChange={(e) => { setImplantArch(e.target.value) }}>
                            <option>amd64</option>
                            <option>386</option>
                            <option>arm</option>
                            <option>arm64</option>
                            <option>mips</option>
                            <option>mips64le</option>
                            <option>mipsle</option>
                            <option>ppc64</option>
                            <option>ppc64le</option>
                            <option>riscv64</option>
                            <option>s390x</option>
                            <option>wasm</option>
                        </select>
                    </div>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">Callback IP</p>
                        <input required title="Test tooltip" type="text" value={serverIp} className="bg-kc2-light-gray rounded-md outline-none text-white pl-1" onChange={(e) => { setServerIp(e.target.value) }}></input>
                    </div>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">Callback Port</p>
                        <input required type="number" value={serverPort} className="bg-kc2-light-gray rounded-md outline-none text-white pl-1" onChange={(e) => { setServerPort(e.target.value) }}></input>
                    </div>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">Implant Name</p>
                        <input type="text" placeholder="EvilImplant" className="bg-kc2-light-gray rounded-md outline-none text-white pl-1" onChange={(e) => { setImplantName(e.target.value) }}></input>
                    </div>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">Callback Interval</p>
                        <input required type="number" value={cbInterval} className="bg-kc2-light-gray rounded-md outline-none text-white pl-1" onChange={(e) => { setCbInterval(e.target.value) }}></input>
                    </div>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">Callback Jitter</p>
                        <input required type="number" value={cbJitter} className="bg-kc2-light-gray rounded-md outline-none text-white pl-1" onChange={(e) => { setCbJitter(e.target.value) }}></input>
                    </div>
                    <div className="flex justify-between py-2 max-w-96">
                        <p className="text-white">Reconnect Try Count</p>
                        <input required type="number" value={maxRetryCount} className="bg-kc2-light-gray rounded-md outline-none text-white pl-1" onChange={(e) => { setMaxRetryCount(e.target.value) }}></input>
                    </div>
                    {
                        generating ? 
                        <p className="bg-[#6c6c6c] text-white rounded-md px-8 py-1 mt-5 text-xl text-center">Generating...</p>
                        :
                        <input type="submit" value="Generate!" className="bg-[#0EC420] text-white w-full rounded-md px-8 py-1 mt-5 text-xl "></input>
                    }
                </form>
            </div>
            <div className="bg-kc2-dark-gray p-3 mb-5 max-w-96 sm:max-w-6xl rounded-lg text-white">
                <h1 className="text-2xl pl-2">Guide:</h1>
                <div className="flex flex-col mt-3 gap-2">
                    <div>
                        <b>OS:</b> The target operating system. Check <a className="text-blue-400" href="https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63">GOOS</a> for more information.
                    </div>
                    <div>
                        <b>Architecture:</b> The target architecture. Check <a className="text-blue-400" href="https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63">GOARCH</a> for more information.
                    </div>
                    <div>
                        <b>Callback IP:</b> IP address of the C2 server. The implant will callback to this address to check for commands.
                    </div>
                    <div>
                        <b>Callback Port:</b> The port that listens for C2 traffic on the C2 server. 
                    </div>
                    <div>
                        <b>Implant Name:</b> Name of the implant. If left empty, a random name will be generated.
                    </div>
                    <div>
                        <b>Callback Interval:</b> Interval between implant check-ins in seconds. 
                    </div>
                    <div>
                        <b>Callback Jitter:</b> Variance between the different callback intervals in seconds.
                    </div>
                    <div>
                        <b>Reconnect Try Count:</b> Amount of times an implant will try to reconnect if it cannot connect to the C2 server. If it cannot contact the server within "Reconnect Try Count" times, the implants terminates.
                    </div>
                </div>
            </div>
        </div>
    )
}